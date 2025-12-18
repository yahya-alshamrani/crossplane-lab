package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type PageData struct {
	Items       []string
	IsAvailable bool
}

// checkEnv ensures all required variables are set before the app starts
func checkEnv() {
	required := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"}
	var missing []string

	for _, v := range required {
		if os.Getenv(v) == "" {
			missing = append(missing, v)
		}
	}

	if len(missing) > 0 {
		// Print a clear message and exit with error code 1
		fmt.Printf("ERROR: Application failed to start.\n")
		fmt.Printf("The following environment variables are missing: %s\n", strings.Join(missing, ", "))
		os.Exit(1) 
	}
}

func getDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	data := PageData{IsAvailable: false}

	db, err := getDB()
	if err != nil {
		log.Printf("Database connection failed: %v", err)
	} else {
		defer db.Close()
		data.IsAvailable = true

		rows, err := db.Query("SELECT name FROM products LIMIT 10")
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var name string
				rows.Scan(&name)
				data.Items = append(data.Items, name)
			}
		}
	}

	tmpl.Execute(w, data)
}

func main() {
	// 1. Validate environment before doing anything else
	checkEnv()

	// 2. Start the server if validation passes
	http.HandleFunc("/", handler)
	fmt.Println("Server successfully started on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
