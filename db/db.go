package db

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/RipperAcskt/crud-api-go/json"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Database struct {
	DbObject *sql.DB
}

func (db Database) List(w http.ResponseWriter, req *http.Request) {
	var id, age int
	var firstName, lastName string

	rows, err := db.DbObject.Query("SELECT * FROM Person ORDER BY id")

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

func (db Database) Create(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}

	var personToCreate json.Person

	errJson := json.JsonUnmarshal(req.Body, &personToCreate, false)

	if errJson != "" {
		http.Error(w, errJson, 500)
		return
	}

	var badRequest string

	if personToCreate.Name == "" {
		badRequest += "Fill name\n"
	}
	if personToCreate.Surname == "" {
		badRequest += "Fill surname\n"
	}
	if personToCreate.Age <= 0 {
		badRequest += "Age need to be upper zero\n"
	}
	if badRequest != "" {
		http.Error(w, badRequest, http.StatusBadRequest)
		return
	}

	_, err := db.DbObject.Exec("INSERT INTO person(firstName, lastName, age) VALUES($1, $2, $3)", personToCreate.Name, personToCreate.Surname, personToCreate.Age)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while creating person: %v", err), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := json.JsonMarshalResponse("Created")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while marshalling: %v\n", err), 500)
		return
	}
	w.Write(response)

}

func (db Database) Delete(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}

	var deleteRequest json.DeleteInfo

	errJson := json.JsonUnmarshal(req.Body, &deleteRequest, true)

	if errJson != "" {
		http.Error(w, errJson, 500)
		return
	}

	if deleteRequest.Id <= 0 && !deleteRequest.DeleteAllTable {
		http.Error(w, "Id should be upper than zero", http.StatusBadRequest)
		return
	}

	stmt, err := db.DbObject.Prepare("DELETE FROM Person WHERE id = $1")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while preparring request to database: %v", err), 500)
		return
	}
	defer stmt.Close()

	if deleteRequest.DeleteAllTable {
		_, err := db.DbObject.Exec("DELETE FROM Person")
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
	response, err := json.JsonMarshalResponse("Deleted")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while marshalling: %v\n", err), 500)
		return
	}
	w.Write(response)
}

func (db Database) UpdateAll(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}

	var personToUpdate json.Person

	errJson := json.JsonUnmarshal(req.Body, &personToUpdate, false)

	if errJson != "" {
		http.Error(w, errJson, 500)
		return
	}

	if personToUpdate.Name == "" || personToUpdate.Surname == "" {
		http.Error(w, "All fields should be filled in", http.StatusBadRequest)
		return
	}

	_, err := db.DbObject.Exec("UPDATE Person SET firstName = $1, lastName = $2, age = $3 WHERE id = $4", personToUpdate.Name, personToUpdate.Surname, personToUpdate.Age, personToUpdate.Id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while updating: %v\n", err), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := json.JsonMarshalResponse("Updated all")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while marshalling: %v\n", err), 500)
		return
	}
	w.Write(response)

}

func (db Database) Update(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPatch {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}

	personToUpdate := json.Person{Age: -1}

	errJson := json.JsonUnmarshal(req.Body, &personToUpdate, false)

	if errJson != "" {
		http.Error(w, errJson, 500)
		return
	}

	var err error

	if personToUpdate.Name != "" && personToUpdate.Surname != "" && personToUpdate.Age != -1 {
		http.Error(w, "For updating all field go to /updateAll", http.StatusBadRequest)
		return
	}

	if personToUpdate.Name != "" {
		_, err = db.DbObject.Exec("UPDATE Person SET firstName = $1 WHERE id = $2", personToUpdate.Name, personToUpdate.Id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while updating: %v\n", err), 500)
			return
		}
	}
	if personToUpdate.Surname != "" {
		_, err = db.DbObject.Exec("UPDATE Person SET lastName = $1 WHERE id = $2", personToUpdate.Surname, personToUpdate.Id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while updating: %v\n", err), 500)
			return
		}
	}
	if personToUpdate.Age != -1 {
		_, err = db.DbObject.Exec("UPDATE Person SET age = $1 WHERE id = $2", personToUpdate.Age, personToUpdate.Id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error while updating: %v\n", err), 500)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := json.JsonMarshalResponse("Updated")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while marshalling: %v\n", err), 500)
		return
	}
	w.Write(response)
}

func (db Database) Home(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Hello on my server</h1>")
	fmt.Fprint(w, "<p>/list - watch all information</p>")
	fmt.Fprint(w, "<p>/create - create person</p>")
	fmt.Fprint(w, "<p>/delete - delete all or id person</p>")
	fmt.Fprint(w, "<p>/updateAll - update all field of person</p>")
	fmt.Fprint(w, "<p>/update - update fields of person seperatly</p>")
	fmt.Fprint(w, "<p>/help - for more information</p>")
}

func (db Database) Help(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "/list - Method: GET\n")
	fmt.Fprint(w, "/create - Method: POST\n\tNeed in request: name and surname and age\n")
	fmt.Fprint(w, "/delete - Method: DELETE\n\tNeed in request: id or flag all which show that all table should be cleared\n")
	fmt.Fprint(w, "/updateAll - Method: PUT\n\tNeed in request: name and surname and age\n")
	fmt.Fprint(w, "/update - Method: PATCH\n\tNeed in request: name or surname or age. Max two params")
}
