package models

import (
	"database/sql"
	"fmt"
)

type Gallery struct {
	Id     int
	UserId int
	Title  string
}

type GalleryService struct {
	DB *sql.DB
}

func (service *GalleryService) Create(userId int, title string) (*Gallery, error) {
	gallery := Gallery{
		UserId: userId,
		Title:  title,
	}

	row := service.DB.QueryRow(`
		INSERT INTO galleries (user_id, title)
		VALUES (
			$1,
			$2
		)
		RETURNING id;
	`, userId, title)

	err := row.Scan(&gallery.Id)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}

	return &gallery, nil
}
