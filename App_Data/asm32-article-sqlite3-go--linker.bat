@echo off

set strCmd=go build -ldflags="-w -s" asm32-article-sqlite3-go.go

echo #%strCmd%
%strCmd%

pause
