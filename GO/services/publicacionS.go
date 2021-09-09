package services

import (
	m "ApiGo/models"
	p "ApiGo/repositories"
)

func Create(publicacion m.Publicacion) error {
	err := p.Create(publicacion)
	if err != nil {
		return err
	}
	return nil
}

