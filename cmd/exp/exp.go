package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "fav-color", "purple")
	value := ctx.Value("fav-color")
	fmt.Println(value)
}

/*

		Connect to Postgres

	db, err := models.Open(models.DefaultPostgresConfig())
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

*/

/*
	Create random numbers by using time.Now().UnixNano() and math/rand package

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fmt.Println(r.Intn(100))
	fmt.Println(r.Intn(100))
	fmt.Println(r.Intn(100))
*/

/*
	Create random 32-byte string for session tokens

	n := 32
	b := make([]byte, n)
	fmt.Println(b)
	numberRead, err := rand.Read(b)

	if numberRead < n {
		panic("Didn't read enough random bytes")
	}

	if err != nil {
		panic(err)
	}

	fmt.Println(base64.URLEncoding.EncodeToString(b))
*/

/*

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


*/
