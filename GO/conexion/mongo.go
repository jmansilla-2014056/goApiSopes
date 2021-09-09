package conexion

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
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