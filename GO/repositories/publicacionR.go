package repositories

import (
	"ApiGo/conexion"
	"ApiGo/models"
	"context"
)

var coleccion = conexion.GetCollection("Publicaciones")
var ctx = context.Background()

func Create(publicacion models.Publicacion) error {

	var err error

	_, err = coleccion.InsertOne(ctx, publicacion)

	if err != nil {
		return err
	}

	return nil
}
