//go:build integration

package expense

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	err := godotenv.Load("../dev.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
}

func TestCreateExpense(t *testing.T) {

	body := bytes.NewBufferString(`{
		"title": "buy a new phone",
		"amount": 39000,
		"note": "buy a new phone",
		"tags": ["gadget", "shopping"]
	}`)

	var e Expense

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "buy a new phone", e.Title)
	assert.Equal(t, 39000, e.Amount)
	assert.Equal(t, "buy a new phone", e.Note)
	assert.Equal(t, 2, len(e.Tags))
}

func TestGetExpenseByID(t *testing.T) {
	c := seedExpense(t)

	var latest Expense
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(c.ID)), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, c.ID, latest.ID)
	assert.NotEmpty(t, latest.Title)
	assert.NotEmpty(t, latest.Amount)
	assert.NotEmpty(t, latest.Note)
	assert.NotEmpty(t, latest.Tags)
}

func TestUpdateExpenseByID(t *testing.T) {
	c := seedExpense(t)

	body := bytes.NewBufferString(`{
		"title": "apple smoothie",
    	"amount": 89,
    	"note": "no discount",
    	"tags": ["beverage"]
	}`)

	var latest Expense
	res := request(http.MethodPut, uri("expenses", strconv.Itoa(c.ID)), body)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "apple smoothie", latest.Title)
	assert.Equal(t, 89, latest.Amount)
	assert.Equal(t, "no discount", latest.Note)
	assert.Equal(t, 1, len(latest.Tags))
}

func TestGetAllExpenses(t *testing.T) {
	seedExpense(t)
	var es []Expense

	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&es)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(es), 0)
}

func seedExpense(t *testing.T) Expense {
	var c Expense
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)
	err := request(http.MethodPost, uri("expenses"), body).Decode(&c)
	if err != nil {
		t.Fatal("can't create expense:", err)
	}
	return c
}

func uri(paths ...string) string {
	host := "http://localhost" + os.Getenv("PORT")
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", os.Getenv("AUTH_SCHEME")+" "+os.Getenv("AUTH_KEY"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	return json.NewDecoder(r.Body).Decode(v)
}
