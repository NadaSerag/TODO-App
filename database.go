package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error

	dsn := "postgres://postgres:Mydatabase123@localhost:5432/todo_app?sslmode=disable"

	//To run on docker, change path to:
	//dsn := "postgres://postgres:Mydatabase123@host.docker.internal:5432/todo_app?sslmode=disable"

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database")
	}

}
