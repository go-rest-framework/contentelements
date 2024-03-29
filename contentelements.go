package contentelements

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-rest-framework/core"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var App core.App

type Contentelements []Contentelement
type Contentcomments []Contentcomment
type Contenttags []Contenttag

type Contentelement struct {
	gorm.Model
	Urld        string           `json:"urld" valid:"ascii,required"`
	UserID      int              `json:"userID"`
	Parent      int              `json:"parent"`
	Title       string           `json:"title" valid:"required"`
	Description string           `json:"description" gorm:"type:varchar(500)"`
	Content     string           `json:"content" gorm:"type:text"`
	Meta_title  string           `json:"meta_title"`
	Meta_descr  string           `json:"meta_descr" gorm:"type:text"`
	Kind        string           `json:"kind"`
	Status      string           `json:"status" valid:"required,in(active|suspend|draft)"`
	Tags        string           `json:"tags"`
	Elements    []Contentelement `json:"elements" gorm:"auto_preload;foreignkey:Parent"`
	Comments    []Contentcomment `json:"comments"`
}

type Contentcomment struct {
	gorm.Model
	Comment          string           `json:"comment" gorm:"type:varchar(500)"`
	UserID           int              `json:"userID"`
	Parent           int              `json:"parent"`
	ContentelementID int              `json:"contentelementID"`
	Comments         []Contentcomment `json:"comments" gorm:"auto_preload;foreignkey:Parent"`
}

type Contenttag struct {
	gorm.Model
	Name   string `json:"name"`
	Weight int    `json:"weight"`
}

type Parent struct {
	Id   uint
	Name string
}

type Parents []Parent

func Configure(a core.App) {
	App = a

	App.DB.AutoMigrate(&Contentelement{}, &Contentcomment{}, &Contenttag{})

	App.R.HandleFunc("/contentelements", actionGetAll).Methods("GET")
	App.R.HandleFunc("/contentelements/{id}", actionGetOne).Methods("GET")

	App.R.HandleFunc("/contentelements", App.Protect(actionCreate, []string{"admin"})).Methods("POST")
	App.R.HandleFunc("/contentelements/{id}", App.Protect(actionUpdate, []string{"admin"})).Methods("PATCH")
	App.R.HandleFunc("/contentelements/{id}", App.Protect(actionDelete, []string{"admin"})).Methods("DELETE")

	App.R.HandleFunc("/contentelements/{id}/comments", actionComments).Methods("GET")
	App.R.HandleFunc("/contentelements/{id}/comments", App.Protect(actionAddComment, []string{"user"})).Methods("POST")
	App.R.HandleFunc("/contentelements/{id}/comments/{cid}", App.Protect(actionUpdateComment, []string{"user"})).Methods("PATCH")
	App.R.HandleFunc("/contentelements/{id}/comments/{cid}", App.Protect(actionDeleteComment, []string{"user"})).Methods("DELETE")

	App.R.HandleFunc("/contenttags", actionTags).Methods("GET")
	App.R.HandleFunc("/parents", actionParents).Methods("GET")
}

func actionGetAll(w http.ResponseWriter, r *http.Request) {
	var (
		elements    Contentelements
		count       int
		rsp         = core.Response{Data: &elements, Req: r}
		all         = r.FormValue("all")
		id          = r.FormValue("id")
		title       = r.FormValue("title")
		description = r.FormValue("description")
		content     = r.FormValue("content")
		sort        = r.FormValue("sort")
		parent      = r.FormValue("parent")
		tree        = r.FormValue("tree")
		limit       = r.FormValue("limit")
		offset      = r.FormValue("offset")
		tags        = r.FormValue("tags")
		status      = r.FormValue("status")
		db          = App.DB
	)

	if all != "" {
		db = db.Where("id LIKE ?", "%"+all+"%")
		db = db.Or("title LIKE ?", "%"+all+"%")
		db = db.Or("description LIKE ?", "%"+all+"%")
		db = db.Or("tags LIKE ?", "%"+all+"%")
	}

	if id != "" {
		db = db.Where("id LIKE ?", "%"+id+"%")
	}

	if title != "" {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}

	if description != "" {
		db = db.Where("description LIKE ?", "%"+description+"%")
	}

	if content != "" {
		db = db.Where("content LIKE ?", "%"+content+"%")
	}

	if tags != "" {
		db = db.Where("tags LIKE ?", "%"+tags+"%")
	}

	if status != "" {
		db = db.Where("status = ?", status)
	}

	if parent != "" || parent == "0" {
		db = db.Where("parent = ?", parent)
	} else {
		if tree != "-1" {
			db = db.Where("parent = ?", 0)
		}
	}

	if tree == "" || tree == "1" {
		db = db.Set("gorm:auto_preload", true)
		db = db.Preload("Elements")
	}

	if sort != "" {
		switch sort {
		case "id":
			db = db.Order("id")
		case "-id":
			db = db.Order("id DESC")
		case "title":
			db = db.Order("title")
		case "-title":
			db = db.Order("title DESC")
		case "created_at":
			db = db.Order("created_at")
		case "-created_at":
			db = db.Order("created_at DESC")
		case "status":
			db = db.Order("status")
		case "-status":
			db = db.Order("status DESC")
		case "user":
			db = db.Order("user_id")
		case "-user":
			db = db.Order("user_id DESC")
		case "kind":
			db = db.Order("kind")
		case "-kind":
			db = db.Order("kind DESC")
		}
	} else {
		db = db.Order("id DESC")
	}

	db.Find(&elements).Count(&count)

	if limit != "" {
		db = db.Limit(limit)
	} else {
		db = db.Limit(5)
	}

	if offset != "" {
		db = db.Offset(offset)
	}

	db.Find(&elements)

	rsp.Data = &elements
	rsp.Count = count

	w.Write(rsp.Make())
}

func actionGetOne(w http.ResponseWriter, r *http.Request) {
	var (
		element Contentelement
		rsp     = core.Response{Data: &element, Req: r}
		db      = App.DB
	)

	vars := mux.Vars(r)

	db = db.Set("gorm:auto_preload", true)
	db = db.Preload("Elements")
	db = db.Preload("Comments")

	db.First(&element, vars["id"])

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
		rsp     = core.Response{Data: &element, Req: r}
	)

	if rsp.IsJsonParseDone(r.Body) && rsp.IsValidate() {
		i, err := strconv.Atoi(r.Header.Get("id"))
		if err != nil {
			rsp.Errors.Add("json", "User getting error"+err.Error())
		}
		element.UserID = i
		App.DB.Create(&element)
	}

	rsp.Data = &element

	s := strings.Split(element.Tags, ",")

	for _, v := range s {
		tag := Contenttag{
			Name:   v,
			Weight: 1,
		}
		App.DB.Where("name = ?", v).First(&tag)
		if tag.ID == 0 {
			App.DB.Create(&tag)
		} else {
			tag.Weight = tag.Weight + 1
			App.DB.Save(&tag)
		}
	}

	w.Write(rsp.Make())
}

func actionUpdate(w http.ResponseWriter, r *http.Request) {
	var (
		data    Contentelement
		element Contentelement
		rsp     = core.Response{Data: &data, Req: r}
	)

	if rsp.IsJsonParseDone(r.Body) {
		if rsp.IsValidate() {

			vars := mux.Vars(r)
			App.DB.First(&element, vars["id"])

			if element.ID == 0 {
				rsp.Errors.Add("ID", "Contentelement not found")
			} else {
				idstring := fmt.Sprintf("%d", element.UserID)
				if idstring != r.Header.Get("id") {
					rsp.Errors.Add("ID", "Only owner can change element")
				} else {
					App.DB.Model(&element).Updates(data)
				}
			}
		}
	}

	rsp.Data = &element

	w.Write(rsp.Make())
}

func actionDelete(w http.ResponseWriter, r *http.Request) {
	var (
		element Contentelement
		rsp     = core.Response{Data: &element, Req: r}
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

func actionComments(w http.ResponseWriter, r *http.Request) {
	var (
		comments Contentcomments
		rsp      = core.Response{Data: &comments, Req: r}
		limit    = r.FormValue("limit")
		offset   = r.FormValue("offset")
		db       = App.DB
	)

	vars := mux.Vars(r)
	db = db.Where("contentelement_id = ?", vars["id"])
	db = db.Where("parent = ?", 0)
	db = db.Set("gorm:auto_preload", true)

	if limit != "" {
		db = db.Limit(limit)
	}

	if offset != "" {
		db = db.Offset(offset)
	}

	db.Preload("Comments").Find(&comments)

	rsp.Data = &comments

	w.Write(rsp.Make())
}

func actionAddComment(w http.ResponseWriter, r *http.Request) {
	var (
		element Contentelement
		comment Contentcomment
		rsp     = core.Response{Data: &comment, Req: r}
		vars    = mux.Vars(r)
	)

	App.DB.First(&element, vars["id"])

	if element.ID == 0 {
		rsp.Errors.Add("ID", "Contentelement not found")
		return
	}

	if rsp.IsJsonParseDone(r.Body) {
		if rsp.IsValidate() {
			comment.ContentelementID = int(element.ID)
			App.DB.Create(&comment)
		}
	}

	rsp.Data = &comment

	w.Write(rsp.Make())
}

func actionUpdateComment(w http.ResponseWriter, r *http.Request) {
	var (
		data    Contentcomment
		comment Contentcomment
		rsp     = core.Response{Data: &data, Req: r}
	)

	if rsp.IsJsonParseDone(r.Body) {
		if rsp.IsValidate() {

			vars := mux.Vars(r)
			App.DB.First(&comment, vars["cid"])

			if comment.ID == 0 {
				rsp.Errors.Add("ID", "Comment not found")
			} else {
				idstring := fmt.Sprintf("%d", comment.UserID)
				if idstring != r.Header.Get("id") {
					rsp.Errors.Add("ID", "Only owner can change element")
				} else {
					App.DB.Model(&comment).Updates(data)
				}
			}
		}
	}

	rsp.Data = &comment

	w.Write(rsp.Make())
}

func actionDeleteComment(w http.ResponseWriter, r *http.Request) {
	var (
		comment Contentcomment
		rsp     = core.Response{Data: &comment, Req: r}
	)

	vars := mux.Vars(r)
	App.DB.First(&comment, vars["cid"])

	if comment.ID == 0 {
		rsp.Errors.Add("ID", "Contentcomment not found")
	} else {
		if App.IsTest {
			App.DB.Unscoped().Delete(&comment)
		} else {
			App.DB.Delete(&comment)
		}
	}

	rsp.Data = &comment

	w.Write(rsp.Make())
}

func actionTags(w http.ResponseWriter, r *http.Request) {
	var (
		tags Contenttags
		rsp  = core.Response{Data: &tags, Req: r}
		sort = r.FormValue("sort")
		db   = App.DB
	)

	if sort != "" {
		db = db.Order(sort)
	}

	db.Find(&tags)

	rsp.Data = &tags

	w.Write(rsp.Make())
}

func actionParents(w http.ResponseWriter, r *http.Request) {
	var (
		elements Contentelements
		res      Parents
		rsp      = core.Response{Data: &res, Req: r}
		db       = App.DB
	)

	db = db.Select("id, title")
	db = db.Where("status = ?", "active")
	db = db.Where("parent = ?", 0)
	db = db.Set("gorm:auto_preload", true)
	db = db.Preload("Elements")

	db.Find(&elements)

	for _, v := range elements {
		res = append(res, Parent{
			Id:   v.ID,
			Name: v.Title + fmt.Sprintf("%d", v.Parent),
		})

		res = genSubParents(res, v.Elements, 1)
	}

	fmt.Println(res)

	rsp.Data = &res

	w.Write(rsp.Make())
}

func genSubParents(list Parents, elements Contentelements, lvl int) Parents {
	for _, v := range elements {
		list = append(list, Parent{
			Id:   v.ID,
			Name: strings.Repeat("--", lvl) + " " + v.Title + fmt.Sprintf("%d", v.Parent),
		})
		list = genSubParents(list, v.Elements, lvl+1)
	}

	return list
}
