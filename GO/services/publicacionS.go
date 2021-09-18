package services

import (
	m "sopes/apigo/models"
	p "sopes/apigo/repositories"
)

func CreateM(publicacion m.Publicacion) error {
	err := p.CreateMongo(publicacion)
	if err != nil {
		return err
	}
	return nil
}

func CreateS(publicacion m.Publicacion) error {
	err := p.CreateSql(publicacion)
	if err != nil {
		return err
	}
	return nil
}


