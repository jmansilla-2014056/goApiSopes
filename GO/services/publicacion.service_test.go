package services_test

import (
	m "sopes/apigo/models"
	ps "sopes/apigo/services"
	"testing"
)

func TestCreate(t *testing.T) {

	var publicacion = m.Publicacion{
		Nombre:     "nombrexxxx",
		Comentario: "comentario",
		Fecha:      "fecha",
		Hashtags:   []string{"a", "b", "c"},
		Upvotes:    2,
		Downvotes:  4,
	}
	err := ps.Create(publicacion)

	if err != nil{
		t.Error("La prueba fallo")
	}else {
		t.Log("La prueba fue un exito")
	}

}
