package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type database struct {
	dbObject *sql.DB
}

type person struct {
	Id      int
	Name    string `json:"name"`
	Surname string `json:"lastName"`
	Age     int    `json:"age"`
}

func (db database) list(w http.ResponseWriter, req *http.Request) {
	var id, age int
	var firstName, lastName string

	rows, err := db.dbObject.Query("SELECT * FROM Person ORDER BY id")

	if err != nil {
		log.Fatalf("Error while doing request to database for output table: %v\n", err)
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&id, &firstName, &lastName, &age); err != nil {
			fmt.Fprintf(w, "%v\n", err)
			continue
		}
		fmt.Fprintf(w, "Id: %d\nName: %s\nSurname: %s\nAge: %d\n", id, firstName, lastName, age)
	}
}

func (db database) create(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while reading body: %v", err), 500)
		return
	}

	var requestJson person

	err = json.Unmarshal(body, &requestJson)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while unmarshal: %v", err), 500)
		return
	}

	var badRequest string

	if requestJson.Name == "" {
		badRequest += "Fill name\n"
	}
	if requestJson.Surname == "" {
		badRequest += "Fill surname\n"
	}
	if requestJson.Age <= 0 {
		badRequest += "Age need to be upper zero\n"
	}
	if badRequest != "" {
		http.Error(w, badRequest, http.StatusBadRequest)
		return
	}

	_, err = db.dbObject.Exec("INSERT INTO person(firstName, lastName, age) VALUES($1, $2, $3)", requestJson.Name, requestJson.Surname, requestJson.Age)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while creating person: %v", err), 500)
		return
	}

}

func (db database) delete(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while reading body: %v", err), 500)
		return
	}

	var deleteRequest struct {
		Id             int
		DeleteAllTable bool `json:"all"`
	}

	err = json.Unmarshal(body, &deleteRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while unmarshal: %v", err), 500)
		return
	}

	if deleteRequest.Id <= 0 && !deleteRequest.DeleteAllTable {
		http.Error(w, "Id should be upper than zero", http.StatusBadRequest)
		return
	}

	stmt, err := db.dbObject.Prepare("DELETE FROM Person WHERE id = $1")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while preparring request to database: %v", err), 500)
		return
	}
	defer stmt.Close()

	if deleteRequest.DeleteAllTable {
		_, err := db.dbObject.Exec("DELETE FROM Person")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while deleting all information from table: %v\n", err), 500)
			return
		}
	}
	_, err = stmt.Exec(deleteRequest.Id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while deleting id-%d: %v\n", deleteRequest.Id, err), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["Status"] = "Deleted"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while marshalling: %v\n", err), 500)
		return
	}
	w.Write(jsonResp)
}

func (db database) updateAll(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while reading body: %v", err), 500)
		return
	}

	var personToUpdate person

	err = json.Unmarshal(body, &personToUpdate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while unmarshal: %v", err), 500)
		return
	}

	if personToUpdate.Name == "" || personToUpdate.Surname == "" {
		http.Error(w, "All fields should be filled in", http.StatusBadRequest)
		return
	}

	_, err = db.dbObject.Exec("UPDATE Person SET firstName = $1, lastName = $2, age = $3 WHERE id = $4", personToUpdate.Name, personToUpdate.Surname, personToUpdate.Age, personToUpdate.Id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while updating: %v\n", err), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["Status"] = "Updated"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while marshalling: %v\n", err), 500)
		return
	}
	w.Write(jsonResp)

}

func main() {
	url := "postgres://ripper:150403@localhost:5432/ripper"

	db := database{}
	db.dbObject = openDB(url)
	defer db.dbObject.Close()

	mux := http.NewServeMux()
	mux.Handle("/list", http.HandlerFunc(db.list))
	mux.Handle("/create", http.HandlerFunc(db.create))
	mux.Handle("/delete", http.HandlerFunc(db.delete))
	mux.Handle("/updateAll", http.HandlerFunc(db.updateAll))
	log.Fatal(http.ListenAndServe("localhost:8080", mux))

}

func openDB(url string) *sql.DB {
	db, err := sql.Open("pgx", url)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return db
}
