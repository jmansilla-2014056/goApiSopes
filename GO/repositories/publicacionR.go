package repositories

import (
	"context"
	"fmt"
	"math/rand"
	"sopes/apigo/conexion"
	"sopes/apigo/models"
	"strconv"
	"time"
)

var coleccion = conexion.GetCollection("Publicaciones")
var db = conexion.GetBase()
var ctx = context.Background()

func CreateMongo(publicacion models.Publicacion) error {

	var err error

	_, err = coleccion.InsertOne(ctx, publicacion)

	if err != nil {
		return err
	}

	return nil
}

func CreateSql(p models.Publicacion) error {

	var err error

	var t = time.Now().Unix()
	var id = n(int(t)) + n(rand.Intn(9)) + n(rand.Intn(9)) + n(rand.Intn(9))

	var query = "call InsertarPublicacion("+id+","+c(p.Nombre)+","+c(p.Comentario)+","+c(p.Fecha)+","+n(p.Upvotes)+","+n(p.Downvotes)+");"
	fmt.Println(query)
	insert, err := db.Query(query)

	if err != nil {
		return err
	}

	for i, s := range p.Hashtags {
		fmt.Println(i, s)
		var query2 = "INSERT INTO hashtag(id_publicacion, hashtag)VALUES("+id+","+c(s)+");"

		insert2, err2 := db.Query(query2)

		if err2 != nil{
			return err2
		}

		defer insert2.Close()
	}

	defer insert.Close()

	return nil
}

func c(x string) string{
	return "'" + x + "'"
}

func n(y int)  string{
	return strconv.Itoa(y)
}