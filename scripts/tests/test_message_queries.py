"""Exercise the production list SQL and visibility predicates without a live database.

Run: python3 -m unittest discover -s scripts/tests -v
SQLite supports the SELECT/UNION/EXISTS subset used here; MySQL-specific locking
and deployment DDL still require integration validation on the target database.
"""
from pathlib import Path
import re
import sqlite3
import unittest

SOURCE = (Path(__file__).resolve().parents[2] / 'app/service/dao/message.go').read_text()
PERSONAL_SQL = re.search(r'db.Raw\(`(SELECT \*, .*?)`, owner, limit, offset\)', SOURCE, re.S).group(1)
LIST_SQL = re.search(r'db.Raw\(`(SELECT o\.\* FROM.*?)`, now, now, limit, offset\)', SOURCE, re.S).group(1)
VISIBILITY = re.search(r'Where\("(deleted_at IS NULL AND status.*?)", now, now\)', SOURCE).group(1)
NOW = '2026-10-31 12:00:00'


class MessageQueriesTest(unittest.TestCase):
    def setUp(self):
        self.db = sqlite3.connect(':memory:')
        self.db.row_factory = sqlite3.Row
        self.addCleanup(self.db.close)
        self.db.executescript('''
          CREATE TABLE message_centers (
            uid TEXT PRIMARY KEY, created_at TEXT, owner_user TEXT,
            customer_uuid TEXT, type INTEGER, event TEXT, content TEXT,
            is_read INTEGER, related_uuid TEXT);
          CREATE TABLE official_messages (
            uid TEXT PRIMARY KEY, title TEXT, content TEXT, status INTEGER,
            published_at TEXT, expires_at TEXT, deleted_at TEXT, sort INTEGER);
          CREATE TABLE official_message_reads (
            user_uuid TEXT, message_uuid TEXT, read_at TEXT,
            PRIMARY KEY (user_uuid, message_uuid));
          INSERT INTO message_centers VALUES
            ('p1', '2026-10-31 11:00:00', 'alice', 'customer', 1, '还款', 'paid', 0, 'pay'),
            ('p2', '2026-10-31 10:00:00', 'alice', 'customer', 2, '撤销', 'revoked', 1, 'pay'),
            ('other', '2026-10-31 11:00:00', 'bob', '', 1, '', '', 0, '');
        ''')
        for uid, status, published, expires, deleted, sort in [
            ('a1', 1, '2026-10-30 12:00:00', None, None, 0),
            ('a2', 1, '2026-10-29 12:00:00', None, None, 10),
            ('draft', 0, '2026-10-30 12:00:00', None, None, 0),
            ('future', 1, '2026-11-01 12:00:00', None, None, 0),
            ('expired', 1, '2026-10-30 12:00:00', NOW, None, 0),
            ('deleted', 1, '2026-10-30 12:00:00', None, NOW, 0),
        ]:
            self.db.execute('INSERT INTO official_messages VALUES (?, ?, ?, ?, ?, ?, ?, ?)',
                            (uid, uid, '公告内容', status, published, expires, deleted, sort))

    def messages(self, owner='alice', offset=0, limit=20):
        return self.db.execute(LIST_SQL, (NOW, NOW, limit, offset)).fetchall()

    def test_official_visibility_and_personal_user_isolation(self):
        rows = self.messages()
        self.assertEqual([r['uid'] for r in rows], ['a2', 'a1'])
        personal = self.db.execute(PERSONAL_SQL, ('alice', 20, 0)).fetchall()
        self.assertEqual([r['uid'] for r in personal], ['p1', 'p2'])
        self.assertEqual([r['source'] for r in personal], ['personal', 'personal'])
        self.assertEqual(rows[0]['title'], 'a2')

    def test_official_pagination_has_no_repeated_announcements(self):
        pages = [self.messages(offset=i, limit=1)[0]['uid'] for i in range(2)]
        self.assertEqual(pages, ['a2', 'a1'])
        self.assertEqual(self.messages(offset=2), [])
        self.assertEqual([r['uid'] for r in self.messages(offset=1, limit=2)], ['a1'])

    def test_bulk_read_is_owner_scoped_and_idempotent(self):
        predicate = re.search(r'Where\("(owner_user = \? AND is_read = 0)", owner\)', SOURCE).group(1)
        sql = 'UPDATE message_centers SET is_read = 1 WHERE ' + predicate
        self.assertEqual(self.db.execute(sql, ('alice',)).rowcount, 1)
        self.assertEqual(self.db.execute(sql, ('alice',)).rowcount, 0)
        self.assertEqual(self.db.execute("SELECT is_read FROM message_centers WHERE uid='other'").fetchone()[0], 0)
        self.assertEqual(self.db.execute('SELECT count(*) FROM official_message_reads').fetchone()[0], 0)



if __name__ == '__main__':
    unittest.main()
