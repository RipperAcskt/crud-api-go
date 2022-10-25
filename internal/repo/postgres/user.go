package postgres

import (
	"fmt"

	"github.com/RipperAcskt/crud-api-go/internal/model"
)

func (p Postgres) SelectAll() ([]model.User, error) {
	rows, err := p.db.Query("SELECT * FROM Person ORDER BY id")

	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}

	defer rows.Close()

	var users []model.User
	var user model.User

	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Name, &user.LastName, &user.Age); err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (p Postgres) SelectById(id []int) ([]model.User, error) {
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

	rows, err := p.db.Query(request, params...)

	if err != nil {
		return nil, fmt.Errorf("query failed: %v", err)
	}

	defer rows.Close()

	var users []model.User
	var user model.User

	for rows.Next() {
		if err = rows.Scan(&user.Id, &user.Name, &user.LastName, &user.Age); err != nil {
			return nil, fmt.Errorf("scan failed: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (p Postgres) Create(u model.User) error {
	_, err := p.db.Exec("INSERT INTO person(firstName, lastName, age) VALUES($1, $2, $3)", u.Name, u.LastName, u.Age)
	if err != nil {
		return fmt.Errorf("exec failed: %v", err)
	}
	return nil
}

func (p Postgres) DeleteAll() error {
	_, err := p.db.Exec("DELETE FROM Person")
	if err != nil {
		return fmt.Errorf("exec failed: %v", err)
	}
	return nil
}

func (p Postgres) DeleteById(id []int) error {
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

	_, err := p.db.Exec(request, params...)

	if err != nil {
		return fmt.Errorf("exec failed: %v", err)
	}
	return nil
}

func (p Postgres) Update(u model.User) error {
	_, err := p.db.Exec("UPDATE Person SET firstName = $1, lastName = $2, age = $3 WHERE id = $4", u.Name, u.LastName, u.Age, u.Id)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("exec failed: %v\n", err), 500)
	}
	return nil
}
