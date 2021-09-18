package conexion

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

)

var (
	usr = "sopes"
	pass = "HXZjzV6k7yvR2x1kAKC3FLw4M79j3SGmUOYAa3X2q6rPImMlMY8jagAAtZWZBtFLpB0jL2BE6VrKlmvFCTmL7g=="
	host = "sopes.mongo.cosmos.azure.com"
	port = "10255"
)


func GetCollection(collection string) *mongo.Collection{
	var cadena string = "mongodb://sopes:HXZjzV6k7yvR2x1kAKC3FLw4M79j3SGmUOYAa3X2q6rPImMlMY8jagAAtZWZBtFLpB0jL2BE6VrKlmvFCTmL7g==@sopes.mongo.cosmos.azure.com:10255/?ssl=true&replicaSet=globaldb&retrywrites=false&maxIdleTimeMS=120000&appName=@sopes@"

	client, err := mongo.NewClient(options.Client().ApplyURI(cadena))

	if err != nil {
		panic(err.Error())
	}

	ctx, _ := context.WithTimeout(context.Background(), 10* time.Second)
	err = client.Connect(ctx)

	if err != nil {
		panic(err.Error())
	}

	return client.Database("BaseS").Collection(collection)

}

func GetBase() *sql.DB {
	fmt.Println("Go MySQL")

	// Open up our database connection.
	// I've set up a database on my local machine using phpmyadmin.
	// The database is called testDb
	db, err := sql.Open("mysql", "admin:admin123@tcp(34.134.27.2:3306)/sopes1p1")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	return db

}