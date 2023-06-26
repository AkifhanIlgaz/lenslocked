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

func (service *GalleryService) ByUserId(userId int) ([]Gallery, error) {
	rows, err := service.DB.Query(`
		SELECT id, title
		FROM galleries
		WHERE user_id = $1;
	`, userId)

	if err != nil {
		return nil, fmt.Errorf("query galleries by user id: %w", err)
	}

	var galleries []Gallery

	for rows.Next() {
		gallery := Gallery{
			UserId: userId,
		}
		err := rows.Scan(&gallery.Id, &gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("query galleries by user id: %w", err)
		}
		galleries = append(galleries, gallery)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("query galleris by user id: %w", err)
	}

	return galleries, nil
}

func (service *GalleryService) Update(gallery *Gallery) error {
	_, err := service.DB.Exec(`
		UPDATE galleries
		SET title = $2
		WHERE id = $1;
	`, gallery.Id, gallery.Title)

	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}
	return nil
}
