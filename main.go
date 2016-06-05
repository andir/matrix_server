package main

import (
        "net/http"
        "encoding/json"
	"strings"
	//"github.com/andir/matrix_server/matrix"
	m "github.com/andir/matrix_server/matrix"
	p "github.com/andir/matrix_server/matrix/protocol"
	"github.com/andir/matrix_server/matrix/protocol/client"
	"github.com/andir/matrix_server/matrix/protocol/client/login"
	"github.com/jinzhu/gorm"
	"github.com/andir/matrix_server/matrix/database"
)



type HandlerFunc func(http.ResponseWriter, *http.Request) (status_code int, obj interface{})


func WriteJSON(data interface{}, w http.ResponseWriter, req *http.Request) {
	if buffer, err := json.Marshal(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(buffer)
	}
}

func MakeSerializedResponseHandler(fun HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		status_code, response := fun(w, req)

		// if the handler doesn't return a valid status code assume 200
		if status_code == 0 {
			status_code = http.StatusOK
		}

		// if the status code isn't 400+ already check if it is an error
		// and adjust the status code to reflect it
		if status_code < 400 {
			if _, ok := response.(p.ErrorResponse); ok {
				status_code = http.StatusBadRequest
			}
		}

		// check the requested content type so we can marshal it as requested
		// We default to application/json to ease the debugging
		// FIXME: this probably should go at the top of this function?!
		content_type := strings.ToLower(req.Header.Get("Content-Type"))
		if content_type == "" {
			content_type = "application/json"
		}

		if content_type == "application/json" {
			w.WriteHeader(status_code)
			WriteJSON(response, w, req)
		} else {
			http.Error(w, "Invalid content type requsted", http.StatusNotAcceptable)
		}
	}
}

// Limit the HTTP methods that can be used with the given handler
func LimitToMethods(fun http.HandlerFunc, methods []string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		for _, b := range methods {
			if b == req.Method {
				fun(w, req)
				return
			}
		}
		http.Error(w, "Invalid Method", http.StatusNotAcceptable)
	}
}

func ClientVersionsHandler(w http.ResponseWriter, req *http.Request) (int, interface{}) {
        return 200, client.ClientVersionsResponse {
		Versions:[]string{"r.0.0.1"},
	}
}

func ClientLoginRequest(w http.ResponseWriter, req *http.Request) (int, interface{}) {
	loginRequest, err := login.ClientLoginRequestFromHTTPRequest(req)
	if err != nil {
		return 400, p.ErrorResponse{
			Errcode: m.M_NOT_JSON,
		}
	}
	return loginRequest.Handle()
}

func setupRoutes() {
        http.HandleFunc("/_matrix/client/versions", LimitToMethods(MakeSerializedResponseHandler(ClientVersionsHandler), []string{"GET"}))
	http.HandleFunc("/_matrix/client/ro0/login", LimitToMethods(MakeSerializedResponseHandler(ClientLoginRequest), []string{"POST"}))
}

var db *gorm.DB

func main() {
        setupRoutes()
	db = database.DatabaseInit()

        http.ListenAndServe(":8000", nil)

}