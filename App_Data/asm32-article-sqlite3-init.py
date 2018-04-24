#!/usr/bin/env python
# -*- coding: utf-8 -*-

# Python 2.7
# asm32-article-sqlite3-init.py

import sqlite3

f = 'asm32.article.sqlite3'

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

cur.execute(_createQuery)
db.commit()

# cur.execute('delete from `%s`' % strTable)
# db.commit()

db.close()
