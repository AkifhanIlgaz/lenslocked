package controllers

import (
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/lenslocked/context"
	"github.com/AkifhanIlgaz/lenslocked/models"
)

type Galleries struct {
	Templates struct {
		New Template
	}
	GalleryService *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")

	g.Templates.New.Execute(w, r, data)
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserId int
		Title  string
	}
	data.Title = r.FormValue("title")
	data.UserId = context.User(r.Context()).ID

	gallery, err := g.GalleryService.Create(data.UserId, data.Title)
	if err != nil {
		g.Templates.New.Execute(w, r, data, err)
		return
	}

	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.Id)
	http.Redirect(w, r, editPath, http.StatusFound)
}
