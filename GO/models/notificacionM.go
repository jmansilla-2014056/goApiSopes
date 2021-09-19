package models

type  Notificacion struct {
	Api     string `json:"api"`
	Estado  string `json:"estado"`
	Numero  int    `json:"numero"`
}

type  Notificaciones []Notificacion