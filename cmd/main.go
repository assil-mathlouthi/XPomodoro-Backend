package main

import (
	"backend/cmd/server"
	"backend/config"
	"backend/database"
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

// 	@title 			XPomodoro Tracker API
// 	@version 		1.0
// 	@description 	A Gamified Pomodoro productivity tracker API.

//	@host			localhost:8000
//	@BasePath		/api/v1

// 	@securityDefinitions.apikey  ApiKeyAuth
// 	@in                          header
// 	@name                        Authorization
// 	@description                 Provide your JWT token in the format: "Bearer <token>"

func main() {
	// connect to the database
	db, err := database.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAdress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	initStorage(db)

	server := server.NewServer("127.0.0.1:8000", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("Database connection established successfully")
}
