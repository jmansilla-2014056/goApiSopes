package models

type  Mensaje struct {
	Api     string `json:"api"`
	Result  bool `json:"result"`
	}

type  Mensajes []Mensaje
