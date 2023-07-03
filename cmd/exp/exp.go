package main

import (
	"fmt"

	"github.com/AkifhanIlgaz/lenslocked/models"
)

func main() {
	galleryService := models.GalleryService{}
	fmt.Println(galleryService.Images(1))
}
