package models

import (
	"database/sql"
	"errors"
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

func (service *GalleryService) ById(id int) (*Gallery, error) {
	gallery := Gallery{
		Id: id,
	}

	row := service.DB.QueryRow(`
		SELECT user_id, title
		FROM galleries
		WHERE id = $1;
	`, gallery.Id)

	err := row.Scan(&gallery.UserId, &gallery.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query gallery by id: %w", err)
	}

	return &gallery, nil
}
