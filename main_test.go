package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)



func TestCheckShortUrl1(t *testing.T)  {
	var url Url
	type testpair struct {
		value 	string
		want 	string
	}
	var tests = []testpair{
		{ "      test  ", "test" },
		{ "test/hrhr",	"test" },
		{ "test",			"test" },
	}
	for _, pair := range tests {
		url.ShortUrl = pair.value
		url.checkShortUrl()
		v := url.ShortUrl
		if v != pair.want {
			t.Error(
				"For", pair,
				"expected", pair.want,
				"got", v,
			)
		}
	}
}

func TestCheckShortUrl2(t *testing.T)  {
	var url Url
	type testpair struct {
		value 	string
		want 	int
	}
	var tests = []testpair{
		{ "test", 4 },
		{ "test_route", 10 },
	}
	for _, pair := range tests {
		url.ShortUrl = pair.value
		url.checkShortUrl()
		v := len(url.ShortUrl)
		if v != pair.want {
			t.Error(
				"For", pair,
				"expected", pair.want,
				"got", v,
			)
		}
	}
}




func TestCheckLongUrl1(t *testing.T)  {
	var url Url
	type testpair struct {
		value 	string
		want 	string
	}
	var tests = []testpair{
		{"", 						"Empty long URL" },
		{"hts://bitly.com/",		"It is not a valid URL" },
		{"https:/bitly.com/", 	"It is not a valid URL" },
		{"https://bitly.com/", 	"" },
	}
	for _, pair := range tests {
		url.ErrorMsg = ""
		url.LongUrl = pair.value
		url.checkLongUrl()
		v := url.ErrorMsg
		if v != pair.want {
			t.Error(
				"For", pair,
				"expected", pair.want,
				"got", v,
			)
		}
	}
}


func TestCheckLongUrl2(t *testing.T)  {
	var url Url
	type testpair struct {
		value 	string
		want 	string
	}
	var tests = []testpair{
		{"  https://bitly.com/",	"https://bitly.com/" },
		{"bitly.com/",			"http://bitly.com/" },
		{"  bitly.com/  ",		"http://bitly.com/" },
		{"https://bitly.com/",	"https://bitly.com/" },
	}
	for _, pair := range tests {
		url.LongUrl = pair.value
		url.checkLongUrl()
		v := url.LongUrl
		if v != pair.want {
			t.Error(
				"For", pair,
				"expected", pair.want,
				"got", v,
			)
		}
	}
}


func TestRandStringRunes(t *testing.T)  {
	type testpair struct {
		value 	int
		want 	int
	}
	var tests = []testpair{
		{3,	3},
		{5,	5},
	}
	for _, pair := range tests {
		v := len(RandStringRunes(pair.value))
		if v != pair.want {
			t.Error(
				"For", pair,
				"expected", pair.want,
				"got", v,
			)
		}
	}
}

func TestAddData(t *testing.T)  {
	var url Url
	type testpair struct {
		value 	string
		want 	string
	}
	var tests = []testpair{
		{"testroute","This route already in use" },
	}
	url.LongUrl = "https://bitly.com/"
	url.ShortUrl = "testroute"
	url.addData()
	for _, pair := range tests {
		url.ErrorMsg = ""
		url.ShortUrl = pair.value
		url.addData()
		v := url.ErrorMsg
		if v != pair.want {
			t.Error(
				"For", pair,
				"expected", pair.want,
				"got", v,
			)
		}
	}
}


func TestCreateRoute(t *testing.T)  {
	var url Url
	type testpair struct {
		value 		string
		wantStatus 	int
		wantData 	string
	}
	var tests = []testpair{
		{`{"longUrl":"https://bitly.com/"}`,	200,"" },
		{`{"longUrl":""}`,					200,"Empty long URL" },
		{`{"longUrl":"bitly.com/"}`,			200,"" },
		{`{"longUrl":"hts://bitly.com/"}`,	200,"It is not a valid URL" },
	}
	for _, pair := range tests {
		r := strings.NewReader(pair.value)
		req, err := http.NewRequest("POST", "/create", r)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateRoute)
		handler.ServeHTTP(rr, req)


		// Проверяем код
		if rr.Code != pair.wantStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, pair.wantStatus)
		}
		// Проверяем данные ответа
		body, _ := ioutil.ReadAll(rr.Body)
		json.Unmarshal(body, &url)
		if url.ErrorMsg != pair.wantData {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), pair.wantData)
		}

	}
}


func TestRoute(t *testing.T)  {
	type testpair struct {
		value 	string
		want 	int
	}
	var tests = []testpair{
		{RandStringRunes(7),		404},
	}
	for _, pair := range tests {
		req, err := http.NewRequest("GET", pair.value, nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Route)
		handler.ServeHTTP(rr, req)
		// Проверяем код
		if rr.Code != pair.want {
			t.Errorf("handler returned wrong status code: got %v want %v",
				rr.Code, pair.want)
		}


	}
}