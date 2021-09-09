package services

import (
	m "sopes/apigo/models"
	p "sopes/apigo/repositories"
)

func Create(publicacion m.Publicacion) error {
	err := p.Create(publicacion)
	if err != nil {
		return err
	}
	return nil
}

