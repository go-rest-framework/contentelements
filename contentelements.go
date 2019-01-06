package contentelements

import (
	"net/http"

	"github.com/go-rest-framework/core"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var App core.App

type Contentelements []Contentelement

type Contentelement struct {
	gorm.Model
	Urld        string
	UserID      int
	Parent      int
	Title       string
	Description string `gorm:"type:varchar(500)"`
	Content     string `gorm:"type:text"`
	Meta_title  string
	Meta_descr  string `gorm:"type:text"`
	Kind        int
	Status      int
	Comments    []Contentcomment
	Tags        []Contenttag
}

type Contentcomment struct {
	gorm.Model
	Comment          string `gorm:"type:varchar(500)"`
	UserID           int
	Parent           int
	ContentelementID int
}

type Contenttag struct {
	gorm.Model
	Name             string
	ContentelementID int
}

func Configure(a core.App) {
	App = a

	App.DB.AutoMigrate(&Contentelement{}, &Contentcomment{}, &Contenttag{})

	App.R.HandleFunc("/api/contentelements", actionGetAll).Methods("GET")
	App.R.HandleFunc("/api/contentelements/{id}", actionGetOne).Methods("GET")

	App.R.HandleFunc("/api/contentelements", App.Protect(actionCreate, []string{"admin"})).Methods("POST")
	App.R.HandleFunc("/api/contentelements/{id}", App.Protect(actionUpdate, []string{"admin"})).Methods("PATCH")
	App.R.HandleFunc("/api/contentelements/{id}", App.Protect(actionDelete, []string{"admin"})).Methods("DELETE")
}

func actionGetAll(w http.ResponseWriter, r *http.Request) {
	var (
		elements    Contentelements
		rsp         = core.Response{Data: &elements}
		all         = r.FormValue("all")
		id          = r.FormValue("id")
		title       = r.FormValue("title")
		description = r.FormValue("description")
		content     = r.FormValue("content")
		sort        = r.FormValue("sort")
		parent      = r.FormValue("parent")
		db          = App.DB
	)

	if all != "" {
		db = db.Where("id LIKE ?", "%"+all+"%")
		db = db.Or("title LIKE ?", "%"+all+"%")
		db = db.Or("description LIKE ?", "%"+all+"%")
		db = db.Or("content LIKE ?", "%"+all+"%")
	}

	if id != "" {
		db = db.Where("id LIKE ?", "%"+id+"%")
	}

	if title != "" {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}

	if description != "" {
		db = db.Where("role LIKE ?", "%"+description+"%")
	}

	if content != "" {
		db = db.Where("role LIKE ?", "%"+content+"%")
	}

	if parent != "" {
		db = db.Where("parent = ?", parent)
	}

	if sort != "" {
		db = db.Order(sort)
	}

	db.Preload("Comments").Preload("Tags").Find(&elements)

	rsp.Data = &elements

	w.Write(rsp.Make())
}

func actionGetOne(w http.ResponseWriter, r *http.Request) {
	var (
		element Contentelement
		rsp     = core.Response{Data: &element}
	)

	vars := mux.Vars(r)
	App.DB.Preload("Comments").Preload("Tags").First(&element, vars["id"])

	if element.ID == 0 {
		rsp.Errors.Add("ID", "Contentelement not found")
	} else {
		rsp.Data = &element
	}

	w.Write(rsp.Make())
}

func actionCreate(w http.ResponseWriter, r *http.Request) {
	var (
		element Contentelement
		rsp     = core.Response{Data: &element}
	)

	if rsp.IsJsonParseDone(r.Body) {
		if rsp.IsValidate() {
			App.DB.Create(&element)
		}
	}

	rsp.Data = &element

	w.Write(rsp.Make())
}

func actionUpdate(w http.ResponseWriter, r *http.Request) {
	var (
		data    Contentelement
		element Contentelement
		rsp     = core.Response{Data: &data}
	)

	if rsp.IsJsonParseDone(r.Body) {
		if rsp.IsValidate() {

			vars := mux.Vars(r)
			App.DB.First(&element, vars["id"])

			if element.ID == 0 {
				rsp.Errors.Add("ID", "Contentelement not found")
			} else {
				App.DB.Model(&element).Updates(data)
			}
		}
	}

	rsp.Data = &element

	w.Write(rsp.Make())
}

func actionDelete(w http.ResponseWriter, r *http.Request) {
	var (
		element Contentelement
		rsp     = core.Response{Data: &element}
	)

	vars := mux.Vars(r)
	App.DB.First(&element, vars["id"])

	if element.ID == 0 {
		rsp.Errors.Add("ID", "Contentelement not found")
	} else {
		if App.IsTest {
			App.DB.Unscoped().Delete(&element)
		} else {
			App.DB.Delete(&element)
		}
	}

	rsp.Data = &element

	w.Write(rsp.Make())
}
