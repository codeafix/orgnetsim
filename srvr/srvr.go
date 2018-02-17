package srvr

import (
	"fmt"
	"net/http"

	"github.com/spaceweasel/mango"
)

//ListenAndServe launches the web server
//rootpath is the root directory where data served by this server is persisted
//port is the port to listen on
func ListenAndServe(rootpath string, port string) {
	fm := NewFileManager(rootpath)
	r := CreateRouter(fm)
	r.RequestLogger = func(l *mango.RequestLog) {
		fmt.Println(l.CombinedFormat())
	}
	r.ErrorLogger = func(err error) {
		fmt.Println(err.Error())
	}
	http.ListenAndServe(":"+port, r)
}

//CreateRouter registers the route handlers. This function allows the route handlers
//to be tested with the mango.Browser
func CreateRouter(fm FileManager) *mango.Router {
	r := mango.NewRouter()

	r.RegisterModules([]mango.Registerer{
		NewSimListHandler(fm),
		NewSimHandler(fm),
		NewStepHandler(fm),
	})

	return r
}
