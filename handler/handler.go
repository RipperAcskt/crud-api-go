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
	case http.MethodPost:
		app.createUsersHandler(w, req)
	case http.MethodDelete:
		app.deleteUserHandler(w, req)
	case http.MethodPut:
		app.updateUserHandler(w, req)
	}
}

func chekQuery(query url.Values) (bool, []int, error) {

	idStringMass := query["id"]
	if len(idStringMass) == 0 {
		return false, nil, nil
	}

	var idIntMass []int
	for _, idString := range idStringMass {
		idInt, err := strconv.Atoi(idString)
		if err != nil {
			return true, nil, err
		}
		idIntMass = append(idIntMass, idInt)
	}
	return true, idIntMass, nil
}

func (app AppHandler) getUsersHandler(w http.ResponseWriter, req *http.Request) {
	var users []json.Person
	var err error

	queryFlag, id, err := chekQuery(req.URL.Query())
	if err != nil {
		log.Fatal(err)
	}
	if queryFlag {
		users, err = db.SelectById(app.DB, id)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		users, err = db.SelectAll(app.DB)
		if err != nil {
			log.Fatal(err)
		}
	}

	resp, err := json.JsonMarshalResponse(users)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (app AppHandler) createUsersHandler(w http.ResponseWriter, req *http.Request) {
	var personToCreate json.Person

	errJson := json.JsonUnmarshal(req.Body, &personToCreate)

	if errJson != nil {
		log.Fatal(errJson)
		return
	}

	if err := validation(personToCreate); err != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")

		resp, err := json.JsonMarshalError(err)
		if err != nil {
			log.Fatal(err)
		}

		w.Write(resp)
		return
	}

	err := db.Create(app.DB, personToCreate)
	if err != nil {
		log.Fatal(err)
	}
}

func validation(p json.Person) string {
	var err string
	if p.Name == "" {
		err += "Fill name."
	}
	if p.Surname == "" {
		err += "Fill surname."
	}
	if p.Age <= 0 {
		err += "Age need to be upper zero."
	}
	return err
}

func (app AppHandler) deleteUserHandler(w http.ResponseWriter, req *http.Request) {

	queryFlag, id, err := chekQuery(req.URL.Query())
	if err != nil {
		log.Fatal(err)
	}
	if queryFlag {
		err := db.DeleteById(app.DB, id)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := db.DeleteAll(app.DB)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func (app AppHandler) updateUserHandler(w http.ResponseWriter, req *http.Request) {
	var personToUpdate json.Person

	errJson := json.JsonUnmarshal(req.Body, &personToUpdate)

	if errJson != nil {
		log.Fatal(errJson)
		return
	}

	if err := validation(personToUpdate); err != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")

		resp, err := json.JsonMarshalError(err)
		if err != nil {
			log.Fatal(err)
		}

		w.Write(resp)
		return
	}

	err := db.Update(app.DB, personToUpdate)
	if err != nil {
		log.Fatal(err)
	}
}
