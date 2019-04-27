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
var CatId uint
var NewsOneId uint
var NewsTwoId uint
var CatTitle string
var NewsOneTitle string
var NewsTwoTitle string
var CommentId uint
var Murl = "http://gorest.ga/api/contentelements"

type TestContentelements struct {
	Errors []core.ErrorMsg
	Data   contentelements.Contentelements
}

type TestContentelement struct {
	Errors []core.ErrorMsg
	Data   contentelements.Contentelement
}

type TestContentcomment struct {
	Errors []core.ErrorMsg
	Data   contentelements.Contentcomment
}

type TestContentcomments struct {
	Errors []core.ErrorMsg
	Data   contentelements.Contentcomments
}

type TestContenttags struct {
	Errors []core.ErrorMsg
	Data   contentelements.Contenttags
}

type TestUser struct {
	Errors []core.ErrorMsg
	Data   users.User
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

	u := readUserBody(resp, t)

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

	url := "http://gorest.ga/api/users/login"
	var userJson = `{"Email":"admin@admin.a", "Password":"adminpass"}`

	resp := doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	AdminToken = u.Data.Token

	return
}

func TestUserLogin(t *testing.T) {

	url := "http://gorest.ga/api/users/login"
	var userJson = `{"Email":"testuser@test.t", "Password":"testpass"}`

	resp := doRequest(url, "POST", userJson, "")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	UserToken = u.Data.Token

	return
}

func TestCreate(t *testing.T) {
	url := Murl
	CatTitle = fake.Title()
	NewsOneTitle = fake.Title()
	NewsTwoTitle = fake.Title()
	el := &contentelements.Contentelement{
		Urld:        fake.Word(),
		UserID:      1,
		Title:       CatTitle,
		Description: fake.ParagraphsN(1),
		Content:     fake.Paragraphs(),
		Meta_title:  fake.Title(),
		Meta_descr:  fake.ParagraphsN(1),
		Kind:        1,
		Status:      1,
		Tags:        "news",
	}

	uj, err := json.Marshal(el)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	resp := doRequest(url, "POST", string(uj), AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	CatId = u.Data.ID

	el = &contentelements.Contentelement{
		Urld:        fake.Word(),
		UserID:      1,
		Parent:      int(CatId),
		Title:       NewsOneTitle,
		Description: fake.ParagraphsN(1),
		Content:     fake.Paragraphs(),
		Meta_title:  fake.Title(),
		Meta_descr:  fake.ParagraphsN(1),
		Kind:        1,
		Status:      1,
		Tags:        "news,test",
	}

	uj, err = json.Marshal(el)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	resp = doRequest(url, "POST", string(uj), AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	NewsOneId = u.Data.ID

	el = &contentelements.Contentelement{
		Urld:        fake.Word(),
		UserID:      1,
		Parent:      int(CatId),
		Title:       NewsTwoTitle,
		Description: fake.ParagraphsN(1),
		Content:     fake.Paragraphs(),
		Meta_title:  fake.Title(),
		Meta_descr:  fake.ParagraphsN(1),
		Kind:        1,
		Status:      1,
		Tags:        "news,test,check",
	}

	uj, err = json.Marshal(el)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	resp = doRequest(url, "POST", string(uj), AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	NewsTwoId = u.Data.ID

	return
}

func TestUpdate(t *testing.T) {
	url := fmt.Sprintf("%s%s%d", Murl, "/", CatId)
	userJson := `{"Title":"` + fake.Title() + `"}`

	resp := doRequest(url, "PATCH", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	return
}

func TestGetAll(t *testing.T) {
	// get count
	url := Murl

	resp := doRequest(url, "GET", "", " ")

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readElementsBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	if len(u.Data) != 3 {
		t.Errorf("Wrong elements count: %d", len(u.Data))
	}

	//---------------

	utitle, _ := toUrlcode(NewsTwoTitle)

	url1 := Murl + "?title=" + utitle

	resp1 := doRequest(url1, "GET", "", " ")

	if resp1.StatusCode != 200 {
		t.Errorf("Success expected: %d%s", resp1.StatusCode, url1)
	}

	u1 := readElementsBody(resp1, t)

	if len(u1.Errors) != 0 {
		t.Fatal(u1.Errors)
	}

	if u1.Data[0].Title != NewsTwoTitle {
		t.Errorf("Wrong title search - : %s", u1.Data[0].Title)
	}

	//---------------

	url2 := Murl + "?parent=" + fmt.Sprintf("%d", CatId)

	resp2 := doRequest(url2, "GET", "", " ")

	if resp2.StatusCode != 200 {
		t.Errorf("Success expected: %d %s", resp2.StatusCode, url2)
	}

	u2 := readElementsBody(resp2, t)

	if len(u2.Errors) != 0 {
		t.Fatal(u2.Errors)
	}

	if len(u2.Data) != 2 {
		t.Errorf("Wrong childrens search: %d %s", len(u2.Data), url2)
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

	url = fmt.Sprintf("%s%s%d", Murl, "/", CatId)

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

//func TestDeleteComments(t *testing.T) {
//url := fmt.Sprintf("%s%s%d%s%d", Murl, "/", int(NewsOneId), "/comments/", int(CommentId))

//resp := doRequest(url, "DELETE", "", UserToken)

//if resp.StatusCode != 200 {
//t.Errorf("Success expected: %d", resp.StatusCode)
//}

//u := readCommentsBody(resp, t)

//if len(u.Errors) != 0 {
//t.Fatal(u.Errors)
//}

//return
//}

func TestGetTags(t *testing.T) {
	url := "http://gorest.ga/api/contenttags"

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

func TestDelete(t *testing.T) {
	url := fmt.Sprintf("%s%s%d", Murl, "/", 0)

	resp := doRequest(url, "DELETE", "", AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) == 0 {
		t.Fatal("wrong id validation dont work")
	}

	deleteElement(t, CatId)
	deleteElement(t, NewsOneId)
	deleteElement(t, NewsTwoId)

	return
}
