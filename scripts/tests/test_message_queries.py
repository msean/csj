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
LIST_SQL = re.search(r'db.Raw\(`(SELECT \* FROM .*?)`, owner, owner, now, now, limit, offset\)', SOURCE, re.S).group(1)
VISIBILITY = re.search(r'Where\("(deleted_at IS NULL AND status.*?)", now, now\)', SOURCE).group(1)
UNREAD = re.search(r'Where\("(NOT EXISTS.*?)", owner\)', SOURCE).group(1)
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
        return self.db.execute(LIST_SQL, (owner, owner, NOW, NOW, limit, offset)).fetchall()

    def unread(self, owner='alice'):
        return self.db.execute('SELECT count(*) FROM official_messages WHERE ' + VISIBILITY + ' AND ' + UNREAD,
                               (NOW, NOW, owner)).fetchone()[0]

    def test_pinned_sort_visibility_and_user_isolation(self):
        rows = self.messages()
        self.assertEqual([r['uid'] for r in rows], ['a2', 'a1', 'p1', 'p2'])
        self.assertEqual([r['source'] for r in rows], ['official', 'official', 'personal', 'personal'])
        self.assertEqual(rows[0]['title'], 'a2')
        self.assertEqual(rows[0]['type'], 3)
        self.assertEqual(self.unread(), 2)

    def test_combined_pagination_has_no_repeated_announcements(self):
        pages = [self.messages(offset=i, limit=1)[0]['uid'] for i in range(4)]
        self.assertEqual(pages, ['a2', 'a1', 'p1', 'p2'])
        self.assertEqual(self.messages(offset=4), [])
        self.assertEqual([r['uid'] for r in self.messages(offset=1, limit=2)], ['a1', 'p1'])

    def test_receipts_are_per_user_and_idempotent(self):
        for _ in range(2):
            self.db.execute('INSERT OR IGNORE INTO official_message_reads VALUES (?, ?, ?)', ('alice', 'a2', NOW))
        self.assertEqual(self.db.execute('SELECT count(*) FROM official_message_reads').fetchone()[0], 1)
        self.assertEqual([r['is_read'] for r in self.messages()], [1, 0, 0, 1])
        self.assertEqual(self.unread('alice'), 1)
        self.assertEqual(self.unread('bob'), 2)
        self.db.execute('INSERT INTO official_message_reads VALUES (?, ?, ?)', ('alice', 'a1', NOW))
        self.assertEqual(self.unread('alice'), 0)
        self.db.execute("UPDATE official_messages SET content = '改文案' WHERE uid = 'a1'")
        self.assertEqual(self.unread('alice'), 0)
        self.db.execute("UPDATE official_messages SET status = 0 WHERE uid = 'a2'")
        self.assertEqual([r['uid'] for r in self.messages()], ['a1', 'p1', 'p2'])


if __name__ == '__main__':
    unittest.main()
