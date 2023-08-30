package main

import (
	"avitosegments/database"
	"avitosegments/handler"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "<password>"
	dbname   = "Avito_test"
)

func main() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("pgx", psqlconn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	*handler.Api = database.API{DB: db}
	r := chi.NewRouter()
	r.Post("/create_segment", handler.CreateHandler)
	r.Post("/delete_segment", handler.DeleteHandler)
	r.Post("/change_segments", handler.ChangeHandler)
	r.Post("/get_segments", handler.GetHandler)
	http.ListenAndServe(":3000", r)
}
