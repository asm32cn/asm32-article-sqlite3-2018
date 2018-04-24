#!/usr/bin/env python
# -*- coding: utf-8 -*-

# Python 2.7
#

import web, datetime, os
import sqlite3

strBaseFolder = 'E:\\PASCAL\\asm32.article.sqlite3-20180111\\App_Data\\'
strFolder = strBaseFolder + 'f/'
strPubDataFile = strBaseFolder + 'asm32-article-web-pub.txt'
strDateFile = strBaseFolder + 'asm32.article.sqlite3'
strTable = 'table_article'

if not os.path.exists( strFolder ):
	os.makedirs( strFolder )

urls = (
	'/(\d*)', 'index',
	'/article-details-(\d+).html', 'details'
)

# http://echarts.baidu.com/feature.html
_index = '''
<html>
<head>
	<meta charset="UTF-8">
	<title>trans text</title>
	<style>
	body{background-color:#a5cbf7;}
	ul.page { display: block; width:100%%; height: 20px; clear:both; }
	ul.page li { display: block; width:30px; height: 20px; float:left; margin:5px; }
	ul.page li a { display: block; width:30px; height: 20px; text-align:center; float:left; border:1px solid #069; border-radius:5px; line-height: 20px; }
	li { line-height: 25px; }
	</style>
</head>
<body>


<h1>asm32.article.sqlite3</h1>


<ul class="page">
%s
</ul>

<ul>
%s
</ul>

<form method="POST">
	<p>
		<input type="submit" />
	</p>
	<dl>
		<dt>strTitle</dt>
		<dd><input type="text" name="txtTitle" value="" size="50" /></dd>
	</dl>
	<dl>
		<dt>strDate</dt>
		<dd><input type="text" name="txtDate" value="%s" /></dd>
	</dl>
	<dl>
		<dt>strFrom</dt>
		<dd><input type="text" name="txtFrom" value="" size="25" /></dd>
	</dl>
	<dl>
		<dt>strFromLink</dt>
		<dd><input type="text" name="txtFromLink" value="" size="50" /></dd>
	</dl>
	<dl>
		<dt>strContent</dt>
		<dd><textarea name="txtContent" cols="200" rows="20">%s</textarea></dd>
	</dl>
</form>
</body>
</html>
'''

def _now():
	dt = datetime.datetime.now() # datetime.datetime(2016, 5, 22, 3, 17, 3, 203000)
	return '%04d-%02d-%02d %02d:%02d:%02d' % (dt.year, dt.month, dt.day, dt.hour, dt.minute, dt.second)

def GetPubData():
	strContent = 'ICollection'
	if ( os.path.exists(strPubDataFile) ):
		with open(strPubDataFile, 'r') as fi:
			strContent = fi.read() 
	return strContent

def getConnection():
	return sqlite3.connect(strDateFile)

def HtmlEncode(s):
	return s.replace('&', '&amp;').replace('>', '&gt;').replace('<', '&lt;').replace('"', '&quot;') if s else ''

class index():

	def GET(self, pn=None):
		# return 'index.GET'
		# pn = web.input().get('pn')
		# pn = 0 if pn == None else int(pn)
		st = int(pn) if pn else 0

		conn = getConnection()
		cur = conn.cursor()

		strQuery = 'select count(*) from `%s`' % strTable
		nCount = int( cur.execute(strQuery).fetchone()[0] )

		if pn == '':
			st = (nCount - 1) / 20 * 20

		def getPageCtrl():
			return '\n'.join( [ '\t<li><a href="/%d">%d</a></li>' % (n, n) for n in xrange(0, nCount, 20) ] )

		strQuery = 'select `id`, `strTitle`, `strDate` from `%s` where `flag`=1 limit %d, %d' % (strTable, st, 20)
		_rows = cur.execute(strQuery)
		strContent = _index % (
			getPageCtrl(),
			'\n'.join( [ '\t<li>%s. <a href="/article-details-%s.html">%s</a> <em>%s</em></li>' % (_row[0], _row[0], HtmlEncode(_row[1]), _row[2]) for _row in _rows ] ).encode('utf-8'),
			_now(),
			GetPubData()
		)
		conn.close()
		return strContent

	def POST(self, id):

		strNow = _now()
		strFile = '%s.txt' % strNow.replace(':', '-')
		_input = web.input()

		def dbs(s):
			return '\'%s\'' % s.replace('\'', '\'\'') if s else 'null'

		def _get(s):
			return _input.get(s)

		strQuery = 'select max(id) from `%s`' % strTable

		conn = getConnection()
		cur = conn.cursor()
		nNewID = cur.execute(strQuery).fetchone()[0]
		nNewID = nNewID + 1 if nNewID else 1

		strQuery = 'insert into `%s`(`id`, `strTitle`, `strDate`, `strFrom`, `strFromLink`, `strContent`, `strDateCreated`, `flag`) values(%d, %s, %s, %s, %s, %s, %s, 1)' % (
			strTable,
			nNewID,
			dbs( _get('txtTitle') ),
			dbs( _get('txtDate') ),
			dbs( _get('txtFrom') ),
			dbs( _get('txtFromLink') ),
			dbs( _get('txtContent') ),
			dbs(strNow)
		)
		strQuery = strQuery #.decode('utf-8')

		# print strQuery

		# strContent = web.input().get('txtContent')
		# with open(strFolder + strFile, 'w') as fo:
		# 	fo.write( strContent.encode('utf-8') )

		# print 'Write, ' + strFile

		cur.execute(strQuery)
		conn.commit()

		conn.close()

		return 'index.POST'

class details():
	def GET(self, id=None):
		strContent = ''
		if id:
			strQuery = 'select `strTitle`, `strDate`, `strFrom`, `strFromLink`, `strContent`, `strDateCreated`, `strDateModified` from `%s` where `id`=%s' % (strTable, id)

			conn = getConnection()
			cur = conn.cursor()
			_row = cur.execute(strQuery).fetchone()
			if _row:
				strContent = '<html>\n' \
					'<head>\n' \
					'<meta charset="UTF-8">\n' \
					'<style>body{background-color:#a5cbf7;}</style>\n' \
					'</head>\n' \
					'<body>\n' \
					'<h1>%s</h1>\n' \
					'<p><strong>Date:</strong> %s</p>\n' \
					'<p><strong>From:</strong> %s</p>\n' \
					'<p><strong>FromLink:</strong> %s</p>\n' \
					'<pre>%s</pre>\n' \
					'<p><strong>DateCreated:</strong> %s</p>\n' \
					'<p><strong>DateModified:</strong> %s</p>\n' \
					'</body>\n' \
					'</html>' % (HtmlEncode(_row[0]), _row[1], HtmlEncode(_row[2]), HtmlEncode(_row[3]), HtmlEncode(_row[4]), _row[5], _row[6])
				strContent = strContent.encode('utf-8')
			else:
				strContent = 'Empty'
		else:
			strContent = 'None'
		return strContent

if __name__ == '__main__':
	web.application(urls, globals()).run()
