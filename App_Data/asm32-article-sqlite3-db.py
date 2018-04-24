#!/usr/bin/env python
# -*- coding: utf-8 -*-

# Python 2.7
#

import sqlite3

_db = 'asm32.article.sqlite3'

_conn = sqlite3.connect(_db)

print '_conn OK!'
_cur = _conn.cursor()

strQuery = '''create table table_article(
	id int,
	strTitle varchar(255),
	strFrom varchar(100) default null,
	strFromLink varchar(255) default null,
	strContent text,
	strDate datetime null,
	strDateCreated datetime null,
	strDateModified datetime null,
	flag int default 1,
	primary key(id)
)'''

try:
	_cur.execute(strQuery)

except Exception, ex:
	print ex

_conn.close()
