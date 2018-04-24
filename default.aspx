<%@ Page Language="C#" ValidateRequest="false" %>
<%@ Import Namespace="System.Data" %>
<%@ Import Namespace="System.Data.SQLite" %>
<script language="C#" runat="server">
string strTable = "table_article";
string strConn;

private string getNow(){
	// return DateTime.Now.ToString("yyyy-MM-dd HH:mm:ss");
	return GetDateString( DateTime.Now );
}

private string GetDateString(DateTime dt){
	string s = null;
	if (dt != null) {
		s = dt.ToString("yyyy-MM-dd HH:mm:ss");
	}
	return s;
}

private string GetArticleList(SQLiteConnection conn, object pn){
	StringBuilder sb = new StringBuilder();
	sb.Append("\n\n<h1>asm32.article.sqlite3</h1>\n\n\n");
	int st = 0;
	int nCount = 0;

	try{
		string strQuery = string.Format("select count(*) from `{0}`", strTable);
		// Response.Write( strQuery );
		using(SQLiteCommand cmd = new SQLiteCommand(strQuery, conn)){
			using(SQLiteDataReader dr = cmd.ExecuteReader()){
				if(dr.Read()){
					nCount = dr.GetInt32(0);
				}
			}
		}
		sb.Append("<ul class=\"page\">\n");
		for(int i = 0; i < nCount; i += 20){
			sb.Append("\t<li><a href=\"?pn=").Append(i);
			sb.Append("\">").Append(i).Append("</a></li>\n");
		}
		sb.Append("</ul>\n\n\n");
		// Response.Write(nCount);

		st = pn == null ? (nCount - 1) / 20 * 20 : Convert.ToInt32(pn);

		strQuery = string.Format("select `id`, `strTitle`, `strDate` from `{0}` where `flag`=1 limit {1}, 20", strTable, st);
		// Response.Write( strQuery );
		sb.Append("<ul>\n");
		using(SQLiteCommand cmd = new SQLiteCommand(strQuery, conn)){
			using(SQLiteDataReader dr = cmd.ExecuteReader()){
				while(dr.Read()){
					int id = dr.GetInt32(0);
					sb.Append("\t<li>").Append(id).Append(". <a href=\"?id=");
					sb.Append(id).Append("\">").Append( HtmlEncode( dr.GetString(1) ) );
					sb.Append("</a> <em>").Append(dr.GetString(2));
					sb.Append("</em></li>\n");
				}
			}
		}
		sb.Append("</ul>\n\n\n");

	}catch(Exception ex){
		Response.Write("Exception: " + ex.Message);
	}

	sb.Append("\n\n");

	sb.Append("<form method=\"POST\" action=\"?\">\n");
	sb.Append("<input type=\"hidden\" name=\"act\" value=\"da\" />\n");
	sb.Append("\t<p>\n");
	sb.Append("\t\t<input type=\"submit\" />\n");
	sb.Append("\t</p>\n");
	sb.Append("\t<dl>\n");
	sb.Append("\t\t<dt><strong>strTitle</strong></dt>\n");
	sb.Append("\t\t<dd><input type=\"text\" name=\"txtTitle\" value=\"\" size=\"50\" /></dd>\n");
	sb.Append("\t</dl>\n");
	sb.Append("\t<dl>\n");
	sb.Append("\t\t<dt><strong>strDate</strong></dt>\n");
	sb.Append("\t\t<dd><input type=\"text\" name=\"txtDate\" value=\"");
	sb.Append( getNow() ).Append("\" /></dd>\n");
	sb.Append("\t</dl>\n");
	sb.Append("\t<dl>\n");
	sb.Append("\t\t<dt><strong>strFrom</strong></dt>\n");
	sb.Append("\t\t<dd><input type=\"text\" name=\"txtFrom\" value=\"\" size=\"25\" /></dd>\n");
	sb.Append("\t</dl>\n");
	sb.Append("\t<dl>\n");
	sb.Append("\t\t<dt><strong>strFromLink</strong></dt>\n");
	sb.Append("\t\t<dd><input type=\"text\" name=\"txtFromLink\" value=\"\" size=\"50\" /></dd>\n");
	sb.Append("\t</dl>\n");
	sb.Append("\t<dl>\n");
	sb.Append("\t\t<dt><strong>strContent</strong></dt>\n");
	sb.Append("\t\t<dd><textarea name=\"txtContent\" cols=\"180\" rows=\"20\"></textarea></dd>\n");
	sb.Append("\t</dl>\n");
	sb.Append("</form>\n");
	sb.Append("<script src=\"/f/article-modify-util.js\"><").Append("/script>\n\n");

	return sb.ToString();
}

private string Dr_GetString(SQLiteDataReader dr, int n){
	return dr[n] == System.DBNull.Value ? "" : HtmlEncode( dr.GetString(n) );
}

private string HtmlEncode(string s){
	return System.Web.HttpUtility.HtmlEncode( s );
}

private string ArticleModifyForm(SQLiteConnection conn, int id){
	string strQuery = string.Format("select `strTitle`, `strDate`, `strFrom`, `strFromLink`, `strContent`, `strDateCreated`, `strDateModified` from `{0}` where `id`={1}", strTable, id);
	StringBuilder sb = new StringBuilder();
	using(SQLiteCommand cmd = new SQLiteCommand(strQuery, conn)){
		using(SQLiteDataReader dr = cmd.ExecuteReader()){
			if(dr.Read()){
				sb.Append("<h1>Modify article details</h1>\n");
				sb.Append("<form method=\"post\" act=\"?\">\n");
				sb.Append("<input type=\"hidden\" name=\"act\" value=\"dm\" />\n");
				sb.Append("<input type=\"hidden\" name=\"id\" value=\"").Append( id ).Append("\" />\n");
				sb.Append("<input type=\"submit\" />\n");
				sb.Append("<dl>\n");
				sb.Append("\t<dt><strong>strTitle:</strong></dt>\n");
				sb.Append("\t<dd><input type=\"text\" name=\"txtTitle\" value=\"");
				sb.Append( Dr_GetString(dr, 0) ).Append("\" size=\"50\" /></dd>\n");
				sb.Append("</dl>\n");

				sb.Append("<dl>\n");
				sb.Append("\t<dt><strong>strDate:</strong></dt>\n");
				sb.Append("\t<dd><input type=\"text\" name=\"txtDate\" value=\"");
				sb.Append( Dr_GetString(dr, 1) ).Append("\" /></dd>\n");
				sb.Append("</dl>\n");

				sb.Append("<dl>\n");
				sb.Append("\t<dt><strong>strFrom:</strong></dt>\n");
				sb.Append("\t<dd><input type=\"text\" name=\"txtFrom\" value=\"");
				sb.Append( Dr_GetString(dr, 2) ).Append("\" size=\"25\" /></dd>\n");
				sb.Append("</dl>\n");

				sb.Append("<dl>\n");
				sb.Append("\t<dt><strong>strFromLink:</strong></dt>\n");
				sb.Append("\t<dd><input type=\"text\" name=\"txtFromLink\" value=\"");
				sb.Append( Dr_GetString(dr, 3) ).Append("\" size=\"50\" /></dd>\n");
				sb.Append("</dl>\n");

				sb.Append("<dl>\n");
				sb.Append("\t<dt><strong>strContent:</strong></dt>\n");
				sb.Append("\t<dd><textarea name=\"txtContent\" cols=\"180\" rows=\"20\">");
				sb.Append( Dr_GetString(dr, 4) );
				sb.Append("</textarea></dl>\n");

				sb.Append("<p><strong>DateCreated:</strong> ").Append(Dr_GetString(dr, 5) ).Append("</p>\n");
				sb.Append("<p><strong>DateModified:</strong> ").Append(Dr_GetString(dr, 6) ).Append("</p>\n");
				sb.Append("</form>\n");
				sb.Append("<script src=\"/f/article-modify-util.js\"/></").Append("script>\n\n");
			}
		}
	}
	return sb.ToString();
}

private string GetArticleDetails(SQLiteConnection conn, int id){
	string strQuery = string.Format("select `strTitle`, `strDate`, `strFrom`, `strFromLink`, `strContent`, `strDateCreated`, `strDateModified` from `{0}` where `id`={1}", strTable, id);
	StringBuilder sb = new StringBuilder();
	using(SQLiteCommand cmd = new SQLiteCommand(strQuery, conn)){
		using(SQLiteDataReader dr = cmd.ExecuteReader()){
			if(dr.Read()){
				sb.Append("\n\n<span class=\"nav\" style=\"float:right\"><a href=\"?act=m&id=");
				sb.Append(id).Append("\" class=\"modify\"></a></span>\n\n\n");
				sb.Append("<h1>").Append(Dr_GetString(dr, 0) ).Append("</h1>\n");
				sb.Append("<p><strong>Date:</strong> ").Append(Dr_GetString(dr, 1) ).Append("</p>\n");
				sb.Append("<p><strong>From:</strong> ").Append(Dr_GetString(dr, 2) ).Append("</p>\n");
				sb.Append("<p><strong>FromLink:</strong> ").Append(Dr_GetString(dr, 3) ).Append("</p>\n");
				sb.Append("<pre>").Append(Dr_GetString(dr, 4) ).Append("</pre>\n");
				sb.Append("<p><strong>DateCreated:</strong> ").Append(Dr_GetString(dr, 5) ).Append("</p>\n");
				sb.Append("<p><strong>DateModified:</strong> ").Append(Dr_GetString(dr, 6) ).Append("</p>\n\n\n");
			}
		}
	}
	return sb.ToString();
}

private string DbString(string s){
	return s == null || string.IsNullOrEmpty(s) ? "null" : "'" + s.Replace("'", "''") + "'";
}

private void DoSaveArticle(SQLiteConnection conn){
	string strQuery;

	Response.Write( "POST\n" );

	object act = Request.Form["act"];
	if(act == null) return;

	string strAct = act.ToString();
	string strDate = string.Empty;
	bool canSave = true;
	if( act.Equals("da") || act.Equals("dm") ){
		try{
			strDate = Convert.ToDateTime(Request.Form["txtDate"]).ToString("yyyy-MM-dd HH:mm:ss");
		}catch(Exception ex){
			// Response.Write("Exception: " + ex.Message);
			strDate = null;
		}
	}

	if ( act.Equals("da") ){

		int nMaxID = 1;

		strQuery = string.Format("select max(id) from `{0}`", strTable);
		using(SQLiteCommand cmd = new SQLiteCommand(strQuery, conn)){
			using(SQLiteDataReader dr = cmd.ExecuteReader()){
				if(dr.Read()){
					nMaxID = dr.GetInt32(0) + 1;
				}
			}


			if( canSave) {
				strQuery = string.Format("insert into `{0}`(`id`, `strTitle`, `strDate`, `strFrom`, " +
					"`strFromLink`, `strContent`, `strDateCreated`, `flag`)" +
					" values({1}, {2}, {3}, {4}, {5}, {6}, {7}, 1)",
					strTable, nMaxID,
					DbString(Request.Form["txtTitle"]),
					DbString(strDate),
					DbString(Request.Form["txtFrom"]),
					DbString(Request.Form["txtFromLink"]),
					DbString(Request.Form["txtContent"]),
					DbString( getNow() )
				);
				cmd.CommandText = strQuery;
				// Response.Write(strQuery);
				cmd.ExecuteNonQuery();
			}
		}
	} else if(strAct.Equals("dm")) {
		object id = Request.Form["id"];
		if ( id != null ){
			int nItemID = Convert.ToInt32(id);
			strQuery = string.Format("update `{0}` set `strTitle`={2}, `strDate`={3}, `strFrom`={4}, " +
				"`strFromLink`={5}, `strContent`={6}, `strDateModified`={7} where `id` = {1}",
				strTable, nItemID,
				DbString(Request.Form["txtTitle"]),
				DbString(strDate),
				DbString(Request.Form["txtFrom"]),
				DbString(Request.Form["txtFromLink"]),
				DbString(Request.Form["txtContent"]),
				DbString( getNow() )
			);
			using(SQLiteCommand cmd = new SQLiteCommand(strQuery, conn)){
				cmd.ExecuteNonQuery();
			}
		}
		Response.Write("dm");
	}
}

protected void Page_Load(Object sender, EventArgs e){

	string strQuery;
	strConn = "Data Source=" + Server.MapPath("App_Data\\asm32.article.sqlite3");
	// Response.Write( strConn );

	using( SQLiteConnection conn = new SQLiteConnection(strConn) ){
		conn.Open();

		if( Request.HttpMethod.Equals("POST") ){
			DoSaveArticle(conn);
			Response.End();
		}

		object id = Request.QueryString["id"];
		object pn = Request.QueryString["pn"];
		object act = Request.QueryString["act"];

		StringBuilder sb = new StringBuilder();

		sb.Append("<!DOCTYPE html>\n");
		sb.Append("<html xmlns=\"http://www.w3.org/1999/xhtml\">\n");
		sb.Append("<head>\n");
		sb.Append("	<meta http-equiv=\"Content-Type\" content=\"text/html;charset=utf-8\"/>\n");
		sb.Append("	<title>asm32.article.sqlite3</title>\n");
		sb.Append("	<style>\n");
		sb.Append("	body { background-color:#a5cbf7; }\n");
		sb.Append("	ul.page { display: block; width:100%; height: 20px; clear:both; }\n");
		sb.Append("	ul.page li, .nav { display: block; width:30px; height: 20px; float:left; margin:5px; }\n");
		sb.Append("	ul.page li a, .nav a { display: block; width:30px; height: 20px; text-align:center; float:left; border:1px solid #069; border-radius:5px; line-height: 20px; }\n");
		sb.Append("	li { line-height: 25px; }\n");
		sb.Append("	.nav a.modify::after { content:'m'; }\n");
		sb.Append("	</style>\n");
		sb.Append("</head>\n");
		sb.Append("<body>\n");


		if( id == null ){
			sb.Append( GetArticleList(conn, pn) );
		} else if( act != null ) {
			string strAct = act.ToString();
			if( strAct.Equals("m") ){
				sb.Append( ArticleModifyForm(conn, Convert.ToInt32(id) ) );
			} else {
				sb.Append("act=" + strAct);
			}
		} else {
			sb.Append( GetArticleDetails(conn, Convert.ToInt32(id) ) );
		}

		sb.Append("</body>\n");
		sb.Append("</html>\n");
		Response.Write( sb.ToString() );
	}

}

</script>