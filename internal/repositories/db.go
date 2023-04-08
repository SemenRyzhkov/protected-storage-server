package repositories

import (
	"database/sql"
	"log"
)

const (
	initUsersTableQuery = "" +
		"CREATE TABLE IF NOT EXISTS public.users (" +
		"id varchar(45) primary key, " +
		"login varchar(45) unique not null, " +
		"password varchar(45) not null" +
		")"
	initRawDataTableQuery = "" +
		"CREATE TABLE IF NOT EXISTS public.raw_data (" +
		"name varchar(45) unique not null, " +
		"data text, " +
		"user_id varchar(45) references public.users (id)" +
		")"
)

var db *sql.DB

func InitDB(dbAddress string) (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	db, connectionErr := sql.Open("postgres", dbAddress)
	if connectionErr != nil {
		log.Println(connectionErr)
		return nil, connectionErr
	}
	createTableErr := createTableIfNotExists(db)
	if createTableErr != nil {
		log.Println(createTableErr)
		return nil, createTableErr
	}
	return db, nil
}

func createTableIfNotExists(db *sql.DB) error {
	_, createUserTableErr := db.Exec(initUsersTableQuery)
	if createUserTableErr != nil {
		return createUserTableErr
	}

	_, createRawDataTableErr := db.Exec(initRawDataTableQuery)
	if createRawDataTableErr != nil {
		return createRawDataTableErr
	}
	return nil
}
