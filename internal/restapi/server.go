package restapi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/RipperAcskt/crud-api-go/internal/model"
	"github.com/RipperAcskt/crud-api-go/internal/repo/postgres"
)

type AppHandler struct {
	postgres *postgres.Postgres
}

func New(p *postgres.Postgres) *AppHandler {
	var a AppHandler
	a.postgres = p
	return &a
}

func (app AppHandler) Close() error {
	return app.postgres.Close()
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

func JsonUnmarshal(body io.ReadCloser, u *[]model.User) error {
	readedBody, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Errorf("readAll faild: %v", err)
	}

	err = json.Unmarshal(readedBody, &u)
	if err != nil {
		return fmt.Errorf("unmarshal faild: %v", err)
	}

	return nil
}

func JsonMarshalError(errText string) ([]byte, error) {
	message := make(map[string]string)
	message["Error"] = errText
	jsonResp, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	return jsonResp, nil
}
