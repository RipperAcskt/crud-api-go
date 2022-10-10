package db

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/RipperAcskt/crud-api-go/json"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Database struct {
	DbObject *sql.DB
}

func SelectAll(DB *sql.DB) ([]json.Person, error) {
	rows, err := DB.Query("SELECT * FROM Person ORDER BY id")

	if err != nil {
		return nil, fmt.Errorf("query faild: %v", err)
	}

	defer rows.Close()

	var users []json.Person
	var user json.Person

	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Name, &user.Surname, &user.Age); err != nil {
			return nil, fmt.Errorf("scan faild: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func SelectById(DB *sql.DB, id []int) ([]json.Person, error) {
	params := make([]interface{}, len(id))
	for i, v := range id {
		params[i] = v
	}

	request := "SELECT * FROM Person WHERE id IN ("
	for i := 0; i < len(id); i++ {
		request += "$" + fmt.Sprint(i+1)
		if i != len(id)-1 {
			request += ", "
		}
	}
	request += ") ORDER BY id"

	rows, err := DB.Query(request, params...)

	if err != nil {
		return nil, fmt.Errorf("query faild: %v", err)
	}

	defer rows.Close()

	var users []json.Person
	var user json.Person

	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Name, &user.Surname, &user.Age); err != nil {
			return nil, fmt.Errorf("scan faild: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func Create(DB *sql.DB, p json.Person) error {

	_, err := DB.Exec("INSERT INTO person(firstName, lastName, age) VALUES($1, $2, $3)", p.Name, p.Surname, p.Age)
	if err != nil {
		return fmt.Errorf("exec faild: %v", err)
	}
	return nil
}

func (db Database) Delete(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}

	var deleteRequest json.DeleteInfo

	// errJson := json.JsonUnmarshal(req.Body, &deleteRequest, true)

	// if errJson != "" {
	// 	http.Error(w, errJson, 500)
	// 	return
	// }

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

	// w.Header().Set("Content-Type", "application/json")
	// response, err := json.JsonMarshalResponse("Deleted")
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Error while marshalling: %v\n", err), 500)
	// 	return
	// }
	// w.Write(response)
}

func (db Database) UpdateAll(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}

	var personToUpdate json.Person

	// errJson := json.JsonUnmarshal(req.Body, &personToUpdate, false)

	// if errJson != "" {
	// 	http.Error(w, errJson, 500)
	// 	return
	// }

	if personToUpdate.Name == "" || personToUpdate.Surname == "" {
		http.Error(w, "All fields should be filled in", http.StatusBadRequest)
		return
	}

	_, err := db.DbObject.Exec("UPDATE Person SET firstName = $1, lastName = $2, age = $3 WHERE id = $4", personToUpdate.Name, personToUpdate.Surname, personToUpdate.Age, personToUpdate.Id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error while updating: %v\n", err), 500)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// response, err := json.JsonMarshalResponse("Updated all")
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Error while marshalling: %v\n", err), 500)
	// 	return
	// }
	// w.Write(response)

}

func (db Database) Update(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPatch {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}

	personToUpdate := json.Person{Age: -1}

	// errJson := json.JsonUnmarshal(req.Body, &personToUpdate, false)

	// if errJson != "" {
	// 	http.Error(w, errJson, 500)
	// 	return
	// }

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

	// w.Header().Set("Content-Type", "application/json")
	// response, err := json.JsonMarshalResponse("Updated")
	// if err != nil {
	// 	http.Error(w, fmt.Sprintf("Error while marshalling: %v\n", err), 500)
	// 	return
	// }
	// w.Write(response)
}
