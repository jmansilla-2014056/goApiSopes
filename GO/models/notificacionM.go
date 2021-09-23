package models

type  Notificacion struct {
	Python_cosmos   int `json:"python_cosmos"`
	Python_sql  	int `json:"python_sql"`
	Go_cosmos  		int `json:"go_cosmos"`
	Go_sql  		int `json:"go_sql"`
	Rust_cosmos 	int `json:"rust_cosmos"`
	Rust_sql		int `json:"rust_sql"`
}

type  Notificaciones []Notificacion