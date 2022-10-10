package handler

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/RipperAcskt/crud-api-go/db"
	"github.com/RipperAcskt/crud-api-go/json"
)

type AppHandler struct {
	DB *sql.DB
}

func (app AppHandler) Controller(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		app.getUsersHandler(w, req)

	}
}

func chekQuery(query url.Values) (bool, []int) {

	idStringMass := query["id"]
	if len(idStringMass) == 0 {
		return false, nil
	}

	var idIntMass []int
	for _, idString := range idStringMass {
		idInt, _ := strconv.Atoi(idString)
		idIntMass = append(idIntMass, idInt)
	}
	return true, idIntMass
}

func (app AppHandler) getUsersHandler(w http.ResponseWriter, req *http.Request) {
	users, err := db.SelectAll(app.DB)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := json.JsonMarshalResponse(users)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
