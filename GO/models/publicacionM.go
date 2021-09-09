package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type  Publicacion struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Nombre     string 		`json:"nombre"`
	Comentario string 		`json:"comentario"`
	Fecha      string 		`json:"fecha"`
	Hashtags   []string 	`json:"hashtags"`
	Upvotes    int 			`json:"upvotes"`
	Downvotes  int 			`json:"downvotes"`
}

type  Publicaciones []Publicacion
