package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
}

func main() {

	cfg := PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "baloo",
		Password: "junglebook",
		DBName:   "lenslocked",
		SSLMode:  "disable",
	}

	db, err := sql.Open("pgx", cfg.String())
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected!")

	// Create table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users(
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT UNIQUE NOT NULL
		);

		CREATE TABLE IF NOT EXISTS orders(
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			amount INT,
			description TEXT
		);
	`)

	if err != nil {
		panic(err)
	}

	fmt.Println("Tables created!")

	/*
		Inserting records

		name := "zoz calhoun"
		email := "zoz@calhoun.io"

		// Insert some data to table
		// We use QueryRow to get the ID of newly created record. Postgres doesn't support LastInsertedID method
		row := db.QueryRow(`
			INSERT INTO users(name, email)
			VALUES(
				$1,
				$2
			) RETURNING id;
		`, name, email)

		// We don't need to check if row.Err() is nil since row.Scan() will return the error if it's not nil
		var id int
		err = row.Scan(&id)
		if err != nil {
			panic(err)
		}

		fmt.Println("User created", id)
	*/

	/*
		Querying Single Record

			id := 7
			var name, email string
			row := db.QueryRow(`
			  SELECT email,name
			  FROM users
			  WHERE id=$1;
			`, id)
			err = row.Scan(&name, &email)
			if err == sql.ErrNoRows {
				fmt.Println("No rows!")
			}
			if err != nil {
				panic(err)
			}
			fmt.Println(name, email)
	*/

	/*
			Creating Multiple Records

		userId := 1
		for i := 1; i <= 5; i++ {
			amount := i * 100
			description := fmt.Sprintf("Fake order #%d", i)
			_, err := db.Exec(`
				INSERT INTO orders(user_id, amount, description)
				VALUES(
					$1,
					$2,
					$3
				);
		`, userId, amount, description)
			if err != nil {
				panic(err)
			}
		}

		fmt.Println("Created fake orders!")
	*/

	type Order struct {
		ID          int
		UserID      int
		Amount      int
		Description string
	}

	var orders []Order

	userId := 1

	rows, err := db.Query(`
		SELECT id,amount, description FROM orders
		WHERE user_id=$1
	`, userId)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var order Order
		order.UserID = userId

		err := rows.Scan(&order.ID, &order.Amount, &order.Description)
		if err != nil {
			panic(err)
		}

		orders = append(orders, order)
	}

	// If any error is encountered during the iteration of 'rows', rows.Err() returns this error
	if err := rows.Err(); err != nil {
		panic(err)
	}

	fmt.Println(orders)
}
