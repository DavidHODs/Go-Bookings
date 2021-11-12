package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	conn, err := sql.Open("pgx", "host=localhost port=5432 dbname=test_connect user=postgres password=password")
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to connect; %v\n", err))
	}
	defer conn.Close()
	log.Println("Connected to the database")

	err = conn.Ping()
	if err != nil {
		log.Fatal("Cannot ping database")
	}
	log.Println("Pinged the database")

	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	query := `insert into users (first_name, last_name) values ($1, $2)`
	_, err = conn.Exec(query, "Jack", "Brown")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted a Row")

	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	stmt := `update users set first_name = $1 where id = $2`
	_, err = conn.Exec(stmt, "Jackie", 5)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Updated one or more row")

	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}

	query = `select id, first_name, last_name from users where id =$1`
	var firstName, lastName string
	var id int
	row := conn.QueryRow(query, 1)
	err = row.Scan(&id, &firstName, &lastName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Query row returns", id, firstName, lastName)

	query = `delete from users where id = $1`
	_, err = conn.Exec(query, 6)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Deleted a Row")

	err = getAllRows(conn)
	if err != nil {
		log.Fatal(err)
	}
}

func getAllRows(conn *sql.DB) error {
	rows, err := conn.Query("select id, first_name, last_name from users")
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close()

	
	var firstName, lastName string
	var id int

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName)
		if err != nil {
			log.Println(err)
			return err
		}
		fmt.Println("Record is ", id, firstName, lastName)
	}

	if err = rows.Err(); err != nil {
		log.Fatal("Error scanning rows", err)
	}

	fmt.Println("-------------------------------------")

	return nil
}
