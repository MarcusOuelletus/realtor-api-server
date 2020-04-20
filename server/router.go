package server

import (
	"encoding/json"
	"net/http"

	"github.com/MarcusOuelletus/rets-server/app/access"
	"github.com/MarcusOuelletus/rets-server/app/areas"
	"github.com/MarcusOuelletus/rets-server/app/frontend"
	"github.com/MarcusOuelletus/rets-server/app/mls"
	"github.com/MarcusOuelletus/rets-server/app/search"
	"github.com/MarcusOuelletus/rets-server/templates"
	"github.com/golang/glog"

	"github.com/gorilla/mux"
)

type handlerFuncTemplate = func(w http.ResponseWriter, r *http.Request)
type handlerFuncLogic = func(r *http.Request) *templates.APIResponse

func handlerTemplate(method handlerFuncLogic) handlerFuncTemplate {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		w.WriteHeader(http.StatusOK)

		// don't run method for pre-flight requests (chrome)
		if r.Method == "OPTIONS" {
			return
		}

		responseObject := method(r)

		json, err := json.Marshal(responseObject)

		if err != nil {
			glog.Errorf("error Marshalling response object: %s\n", err.Error())
		}

		w.Write(json)
	}
}

func generateRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/signup", handlerTemplate(access.Signup))
	r.HandleFunc("/search", handlerTemplate(search.Handler))
	r.HandleFunc("/areas", handlerTemplate(areas.Areas))
	r.HandleFunc("/mls", handlerTemplate(mls.MLS))
	r.HandleFunc("/js", frontend.JS)
	r.HandleFunc("/test", test)
	r.HandleFunc("/", frontend.Frontend)

	return r
}

func test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("success"))
}
