package postgres

import (
	"database/sql"
	"fmt"

	"github.com/RipperAcskt/crud-api-go/internal/model"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func SelectAll(DB *sql.DB) ([]model.User, error) {
	rows, err := DB.Query("SELECT * FROM Person ORDER BY id")

	if err != nil {
		return nil, fmt.Errorf("query faild: %v", err)
	}

	defer rows.Close()

	var users []model.User
	var user model.User

	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Name, &user.LastName, &user.Age); err != nil {
			return nil, fmt.Errorf("scan faild: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func SelectById(DB *sql.DB, id []int) ([]model.User, error) {
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

	var users []model.User
	var user model.User

	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Name, &user.LastName, &user.Age); err != nil {
			return nil, fmt.Errorf("scan faild: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func Create(DB *sql.DB, p model.User) error {

	_, err := DB.Exec("INSERT INTO person(firstName, lastName, age) VALUES($1, $2, $3)", p.Name, p.LastName, p.Age)
	if err != nil {
		return fmt.Errorf("exec faild: %v", err)
	}
	return nil
}

func DeleteAll(DB *sql.DB) error {

	_, err := DB.Exec("DELETE FROM Person")
	if err != nil {
		return fmt.Errorf("exec faild: %v", err)
	}
	return nil
}

func DeleteById(DB *sql.DB, id []int) error {
	params := make([]interface{}, len(id))
	for i, v := range id {
		params[i] = v
	}

	request := "DELETE FROM Person WHERE id IN ("
	for i := 0; i < len(id); i++ {
		request += "$" + fmt.Sprint(i+1)
		if i != len(id)-1 {
			request += ", "
		}
	}
	request += ")"

	_, err := DB.Exec(request, params...)

	if err != nil {
		return fmt.Errorf("exec faild: %v", err)
	}
	return nil
}

func Update(DB *sql.DB, p model.User) error {

	_, err := DB.Exec("UPDATE Person SET firstName = $1, lastName = $2, age = $3 WHERE id = $4", p.Name, p.LastName, p.Age, p.Id)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("exec faild: %v\n", err), 500)
	}
	return nil
}
