package models

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	GalleryID int
	Path      string
	Filename  string
}

type Gallery struct {
	Id     int
	UserId int
	Title  string
}

type GalleryService struct {
	DB *sql.DB
	// ImagesDir is used to tell GalleryService where to store and locate images.
	// If not set, GalleryService will use "images" by default.
	// Make sure to add new custom ImagesDir to .gitignore file if you use custom ImagesDir
	ImagesDir string
}

func (service *GalleryService) Images(galleryId int) ([]Image, error) {
	globPattern := filepath.Join(service.galleryDir(galleryId), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving gallery images: %w", err)
	}

	var imagePaths []Image
	extensions := service.extensions()
	for _, file := range allFiles {
		if hasExtension(file, extensions) {
			imagePaths = append(imagePaths, Image{
				GalleryID: galleryId,
				Path:      file,
				Filename:  filepath.Base(file),
			})
		}
	}

	return imagePaths, nil
}

func (service *GalleryService) Image(galleryID int, filename string) (Image, error) {
	imagePath := filepath.Join(service.galleryDir(galleryID), filename)

	_, err := os.Stat(imagePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Image{}, ErrNotFound
		}
		return Image{}, fmt.Errorf("query image: %w", err)
	}

	return Image{
		Filename:  filename,
		GalleryID: galleryID,
		Path:      imagePath,
	}, nil
}

func hasExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}

func (service *GalleryService) extensions() []string {
	return []string{".png", ".jpg", ".jpeg", ".gif"}
}

func (service *GalleryService) galleryDir(id int) string {
	imagesDir := service.ImagesDir

	if imagesDir == "" {
		imagesDir = "images"
	}

	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
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
		return nil, fmt.Errorf("query galleries by user id: %w", err)
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

func (service *GalleryService) Delete(id int) error {
	_, err := service.DB.Exec(`
		DELETE FROM galleries
		WHERE id = $1;
	`, id)

	if err != nil {
		return fmt.Errorf("delete gallery: %w", err)
	}
	return nil
}
