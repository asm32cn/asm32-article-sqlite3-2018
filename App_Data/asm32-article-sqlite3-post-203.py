#!/usr/bin/env python
# -*- coding: utf-8 -*-

# Python 2.7
# asm32-article-sqlite3-post-120.py

import sqlite3, urllib

f = 'asm32.article.sqlite3'

conf_strUrl = 'http://10.52.9.203:8081/'

strTable = 'table_article'


db = sqlite3.connect(f)

cur = db.cursor()

cur.execute('select `id`, `strTitle`, `strDate`, `strFrom`, `strFromLink`, `strContent`, `strDateCreated`, `flag` from `%s`' % strTable)
_rows = cur.fetchall()

db.close()

def _dbs(s):
	return s.encode('utf-8') if s else ''

for r in _rows:
	# print r
	data = {
		'txtTitle': _dbs(r[1]),
		'txtDate': _dbs(r[2]),
		'txtFrom': _dbs(r[3]),
		'txtFromLink': _dbs(r[4]),
		'txtContent': _dbs(r[5])
	}

	print r[1]
	print urllib.urlopen(conf_strUrl, urllib.urlencode(data)).read()

db.close()