package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"net/http"
)

type Product struct {
	Id      int
	Model   string
	Company string
	Price   int
}

var database *sql.DB

// функция добавления данных
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		model := r.FormValue("model")
		company := r.FormValue("company")
		price := r.FormValue("price")

		_, err = database.Exec("INSERT INTO products (model, company, price) VALUES ($1, $2, $3)",
			model, company, price)

		if err != nil {
			panic(err)
		}
		http.Redirect(w, r, "/", 301)
	} else {
		http.ServeFile(w, r, "./templates/create.html")
	}
}

func IndexHandler(w http.ResponseWriter, _ *http.Request) {

	rows, err := database.Query("SELECT * FROM products")
	if err != nil {
		panic(err)
	}
	defer func(rows *sql.Rows) {
		var err = rows.Close()
		if err != nil {

		}
	}(rows)

	var products []Product
	for rows.Next() {
		p := Product{}
		err := rows.Scan(&p.Id, &p.Model, &p.Company, &p.Price)
		if err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, p)
	}

	tmpl, _ := template.ParseFiles("./templates/user.html")
	err = tmpl.Execute(w, products)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}

func main() {
	db, err := sql.Open("postgres", "user=todo password=todo dbname=for_to_do sslmode=disable")
	if err != nil {
		panic(err)
	}
	database = db
	defer func(db *sql.DB) {
		var err = db.Close()
		if err != nil {

		}
	}(db)

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/create", CreateHandler)

	fmt.Println("Server is listening...")
	err = http.ListenAndServe(`:8181`, nil)
	if err != nil {
		return
	}
}
