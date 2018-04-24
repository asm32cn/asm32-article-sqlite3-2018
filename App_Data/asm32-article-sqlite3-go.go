// asm32-article-sqlite3-go.go
package main

import (
	"fmt"
	"io"
	"strconv"
	"log"
	"time"
	"regexp"
	"net/http"
	"strings"
	"database/sql"
	"encoding/base64"
	_ "github.com/mattn/go-sqlite3"
)

const (
	strTitle = "asm32.article.sqlite3"
	strHtmlTemplate = "<!DOCTYPE html>\n" +
		"<html xmlns=\"http://www.w3.org/1999/xhtml\">\n" +
		"<head>\n" +
		"	<meta http-equiv=\"Content-Type\" content=\"text/html;charset=utf-8\"/>\n" +
		"	<title>%s</title>\n" +
		"	<style>\n" +
		"	body { background-color:#a5cbf7; }\n" +
		"	ul.page { display: block; width:100%%; height: 20px; clear:both; }\n" +
		"	ul.page li, .nav { display: block; width:30px; height: 20px; float:left; margin:5px; }\n" +
		"	ul.page li a, .nav a { display: block; width:30px; height: 20px; text-align:center; float:left; border:1px solid #069; border-radius:5px; line-height: 20px; }\n" +
		"	li { line-height: 25px; }\n" +
		"	</style>\n" +
		"</head>\n" +
		"<body>\n\n\n\n" +
		"%s\n\n\n\n" +
		"</body>\n</html>"
	strDateTimeFormat = "2006-01-02 15:04:05"
	strDatabase = "./asm32.article.sqlite3"
	strTableName = "table_article"
)

type NullTime struct {
	Time time.Time
	Valid bool // Valid is true if Time is not NULL
}
// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
    nt.Time, nt.Valid = value.(time.Time)
    return nil
}
// // Value implements the driver Valuer interface.
// func (nt NullTime) Value() (driver.Value, error) {
//     if !nt.Valid {
//         return nil, nil
//     }
//     return nt.Time, nil
// }

func Base64Decode(s string) string {
	bs, err := base64.StdEncoding.DecodeString(s)
	ns := ""
	if err == nil {
		ns = string(bs)
	}
	return ns
}

func getNowString() string {
	return time.Now().Format(strDateTimeFormat)
}

func ConvertTimeString(tm NullTime) string {
	var ts string
	if tm.Valid {
		ts = tm.Time.Format(strDateTimeFormat)
	} else {
		ts = ""
	}
	return ts
}

var strTemplate_addForm = Base64Decode (
	"CgoKPGZvcm0gbWV0aG9kPSJQT1NUIiBhY3Rpb249Ij8iPgo8aW5wdXQgdHlwZT0iaGlkZGVuIiBuYW1l" +
	"PSJhY3QiIHZhbHVlPSJkYSIgLz4KCTxwPgoJCTxpbnB1dCB0eXBlPSJzdWJtaXQiIC8+Cgk8L3A+Cgk8" +
	"ZGw+CgkJPGR0PjxzdHJvbmc+c3RyVGl0bGU8L3N0cm9uZz48L2R0PgoJCTxkZD48aW5wdXQgdHlwZT0i" +
	"dGV4dCIgbmFtZT0idHh0VGl0bGUiIHZhbHVlPSIlcyIgc2l6ZT0iNTAiIC8+PC9kZD4KCTwvZGw+Cgk8" +
	"ZGw+CgkJPGR0PjxzdHJvbmc+c3RyRGF0ZTwvc3Ryb25nPjwvZHQ+CgkJPGRkPjxpbnB1dCB0eXBlPSJ0" +
	"ZXh0IiBuYW1lPSJ0eHREYXRlIiB2YWx1ZT0iJXMiIC8+PC9kZD4KCTwvZGw+Cgk8ZGw+CgkJPGR0Pjxz" +
	"dHJvbmc+c3RyRnJvbTwvc3Ryb25nPjwvZHQ+CgkJPGRkPjxpbnB1dCB0eXBlPSJ0ZXh0IiBuYW1lPSJ0" +
	"eHRGcm9tIiB2YWx1ZT0iJXMiIHNpemU9IjI1IiAvPjwvZGQ+Cgk8L2RsPgoJPGRsPgoJCTxkdD48c3Ry" +
	"b25nPnN0ckZyb21MaW5rPC9zdHJvbmc+PC9kdD4KCQk8ZGQ+PGlucHV0IHR5cGU9InRleHQiIG5hbWU9" +
	"InR4dEZyb21MaW5rIiB2YWx1ZT0iJXMiIHNpemU9IjUwIiAvPjwvZGQ+Cgk8L2RsPgoJPGRsPgoJCTxk" +
	"dD48c3Ryb25nPnN0ckNvbnRlbnQ8L3N0cm9uZz48L2R0PgoJCTxkZD48dGV4dGFyZWEgbmFtZT0idHh0" +
	"Q29udGVudCIgY29scz0iMjAwIiByb3dzPSIyMCI+JXM8L3RleHRhcmVhPjwvZGQ+Cgk8L2RsPgo8L2Zv" +
	"cm0+Cg==")
var strTemplate_modifyForm = Base64Decode(
	"Cgo8aDE+TW9kaWZ5IGFydGljbGUgZGV0YWlsczwvaDE+Cgo8Zm9ybSBtZXRob2Q9InBvc3QiIGFjdD0i" +
	"PyI+CjxpbnB1dCB0eXBlPSJoaWRkZW4iIG5hbWU9ImFjdCIgdmFsdWU9ImRtIiAvPgo8aW5wdXQgdHlw" +
	"ZT0iaGlkZGVuIiBuYW1lPSJpZCIgdmFsdWU9IiVkIiAvPgo8aW5wdXQgdHlwZT0ic3VibWl0IiAvPgo8" +
	"ZGw+Cgk8ZHQ+PHN0cm9uZz5zdHJUaXRsZTo8L3N0cm9uZz48L2R0PgoJPGRkPjxpbnB1dCB0eXBlPSJ0" +
	"ZXh0IiBuYW1lPSJ0eHRUaXRsZSIgdmFsdWU9IiVzIiBzaXplPSI1MCIgLz48L2RkPgo8L2RsPgo8ZGw+" +
	"Cgk8ZHQ+PHN0cm9uZz5zdHJEYXRlOjwvc3Ryb25nPjwvZHQ+Cgk8ZGQ+PGlucHV0IHR5cGU9InRleHQi" +
	"IG5hbWU9InR4dERhdGUiIHZhbHVlPSIlcyIgLz48L2RkPgo8L2RsPgo8ZGw+Cgk8ZHQ+PHN0cm9uZz5z" +
	"dHJGcm9tOjwvc3Ryb25nPjwvZHQ+Cgk8ZGQ+PGlucHV0IHR5cGU9InRleHQiIG5hbWU9InR4dEZyb20i" +
	"IHZhbHVlPSIlcyIgc2l6ZT0iMjUiIC8+PC9kZD4KPC9kbD4KPGRsPgoJPGR0PjxzdHJvbmc+c3RyRnJv" +
	"bUxpbms6PC9zdHJvbmc+PC9kdD4KCTxkZD48aW5wdXQgdHlwZT0idGV4dCIgbmFtZT0idHh0RnJvbUxp" +
	"bmsiIHZhbHVlPSIlcyIgc2l6ZT0iNTAiIC8+PC9kZD4KPC9kbD4KPGRsPgoJPGR0PjxzdHJvbmc+c3Ry" +
	"Q29udGVudDo8L3N0cm9uZz48L2R0PgoJPGRkPjx0ZXh0YXJlYSBuYW1lPSJ0eHRDb250ZW50IiBjb2xz" +
	"PSIyMDAiIHJvd3M9IjIwIj4lczwvdGV4dGFyZWE+PC9kbD4KPHA+PHN0cm9uZz5EYXRlQ3JlYXRlZDo8" +
	"L3N0cm9uZz4gJXM8L3A+CjxwPjxzdHJvbmc+RGF0ZU1vZGlmaWVkOjwvc3Ryb25nPiAlczwvcD4KCg==")

var regexNumber, _ = regexp.Compile(`\d+`)
var regexArticleList = regexp.MustCompile(`^/(\d*)$`)
var regexArticleDetails = regexp.MustCompile(`^/article-details-(\d+).html$`)
var regexArticleModify = regexp.MustCompile(`^/article-modify-(\d+).html$`)

func HtmlEncode(s string) string {
	ns := s
	ns = strings.Replace(ns, "&", "&amp;", -1)
	ns = strings.Replace(ns, ">", "&gt;", -1)
	ns = strings.Replace(ns, "<", "&lt;", -1)
	ns = strings.Replace(ns, "\"", "&quot;", -1)
	return ns
}

func HtmlEncodeNS(s sql.NullString) string {
	if s.Valid {
		return HtmlEncode(s.String)
	}
	return ""
}

func DoRenderHtml(w http.ResponseWriter, strPageTitle string, htmlBody string) {
	strHtmlData := fmt.Sprintf(strHtmlTemplate, HtmlEncode(strPageTitle), htmlBody)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, strHtmlData)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	strHtmlBody := "<h1>" + strTitle + "</h1>\n\n\n"

	nCount := 0
	var strQuery string;
	db, err := sql.Open("sqlite3", strDatabase)
	if err != nil {
		strHtmlBody += "Error: " + err.Error()
		DoRenderHtml(w, strTitle, strHtmlBody)
		return
	}

	strQuery = fmt.Sprintf("select count(*) from `%s`", strTableName)
	rows, err := db.Query(strQuery)
	if err != nil {
		strHtmlBody += "Error: " + err.Error()
	} else if rows.Next() {
		err = rows.Scan(&nCount)
		// println("nCount =", nCount)
		rows.Close()
	}

	strHtmlBody += fmt.Sprintf("<!-- nCount = %d -->\n\n", nCount) +
		"<ul class=\"page\">\n"
	for i := 0; i < nCount; i += 20 {
		strHtmlBody += fmt.Sprintf("\t<li><a href=\"/%d\">%d</a></li>\n", i, i)
	}
	strHtmlBody += "</ul>\n\n\n"

	pn := 0

	pn, err = strconv.Atoi( regexNumber.FindString( r.URL.Path ) )
	// pn, err = strconv.Atoi( r.Form.Get("pn") )
	if err != nil {
		pn = (nCount - 1) / 20 * 20
	}

	// println("pn =", pn)

	strQuery = fmt.Sprintf("select `id`, `strTitle`, `strDate` from `%s` where `flag`=1 limit %d, 20", strTableName, pn)
	rows, err = db.Query(strQuery)
	if err != nil {
		strHtmlBody += "Error: " + err.Error()
	} else {
		strHtmlBody += "<ul>\n"
		for rows.Next() {
			var id int
			var ns_strTitle sql.NullString
			var tm_strDate NullTime
			err = rows.Scan(&id, &ns_strTitle, &tm_strDate)
			strHtmlBody += fmt.Sprintf("\t<li>%d. <a href=\"/article-details-%d.html\">" +
				"%s</a> <em>%s</em></li>\n", id, id, HtmlEncodeNS(ns_strTitle), ConvertTimeString(tm_strDate))
		}
		strHtmlBody += "</ul>\n\n"
		rows.Close()
	}
	db.Close()

	strHtmlBody += fmt.Sprintf(strTemplate_addForm, "", getNowString(), "", "", "")

	DoRenderHtml(w, strTitle, strHtmlBody)
	return
}

func detailsHandler(w http.ResponseWriter, r *http.Request) {
	var strQuery string
	var strHtmlBody string

	s_strTitle := strTitle
	db, err := sql.Open("sqlite3", strDatabase)
	id, err := strconv.Atoi( regexNumber.FindString( r.URL.Path ) )
	if err != nil {
		strHtmlBody = "Error:1 " + err.Error()
	} else {
		strQuery = fmt.Sprintf("select `strTitle`, `strDate`, `strFrom`, `strFromLink`, " +
			"`strContent`, `strDateCreated`, `strDateModified` from `%s` where `id`=%d", strTableName, id)
		rows, err := db.Query(strQuery)
		if err != nil {
			strHtmlBody = "Error:2 " + err.Error()
		} else if rows.Next() {
			var ns_strTitle, ns_strFrom, ns_strFromLink, ns_strContent sql.NullString
			var tm_strDate, tm_strDateCreated, tm_strDateModified NullTime

			err = rows.Scan(&ns_strTitle, &tm_strDate, &ns_strFrom, &ns_strFromLink, &ns_strContent, &tm_strDateCreated, &tm_strDateModified)
			if err != nil {
				println("Error: " + err.Error())
			}
			s_strTitle = ns_strTitle.String
			strHtmlBody = fmt.Sprintf("<span class=\"nav\" style=\"float:right\">" +
				"<a href=\"/article-modify-%d.html\">m</a></span>\n" +
				"<h1>%s</h1>\n" +
				"<p><strong>Date:</strong> %s</p>\n" +
				"<p><strong>From:</strong> %s</p>\n" +
				"<p><strong>FromLink:</strong> %s</p>\n" +
				"<pre>%s</pre>\n" +
				"<p><strong>DateCreated:</strong> %s</p>\n" +
				"<p><strong>DateModified:</strong> %s</p>",
				id, HtmlEncodeNS(ns_strTitle), ConvertTimeString(tm_strDate),
				HtmlEncodeNS(ns_strFrom), HtmlEncodeNS(ns_strFromLink),
				HtmlEncodeNS(ns_strContent), ConvertTimeString(tm_strDateCreated), ConvertTimeString(tm_strDateModified))
			rows.Close()
		}
	}
	db.Close()

	DoRenderHtml(w, s_strTitle, strHtmlBody)
	return
}

func modifyFormHandler(w http.ResponseWriter, r *http.Request) {
	var strQuery string
	var strHtmlBody string

	s_strTitle := "Modify article details"
	db, err := sql.Open("sqlite3", strDatabase)
	id, err := strconv.Atoi( regexNumber.FindString( r.URL.Path ) )
	if err != nil {
		strHtmlBody = "Error:1 " + err.Error()
	} else {
		strQuery = fmt.Sprintf("select `strTitle`, `strDate`, `strFrom`, `strFromLink`, " +
			"`strContent`, `strDateCreated`, `strDateModified` from `%s` where `id`=%d", strTableName, id)
		rows, err := db.Query(strQuery)
		if err != nil {
			strHtmlBody = "Error:2 " + err.Error()
		} else if rows.Next() {
			var ns_strTitle, ns_strFrom, ns_strFromLink, ns_strContent sql.NullString
			var tm_strDate, tm_strDateCreated, tm_strDateModified NullTime

			err = rows.Scan(&ns_strTitle, &tm_strDate, &ns_strFrom, &ns_strFromLink, &ns_strContent, &tm_strDateCreated, &tm_strDateModified)
			if err != nil {
				println("Error: " + err.Error())
			}
			// s_strTitle = ns_strTitle.String
			strHtmlBody = fmt.Sprintf(strTemplate_modifyForm,
				id, HtmlEncodeNS(ns_strTitle), ConvertTimeString(tm_strDate),
				HtmlEncodeNS(ns_strFrom), HtmlEncodeNS(ns_strFromLink),
				HtmlEncodeNS(ns_strContent), ConvertTimeString(tm_strDateCreated), ConvertTimeString(tm_strDateModified))
			rows.Close()
		}
	}
	db.Close()

	DoRenderHtml(w, s_strTitle, strHtmlBody)
	return
}

func _dbs(s string) string {
	ns := strings.TrimSpace(s)
	// if ns == "" {
	// 	ns = "null"
	// } else {
	// 	ns = "'" + strings.Replace(ns, "'", "''", -1) + "'"
	// }
	return ns
}

func DoUpdateDetails(w http.ResponseWriter, r *http.Request) {
	act := r.PostFormValue("act")
	// io.WriteString(w, act)
	if act != "da" && act != "dm" {
		return
	}

	var strQuery, strHtmlTitle, strHtmlBody string
	var nItemID int
	var s_strTitle, s_strDate, s_strFrom, s_strFromLink, s_strContent string

	strHtmlTitle = "Update article details - "

	db, err := sql.Open("sqlite3", strDatabase)
	if err != nil {
		strHtmlBody = "Error:1 " + err.Error()
	}

	s_strTitle    = _dbs(r.PostFormValue("txtTitle"))
	s_strDate     = _dbs(r.PostFormValue("txtDate"))
	s_strFrom     = _dbs(r.PostFormValue("txtFrom"))
	s_strFromLink = _dbs(r.PostFormValue("txtFromLink"))
	s_strContent  = _dbs(r.PostFormValue("txtContent"))

	if act == "da" {

		strHtmlTitle += "a"

		strQuery = fmt.Sprintf("select max(id) from `%s`", strTableName)
		rows, err := db.Query(strQuery)
		if err != nil {
			strHtmlBody = "Error:2 " + err.Error()
		} else if rows.Next() {
			err := rows.Scan(&nItemID)
			if err != nil {
				strHtmlBody = "Error:3 " + err.Error()
			} else {
				rows.Close()
				nItemID++

				strQuery = fmt.Sprintf("insert into `%s`(`id`, `strTitle`, `strDate`, `strFrom`, " +
					"`strFromLink`, `strContent`, `strDateCreated`, `flag`)" +
					" values(?, ?, ?, ?, ?, ?, ?, 1)", strTableName)

				stmt, err := db.Prepare(strQuery)
				if err != nil {
					strHtmlBody = "Error:4 " + err.Error()
				} else {
					res, err := stmt.Exec(nItemID, s_strTitle, s_strDate, s_strFrom, s_strFromLink, s_strContent, getNowString())
					if err != nil {
						strHtmlBody = "Error:5 " + err.Error()
					} else {
						affect, err := res.RowsAffected()
						if err != nil {
							strHtmlBody = "Error:6 " + err.Error()
						} else {
							// fmt.Println(affect)
							strHtmlBody = fmt.Sprintf("Update article details OK! affect %d", affect)
						}
					}
				}

			}
		}

	} else if act == "dm" {

		nItemID, _ = strconv.Atoi( r.PostFormValue("id") )

		strHtmlTitle += fmt.Sprintf("m (%d)", nItemID)

		strQuery = fmt.Sprintf("update `%s` set `strTitle`=?, `strDate`=?, `strFrom`=?, " +
			"`strFromLink`=?, `strContent`=?, `strDateModified`=? where `id`=?", strTableName)

		stmt, err := db.Prepare(strQuery)
		if err != nil {
			strHtmlBody = "Error:4 " + err.Error()
		} else {
			res, err := stmt.Exec(s_strTitle, s_strDate, s_strFrom, s_strFromLink, s_strContent, getNowString(), nItemID)
			if err != nil {
				strHtmlBody = "Error:5 " + err.Error()
			} else {
				affect, err := res.RowsAffected()
				if err != nil {
					strHtmlBody = "Error:6 " + err.Error()
				} else {
					// fmt.Println(affect)
					strHtmlBody = fmt.Sprintf("Update article details OK! affect %d", affect)
				}
			}
		}
	}

	db.Close()

	DoRenderHtml(w, strHtmlTitle, strHtmlBody)

	return
}

func faviconIconHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func routeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		DoUpdateDetails(w, r)
		return
	}

	switch {
	case regexArticleList.MatchString(r.URL.Path):
		indexHandler(w, r)
	case regexArticleDetails.MatchString(r.URL.Path):
		detailsHandler(w, r)
	case regexArticleModify.MatchString(r.URL.Path):
		modifyFormHandler(w, r)
	default:
		io.WriteString(w, r.URL.Path)
	}
}

func main() {
	// http.HandleFunc("/", indexHandler)
	http.HandleFunc("/", routeHandler)
	http.HandleFunc("/favicon.ico", faviconIconHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}