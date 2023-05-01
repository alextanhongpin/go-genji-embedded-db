package main

import (
	"fmt"

	"github.com/genjidb/genji"
	"github.com/genjidb/genji/document"
	"github.com/genjidb/genji/types"
)

func main() {
	db, err := genji.Open("mydb")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Exec(`
		CREATE TABLE user (
			id INT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			address (
				city TEXT DEFAULT '?',
				zipcode TEXT
			),
			friends ARRAY
		)`,
	)
	if err != nil {
		panic(err)
	}
	type User struct {
		ID              uint
		Name            string
		TheAgeOfTheUser int64 `genji:"age"`
		Address         struct {
			City    string
			ZipCode string
		}
	}
	u := User{
		Name: "jessie",
	}
	u.Address.City = "Lyon"
	u.Address.ZipCode = "90001"

	err = db.Exec(`INSERT INTO user VALUES ?`, &u)
	if err != nil {
		if genji.IsAlreadyExistsError(err) {
			fmt.Println("exists")
		} else {
			panic(err)
		}
	}
	res, err := db.Query(`SELECT * FROM user WHERE age IS NULL`)
	if err != nil {
		panic(err)
	}
	defer res.Close()

	err = res.Iterate(func(d types.Document) error {
		var u User
		err = document.StructScan(d, &u)
		if err != nil {
			return err
		}

		fmt.Println("id", u.ID)
		fmt.Printf("query user: %#v\n", u)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
