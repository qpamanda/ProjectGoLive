package server

import "net/http"

func addrequest(res http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(res, "addrequest.gohtml", nil)
}
