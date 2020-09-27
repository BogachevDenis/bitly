package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/bitly/database"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/asaskevich/govalidator"
	log "github.com/sirupsen/logrus"
)

const (
	emptyLine 	string = "Empty long URL"
	routeinuse 	string = "This route already in use"
	errRoute 	string = "It is not a valid URL"
)

type config struct {
	ServerHost		string
	HTTPServerPort	string
	PgPort     		string
	PgUser     		string
	PgPass     		string
	PgBase     		string
}

type Url struct{
	LongUrl 		string 	`json:"longUrl"`
	ShortUrl 		string 	`json:"shortUrl"`
	ErrorMsg		string 	`json:"errorMsg"`
}

var cfg config

func init() {
	//Read config file
	file, err := os.Open("config.cfg")
	if err != nil{
		log.WithFields(log.Fields{
			"file" : "config.cfg",
			"error" : err,
		}).Fatal("File can`t be opened")
	}
	log.WithFields(log.Fields{
		"file" : "config.cfg",
	}).Info("Config file was opened")

	defer file.Close()
	stat, _ := file.Stat()
	readByte := make([]byte, stat.Size())
	file.Read(readByte)
	json.Unmarshal(readByte, &cfg)

	rand.Seed(time.Now().UnixNano())
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",http.FileServer(http.Dir("static/"))))
	r.HandleFunc("/",mainPage)
	r.HandleFunc("/create", CreateRoute).Methods("POST")
	r.HandleFunc("/{route}", Route).Methods("GET")
	log.WithFields(log.Fields{
		"port" : "8080",
	}).Info("Starting Server")
	http.ListenAndServe(":8080",r)
}

// Main page
func mainPage(w http.ResponseWriter, r *http.Request)  {
	tmpl := template.Must(template.ParseFiles("static/index.html"))
	tmpl.Execute(w, nil)
}

// POST processing, Create new route
func CreateRoute(w http.ResponseWriter, r *http.Request)  {
	url := new(Url)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"file" : r.Body,
			"error": err,
		}).Warning("Read request error")
	}
	err = json.Unmarshal(body, &url)
	if err != nil {
		log.WithFields(log.Fields{
			"file" : body,
			"error": err,
		}).Warning("Unmarshal error")
	}
	url.checkShortUrl()
	url.checkLongUrl()
	url.addData()
	savedUrl , _ := json.Marshal(url)
	w.WriteHeader(200)
	w.Write(savedUrl)
}


//Processing new route
func Route(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	route := vars["route"]
	database.Connect(cfg.PgUser, cfg.PgPass, cfg.PgBase)
	longURL := database.GetlongURL(route)
	if longURL == "" {
		w.WriteHeader(404)
		w.Write([]byte("Page not found"))
	} else {
		http.Redirect(w, r, fmt.Sprintf(longURL), http.StatusSeeOther)
		log.WithFields(log.Fields{
			"redirect to" : longURL,
		}).Info("Redirect OK")
	}
}
// Check short URL
func (u *Url) checkShortUrl() {
	if u.ErrorMsg != "" {
		return
	}
	if len(u.ShortUrl) < 1 {
		u.ShortUrl = RandStringRunes(5)
		return
	}
	u.ShortUrl = strings.TrimSpace(u.ShortUrl)
	index := strings.Index(u.ShortUrl, "/")
	if index != -1 {
		u.ShortUrl = u.ShortUrl[:index]
	}
}
// Check long URL
func (u *Url) checkLongUrl() {
	if u.LongUrl == "" {
		u.ErrorMsg = emptyLine
		return
	}
	u.LongUrl = strings.TrimSpace(u.LongUrl)
	if !govalidator.IsURL(u.LongUrl){
		u.ErrorMsg = errRoute
		return
	} else if !(strings.HasPrefix(u.LongUrl, "http://") || strings.HasPrefix(u.LongUrl, "https://")){
		u.LongUrl = strings.Join([]string{"http://", u.LongUrl}, "")
	}
}

// Send data to database
func (u *Url) addData() {
	if u.ErrorMsg != "" {
		return
	}
	database.Connect(cfg.PgUser, cfg.PgPass, cfg.PgBase)
	longUrl := database.GetlongURL(u.ShortUrl)
	if longUrl != "" {
		u.ErrorMsg = routeinuse
		return
	}
	err := database.InsertData(u.LongUrl, u.ShortUrl)
	if err == nil{
		log.WithFields(log.Fields{
			"data" : u,
		}).Info("InsertData to DB")
	}
}
// Generate random route
func RandStringRunes(n int) string {
	var Runes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	mas := make([]rune, n)
	for i := range mas {
		mas[i] = Runes[rand.Intn(len(Runes))]
	}
	return string(mas)
}

