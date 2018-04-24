#!/usr/bin/env python
# -*- coding: utf-8 -*-

# Python 2.7
# urllib-demo-1.py

import urllib

for i in xrange(49):
	url = 'http://localhost:8080/article-details-%s.html' % i
	_content = urllib.urlopen(url).read()

	# with open('f/' + url.split('/')[-1], 'wb') as fo:
	# 	fo.write(_content)

	print url
	# print _content.decode('utf-8')