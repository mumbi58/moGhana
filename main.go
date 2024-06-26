package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
)

type payload struct {
	IncomingNumber string `json:"incomingNumber"`
	MessageText    string `json:"messageText"`
	ChannelNumber  string `json:"channelNumber"`
	Keyword        string `json:"keyword"`
}

func Process(c echo.Context, db *sql.DB) error {
	var p payload
	if err := c.Bind(&p); err != nil {
		return err
	}
	tableName := os.Getenv("TABLE_NAME")
	incomingNumber := p.IncomingNumber
	messageText := p.MessageText
	channelNumber := p.ChannelNumber

	// Use placeholders in the SQL query to avoid SQL injection
	query := fmt.Sprintf("INSERT INTO %v (message, sender_address, dest_address) VALUES (?, ?, ?)", tableName)

	_, err := db.Exec(query, messageText, incomingNumber, channelNumber)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}

	return c.String(http.StatusOK, "Payload processed successfully")
}

func main() {
	e := echo.New()

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Set up your MySQL database connection here
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	dbName := os.Getenv("DB_NAME")

	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", username, password, host, port, dbName)
	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v", err)
	}
	defer db.Close()

	e.POST("/moGhana", func(c echo.Context) error {
		return Process(c, db)
	})
	h := "0.0.0.0"
	p := os.Getenv("PORT")

	hst := fmt.Sprintf("%s:%s", h, p)
	e.Logger.Fatal(e.Start(hst))
}
