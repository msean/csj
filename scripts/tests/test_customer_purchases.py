from pathlib import Path
import re
import sqlite3
import unittest

SOURCE = (Path(__file__).resolve().parents[2] / 'app/service/dao/customer_purchases.go').read_text()
BASE = re.search(r'Where\("(owner_user = .*?)", owner, customer\)', SOURCE).group(1)
SUMMARY = re.search(r'Select\("(COUNT\(\*\).*?)"\)', SOURCE).group(1)

class PurchasesTest(unittest.TestCase):
    def setUp(self):
        self.db = sqlite3.connect(':memory:')
        self.addCleanup(self.db.close)
        self.db.execute('CREATE TABLE orders(uid TEXT,owner_user TEXT,user_uuid TEXT,deleted_at TEXT,status INT,order_time INT,amount NUMERIC,credit_amount NUMERIC,tricycle_fee NUMERIC DEFAULT 0)')
        for i in range(25):
            self.db.execute('INSERT INTO orders(uid,owner_user,user_uuid,deleted_at,status,order_time,amount,credit_amount) VALUES(?,?,?,?,?,?,?,?)',(str(i),'owner','buyer',None,1,100+i,10.5,3 if i%2==0 else 0))
        for row in [('other-owner','else','buyer',None,1,100,999,99),('other-buyer','owner','else',None,1,100,999,99),('cancel','owner','buyer',None,100,100,999,99),('deleted','owner','buyer','deleted',1,100,999,99)]:
            self.db.execute('INSERT INTO orders(uid,owner_user,user_uuid,deleted_at,status,order_time,amount,credit_amount) VALUES(?,?,?,?,?,?,?,?)',row)
    def query(self,start=0,end=0,credit=False):
        clauses=[BASE];args=['owner','buyer']
        for text,param in [('order_time >= ?',start),('order_time < ?',end)]:
            if param:
                self.assertIn(text,SOURCE);clauses.append(text);args.append(param)
        if credit: clauses.append('credit_amount > 0')
        return self.db.execute('SELECT '+SUMMARY+' FROM orders WHERE '+' AND '.join(clauses),args).fetchone()
    def test_total_covers_all_pages_and_is_owner_scoped(self):
        self.assertEqual(self.query(),(25,262.5))
    def test_debt_filter_uses_remaining_amount(self):
        self.assertEqual(self.query(credit=True),(13,136.5))
    def test_time_range_is_start_inclusive_end_exclusive(self):
        self.assertEqual(self.query(100,110),(10,105))
        self.assertEqual(self.query(125,126),(0,0))

    def test_transport_fees_use_the_same_total_as_customer_summary(self):
        self.db.execute("UPDATE orders SET tricycle_fee=0.5 WHERE owner_user='owner' AND user_uuid='buyer' AND status=1 AND deleted_at IS NULL")
        self.assertEqual(self.query(),(25,275))
