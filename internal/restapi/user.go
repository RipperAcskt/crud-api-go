package restapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/RipperAcskt/crud-api-go/internal/model"
	"github.com/RipperAcskt/crud-api-go/internal/repo/postgres"
)

func (app AppHandler) getUsersHandler(w http.ResponseWriter, req *http.Request) {
	var users []model.User
	var err error

	queryFlag, id, err := chekQuery(req.URL.Query())
	if err != nil {
		log.Fatal(err)
	}
	if queryFlag {
		users, err = postgres.SelectById(app.DB, id)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		users, err = postgres.SelectAll(app.DB)
		if err != nil {
			log.Fatal(err)
		}
	}

	resp, err := json.Marshal(users)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (app AppHandler) createUsersHandler(w http.ResponseWriter, req *http.Request) {
	var userToCreate []model.User

	errJson := JsonUnmarshal(req.Body, &userToCreate)

	if errJson != nil {
		log.Fatal(errJson)
		return
	}
	for _, u := range userToCreate {
		if err := u.Validation(); err != "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")

			resp, err := JsonMarshalError(err)
			if err != nil {
				log.Fatal(err)
			}

			w.Write(resp)
			return
		}

		err := postgres.Create(app.DB, u)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func (app AppHandler) deleteUserHandler(w http.ResponseWriter, req *http.Request) {

	queryFlag, id, err := chekQuery(req.URL.Query())
	if err != nil {
		log.Fatal(err)
	}
	if queryFlag {
		err := postgres.DeleteById(app.DB, id)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := postgres.DeleteAll(app.DB)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func (app AppHandler) updateUserHandler(w http.ResponseWriter, req *http.Request) {
	var userToUpdate []model.User

	errJson := JsonUnmarshal(req.Body, &userToUpdate)

	if errJson != nil {
		log.Fatal(errJson)
		return
	}

	for _, u := range userToUpdate {
		if err := u.Validation(); err != "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")

			resp, err := JsonMarshalError(err)
			if err != nil {
				log.Fatal(err)
			}

			w.Write(resp)
			return
		}

		err := postgres.Update(app.DB, u)
		if err != nil {
			log.Fatal(err)
		}
	}

}
