#!/usr/bin/env python
# -*- coding: utf-8 -*-

# Python 2.7
# asm32-article-sqlite3-sort.py

import sqlite3

f = 'asm32.article.sqlite3'
f2 = 'asm32.article.2.sqlite3'

strTable = 'table_article'


_createQuery = '''CREATE TABLE if not exists `%s`(
	`id` int,
	`strTitle` varchar(255),
	`strFrom` varchar(100) default null,
	`strFromLink` varchar(255) default null,
	`strContent` text,
	`strDate` datetime null,
	`strDateCreated` datetime null,
	`strDateModified` datetime null,
	`flag` int default 1,
	primary key(`id`)
);
''' % strTable

db = sqlite3.connect(f)

cur = db.cursor()

cur.execute('select `id`, `strTitle`, `strDate`, `strFrom`, `strFromLink`, `strContent`, `strDateCreated`, `flag` from `%s`' % strTable)
_rows = cur.fetchall()

db.close()

# print _rows

db = sqlite3.connect(f2)

cur = db.cursor()

cur.execute(_createQuery)
db.commit()

cur.execute('delete from `%s`' % strTable)
db.commit()

def dbs(s):
	return '\'%s\'' % s.replace('\'', '\'\'') if s else 'null'

n = 0

def dbInsertRow(r):
	global n
	n += 1
	strQuery = 'insert into `%s`(`id`, `strTitle`, `strDate`, `strFrom`, `strFromLink`, `strContent`, `strDateCreated`, `flag`) values(%d, %s, %s, %s, %s, %s, %s, 1)' % (
		strTable,
		n,
		dbs( r[1] ),
		dbs( r[2] ),
		dbs( r[3] ),
		dbs( r[4] ),
		dbs( r[5] ),
		dbs( r[6] )
	)
	cur.execute(strQuery)
	db.commit()

# _list = [112, 93, 114, 85, 47, 115, 87, 116, 78, 79, 99, 113]
_list = [1, 2, 3, 7, 4, 5, 6, 7, 8, 119]

for i in _list:
	for r in _rows:
		if r[0] == i:
			print r[0]
			dbInsertRow(r)

for r in _rows:
	if r[0] not in _list:
		print r[0]
		dbInsertRow(r)

db.close()