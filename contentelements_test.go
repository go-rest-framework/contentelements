package contentelements_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/go-rest-framework/contentelements"
	"github.com/go-rest-framework/core"
	"github.com/go-rest-framework/users"
	"github.com/icrowley/fake"
)

var AdminToken string
var UserToken string
var CatId1 uint
var CatId2 uint
var NewsOneId uint
var NewsTwoId uint
var PostOneId uint
var PostTwoId uint
var CatTitle string
var NewsOneTitle string
var NewsOneOneTitle string
var NewsOneTwoTitle string
var NewsTwoTitle string
var NewsTwoOneTitle string
var NewsTwoTwoTitle string

var OneTags string
var OneOneTags string
var OneTwoTags string
var TwoTags string
var TwoOneTags string
var TwoTwoTags string

var NewTags string

var CommentId uint
var Murl = "http://localhost/api/contentelements"

type TestContentelements struct {
	Errors []core.ErrorMsg                 `json:"errors"`
	Data   contentelements.Contentelements `json:"data"`
}

type TestContentelement struct {
	Errors []core.ErrorMsg                `json:"errors"`
	Data   contentelements.Contentelement `json:"data"`
}

type TestContentcomment struct {
	Errors []core.ErrorMsg                `json:"errors"`
	Data   contentelements.Contentcomment `json:"data"`
}

type TestContentcomments struct {
	Errors []core.ErrorMsg                 `json:"errors"`
	Data   contentelements.Contentcomments `json:"data"`
}

type TestContenttags struct {
	Errors []core.ErrorMsg             `json:"errors"`
	Data   contentelements.Contenttags `json:"data"`
}

type TestUser struct {
	Errors []core.ErrorMsg `json:"errors"`
	Data   users.User      `json:"data"`
}

func doRequest(url, proto, userJson, token string) *http.Response {
	reader := strings.NewReader(userJson)
	request, err := http.NewRequest(proto, url, reader)
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func readElementsBody(r *http.Response, t *testing.T) TestContentelements {
	var u TestContentelements
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(body), &u)
	return u
}

func readCommentBody(r *http.Response, t *testing.T) TestContentcomment {
	var u TestContentcomment
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(body), &u)
	return u
}

func readCommentsBody(r *http.Response, t *testing.T) TestContentcomments {
	var u TestContentcomments
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(body), &u)
	return u
}

func readTagsBody(r *http.Response, t *testing.T) TestContenttags {
	var u TestContenttags
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(body), &u)
	return u
}

func readElementBody(r *http.Response, t *testing.T) TestContentelement {
	var u TestContentelement
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(body), &u)
	defer r.Body.Close()
	return u
}

func readUserBody(r *http.Response, t *testing.T) TestUser {
	var u TestUser
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal([]byte(body), &u)
	defer r.Body.Close()
	return u
}

func deleteElement(t *testing.T, id uint) {
	url := fmt.Sprintf("%s%s%d", Murl, "/", id)

	resp := doRequest(url, "DELETE", "", AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readElementBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	return
}

func toUrlcode(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func TestAdminLogin(t *testing.T) {

	url := "http://localhost/api/users/login"
	var userJson = `{"email":"admin@admin.a", "password":"adminpass"}`

	resp := doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	AdminToken = u.Data.Token

	return
}

func TestUserLogin(t *testing.T) {

	url := "http://localhost/api/users/login"
	var userJson = `{"email":"testuser@test.t", "password":"testpass"}`

	resp := doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	UserToken = u.Data.Token

	return
}

func CreateOne(t *testing.T, parent int, title string, tags string) uint {
	url := Murl
	el := &contentelements.Contentelement{
		Urld:        fake.Word(),
		Parent:      parent,
		Title:       title,
		Description: fake.ParagraphsN(1),
		Content:     fake.Paragraphs(),
		Meta_title:  fake.Title(),
		Meta_descr:  fake.ParagraphsN(1),
		Kind:        "standart",
		Status:      "active",
		Tags:        tags,
	}

	uj, err := json.Marshal(el)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return 0
	}

	resp := doRequest(url, "POST", string(uj), AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	return u.Data.ID
}

func TestCreate(t *testing.T) {

	NewsOneTitle = fake.Title()
	NewsOneOneTitle = fake.Title()
	NewsOneTwoTitle = fake.Title()
	NewsTwoTitle = fake.Title()
	NewsTwoOneTitle = fake.Title()
	NewsTwoTwoTitle = fake.Title()

	OneTags = fake.Word()
	OneOneTags = fake.Word()
	OneTwoTags = fake.Word()
	TwoTags = fake.Word()
	TwoOneTags = fake.Word()
	TwoTwoTags = fake.Word()

	CatId1 = CreateOne(t, 0, NewsOneTitle, OneTags)
	NewsOneId = CreateOne(t, int(CatId1), NewsOneOneTitle, OneOneTags)
	NewsTwoId = CreateOne(t, int(CatId1), NewsOneTwoTitle, OneTwoTags)
	CatId2 = CreateOne(t, 0, NewsTwoTitle, TwoTags)
	PostOneId = CreateOne(t, int(CatId2), NewsTwoOneTitle, TwoOneTags)
	PostTwoId = CreateOne(t, int(CatId2), NewsTwoTwoTitle, TwoTwoTags)

	return
}

func TestUpdate(t *testing.T) {
	url := fmt.Sprintf("%s%s%d", Murl, "/", CatId1)

	NewTags = fake.Word()

	el := &contentelements.Contentelement{
		Urld:        fake.Word(),
		Parent:      0,
		UserID:      1,
		Title:       NewsOneTitle,
		Description: fake.ParagraphsN(1),
		Content:     fake.Paragraphs(),
		Meta_title:  fake.Title(),
		Meta_descr:  fake.ParagraphsN(1),
		Kind:        "standart",
		Status:      "active",
		Tags:        NewTags,
	}

	uj, err := json.Marshal(el)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	resp := doRequest(url, "PATCH", string(uj), AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	return
}

func GetOne(t *testing.T, url string) TestContentelements {
	resp := doRequest(url, "GET", "", " ")

	fmt.Println(url)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readElementsBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	return u
}

func TestGetAll(t *testing.T) {
	//find limit 1 offset 1 and get second element
	u := GetOne(t, Murl+"?limit=1&offset=1")

	if len(u.Data) != 1 {
		t.Errorf("Wrong limit and offset search count: %d, need 1", len(u.Data))
	}
	//find by title in all
	utitle, _ := toUrlcode(NewsTwoTitle)
	u = GetOne(t, Murl+"?all="+utitle)

	if u.Data[0].Title != NewsTwoTitle {
		t.Errorf("Wrong title search - : %s", u.Data[0].Title)
	}
	//find by tags in all
	u = GetOne(t, Murl+"?all="+OneOneTags)

	if len(u.Data) != 1 {
		t.Errorf("Wrong tag search count: %d, need 1", len(u.Data))
	}
	//find by parent and title
	utitle, _ = toUrlcode(NewsOneOneTitle)
	u = GetOne(t, Murl+"?parent="+fmt.Sprintf("%d", CatId1)+"&title="+utitle)

	if len(u.Data) != 1 {
		t.Errorf("Wrong parent and title search count: %d, need 1", len(u.Data))
	}
	//find by parent, title and status and get 0 elements
	utitle, _ = toUrlcode(NewsOneOneTitle)
	u = GetOne(t, Murl+"?parent="+fmt.Sprintf("%d", CatId1)+"&title="+utitle+"&status=draft")

	if len(u.Data) != 0 {
		t.Errorf("Wrong parent, title and status search count: %d, need 0", len(u.Data))
	}
	//set tree = 0 and get more elements in counts
	u = GetOne(t, Murl+"?tree=0")

	if len(u.Data) != 6 {
		t.Errorf("Wrong parent, title and status search count: %d, need 6", len(u.Data))
	}

	//sort by id
	u = GetOne(t, Murl+"?sort=id")

	if u.Data[0].ID != CatId1 {
		t.Errorf("Wrong sorting by id first id: %d, need %d", u.Data[0].ID, CatId1)
	}
	//sort by -id
	u = GetOne(t, Murl+"?sort=-id")

	if u.Data[0].ID != CatId2 {
		t.Errorf("Wrong sorting by -id first id: %d, need %d", u.Data[0].ID, CatId2)
	}
	//sort by title
	u = GetOne(t, Murl+"?sort=title")

	if u.Data[0].ID != CatId1 {
		t.Errorf("Wrong sorting by title first id: %d, need %d", u.Data[0].ID, CatId1)
	}
	//sort by -title
	u = GetOne(t, Murl+"?sort=-title")

	if u.Data[0].ID != CatId2 {
		t.Errorf("Wrong sorting by -title first id: %d, need %d", u.Data[0].ID, CatId2)
	}
	//sort by created_at
	u = GetOne(t, Murl+"?sort=created_at")

	if u.Data[0].ID != CatId1 {
		t.Errorf("Wrong sorting by created_at first id: %d, need %d", u.Data[0].ID, CatId1)
	}
	//sort by -created_at
	u = GetOne(t, Murl+"?sort=-created_at")

	if u.Data[0].ID != CatId2 {
		t.Errorf("Wrong sorting by -created_at first id: %d, need %d", u.Data[0].ID, CatId2)
	}

	return
}

func TestGetOne(t *testing.T) {
	url := Murl + "/0"
	resp := doRequest(url, "GET", "", " ")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readElementBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("element not found dont work")
	}

	url = fmt.Sprintf("%s%s%d", Murl, "/", CatId1)

	resp = doRequest(url, "GET", "", " ")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readElementBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	return
}

func TestAddComments(t *testing.T) {
	url := fmt.Sprintf("%s%s%d%s", Murl, "/", int(NewsOneId), "/comments")
	el := &contentelements.Contentcomment{
		Comment: fake.ParagraphsN(1),
		UserID:  2,
		Parent:  0,
	}

	uj, err := json.Marshal(el)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	resp := doRequest(url, "POST", string(uj), UserToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readCommentBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	CommentId = u.Data.ID

	el1 := &contentelements.Contentcomment{
		Comment: fake.ParagraphsN(1),
		UserID:  2,
		Parent:  int(CommentId),
	}

	uj1, err := json.Marshal(el1)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	resp1 := doRequest(url, "POST", string(uj1), UserToken)

	if resp1.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp1.StatusCode)
	}

	u1 := readCommentBody(resp1, t)

	if len(u1.Errors) != 0 {
		t.Fatal(u1.Errors)
	}

	el2 := &contentelements.Contentcomment{
		Comment: fake.ParagraphsN(1),
		UserID:  2,
		Parent:  int(CommentId),
	}

	uj2, err := json.Marshal(el2)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	resp2 := doRequest(url, "POST", string(uj2), UserToken)

	if resp2.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp2.StatusCode)
	}

	el3 := &contentelements.Contentcomment{
		Comment: fake.ParagraphsN(1),
		UserID:  2,
		Parent:  int(u1.Data.ID),
	}

	uj3, err := json.Marshal(el3)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	resp3 := doRequest(url, "POST", string(uj3), UserToken)

	if resp3.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp3.StatusCode)
	}

	return
}

func TestReadComments(t *testing.T) {
	url := fmt.Sprintf("%s%s%d%s", Murl, "/", int(NewsOneId), "/comments")

	resp := doRequest(url, "GET", "", "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readCommentsBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if len(u.Data) < 1 {
		t.Errorf("Wrong comments count: %d", len(u.Data))
	}

	return
}

func TestUpdateComments(t *testing.T) {
	url := fmt.Sprintf("%s%s%d%s%d", Murl, "/", int(NewsOneId), "/comments/", int(CommentId))
	el := &contentelements.Contentcomment{
		Comment: "testtest",
	}

	uj, err := json.Marshal(el)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	resp := doRequest(url, "PATCH", string(uj), UserToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readCommentBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if u.Data.Comment != "testtest" {
		t.Errorf("Wrong new comment: %s", u.Data.Comment)
	}

	return
}

func TestDeleteComments(t *testing.T) {
	url := fmt.Sprintf("%s%s%d%s%d", Murl, "/", int(NewsOneId), "/comments/", int(CommentId))

	resp := doRequest(url, "DELETE", "", UserToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readCommentsBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	return
}

func TestGetTags(t *testing.T) {
	url := "http://localhost/api/contenttags"

	resp := doRequest(url, "GET", "", "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readTagsBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if len(u.Data) < 1 {
		t.Errorf("Wrong comments count: %d", len(u.Data))
	}

	return
}

//func TestDelete(t *testing.T) {
//url := fmt.Sprintf("%s%s%d", Murl, "/", 0)

//resp := doRequest(url, "DELETE", "", AdminToken)

//if resp.StatusCode != 200 {
//t.Errorf("Success expected: %d", resp.StatusCode)
//}

//u := readUserBody(resp, t)

//if len(u.Errors) == 0 {
//t.Fatal("wrong id validation dont work")
//}

//deleteElement(t, CatId)
//deleteElement(t, NewsOneId)
//deleteElement(t, NewsTwoId)

//return
//}
