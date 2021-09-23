package models

type  Mensaje struct {
	Api     string `json:"api"`
	Cosmos  bool `json:"cosmos"`
	Sql  bool `json:"sql"`
	}

type  Mensajes []Mensaje
