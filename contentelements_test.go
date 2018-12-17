package contentelements_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/go-rest-framework/contentelements"
	"github.com/go-rest-framework/core"
	"github.com/go-rest-framework/users"
	"github.com/icrowley/fake"
)

var AdminToken string
var CatId uint
var NewsOneId uint
var NewsTwoId uint
var CatTitle string
var NewsOneTitle string
var NewsTwoTitle string
var Murl = "http://gorest.ga/api/contentelements"

type TestContentelements struct {
	Errors []core.ErrorMsg
	Data   contentelements.Contentelements
}

type TestContentelement struct {
	Errors []core.ErrorMsg
	Data   contentelements.Contentelement
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

func TestCreate(t *testing.T) {
	url := Murl
	CatTitle = fake.Title()
	NewsOneTitle = fake.Title()
	NewsTwoTitle = fake.Title()
	userJson := `{
		"Urld" : "` + fake.Word() + `",
		"UserID" : 1,
		"Parent" : "",
		"Title" : "` + CatTitle + `",
		"Description" : "` + fake.ParagraphsN(1) + `",
		"Content" : "` + fake.Paragraphs() + `",
		"Meta_title" : "` + fake.Title() + `",
		"Meta_descr" : "` + fake.ParagraphsN(1) + `",
		"Kind" : 1,
		"Status" : 1,
	}`

	resp := doRequest(url, "POST", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u := readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	CatId = u.Data.ID

	userJson = `{
		"Urld" : "` + fake.Word() + `",
		"UserID" : 1,
		"Parent" : ` + fmt.Sprintf("%d", CatId) + `,
		"Title" : "` + NewsOneTitle + `",
		"Description" : "` + fake.ParagraphsN(1) + `",
		"Content" : "` + fake.Paragraphs() + `",
		"Meta_title" : "` + fake.Title() + `",
		"Meta_descr" : "` + fake.ParagraphsN(1) + `",
		"Kind" : 1,
		"Status" : 1,
	}`

	resp = doRequest(url, "POST", userJson, AdminToken)

	if resp.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp.StatusCode)
	}

	u = readUserBody(resp, t)

	if len(u.Errors) != 0 {
		t.Fatal(u.Errors)
	}

	NewsOneId = u.Data.ID

	userJson = `{
		"Urld" : "` + fake.Word() + `",
		"UserID" : 1,
		"Parent" : ` + fmt.Sprintf("%d", CatId) + `,
		"Title" : "` + NewsTwoTitle + `",
		"Description" : "` + fake.ParagraphsN(1) + `",
		"Content" : "` + fake.Paragraphs() + `",
		"Meta_title" : "` + fake.Title() + `",
		"Meta_descr" : "` + fake.ParagraphsN(1) + `",
		"Kind" : 1,
		"Status" : 1,
	}`

	resp = doRequest(url, "POST", userJson, AdminToken)

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

	url1 := Murl + "?title=" + NewsTwoTitle

	resp1 := doRequest(url1, "GET", "", " ")

	if resp1.StatusCode != 200 {
		t.Errorf("Success expected: %d", resp1.StatusCode)
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
		t.Errorf("Success expected: %d", resp2.StatusCode)
	}

	u2 := readElementsBody(resp2, t)

	if len(u2.Errors) != 0 {
		t.Fatal(u2.Errors)
	}

	if len(u.Data) != 2 {
		t.Errorf("Wrong childrens search: %d", len(u.Data))
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
