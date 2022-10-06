package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	url := "postgres://ripper:150403@localhost:5432/ripper"

	db := openDB(url)
	defer db.Close()

	var n int

L:
	for {
		fmt.Printf("1-Создать запись в таблице\n2-Вывести всю таблицу\n3-Модифицировать данные таблицы\n4-Удалить данные из таблицы\n5-Выход\n")
		fmt.Scanf("%d", &n)

		switch n {
		case 1:
			write(db)
		case 2:
			printAllTable(db)
		case 3:
			update(db)
		case 4:
			delete(db)
		case 5:
			break L
		}

	}

}

func openDB(url string) *sql.DB {
	db, err := sql.Open("pgx", url)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return db
}

func write(db *sql.DB) {
	var age int
	var firstName, lastName string

	fmt.Print("Имя: ")
	fmt.Scanln(&firstName)
	fmt.Print("Фамилия: ")
	fmt.Scanln(&lastName)
	fmt.Print("Возраст: ")
	fmt.Scanln(&age)

	_, err := db.Exec("INSERT INTO person(firstName, lastName, age) VALUES($1, $2, $3)", firstName, lastName, age)
	if err != nil {
		db.Close()
		log.Fatalf("Error while writing: %v\n", err)
	}
}

func printAllTable(db *sql.DB) {
	var id, age int
	var firstName, lastName string

	rows, err := db.Query("SELECT * FROM Person ORDER BY id")

	if err != nil {
		log.Fatalf("Error while doing request to database for output table: %v\n", err)
	}

	defer rows.Close()

	fmt.Println()
	for rows.Next() {
		if err = rows.Scan(&id, &firstName, &lastName, &age); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		fmt.Println(id, firstName, lastName, age)
	}
	fmt.Println()
}

func update(db *sql.DB) {
	var id, age, n int
	var firstName, lastName string

	printAllTable(db)
	fmt.Print("Id: ")
	fmt.Scanln(&id)
	fmt.Printf("1-Имя\n2-Фамилия\n3-Возвраст\n")
	fmt.Scanln(&n)

	switch n {
	case 1:
		fmt.Print("Имя: ")
		fmt.Scanln(&firstName)
		_, err := db.Exec("UPDATE Person SET firstName = $1 WHERE id = $2", firstName, id)
		if err != nil {
			log.Fatalf("Error while updating firstName: %v\n", err)
		}
	case 2:
		fmt.Print("Фамилия: ")
		fmt.Scanln(&lastName)
		_, err := db.Exec("UPDATE Person SET lastName = $1 WHERE id = $2", lastName, id)
		if err != nil {
			log.Fatalf("Error while updating lastName: %v\n", err)
		}
	case 3:
		fmt.Print("Возраст: ")
		fmt.Scanln(&age)
		_, err := db.Exec("UPDATE Person SET age = $1 WHERE id = $2", age, id)
		if err != nil {
			log.Fatalf("Error while updating age: %v\n", err)
		}

	}
}

func delete(db *sql.DB) {
	var id, maxId int
	printAllTable(db)

	db.QueryRow("SELECT MAX(id) FROM Person").Scan(&maxId)
	fmt.Printf("%v-Отчистить таблицу\n\n", maxId+1)

	fmt.Print("Id: ")
	fmt.Scanln(&id)

	stmt, err := db.Prepare("DELETE FROM Person WHERE id = $1")
	if err != nil {
		log.Fatalf("Error while preparing query: %v\n", err)
	}
	defer stmt.Close()

	if id > 0 && id <= maxId {
		_, err = stmt.Exec(id)
		if err != nil {
			log.Fatalf("Error while deleting: %v\n", err)
		}
	} else if id == maxId+1 {
		_, err := db.Exec("DELETE FROM Person")
		if err != nil {
			log.Fatalf("Error while deleting all information from table: %v\n", err)
		}
	}
}
