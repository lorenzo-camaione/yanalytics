package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

const (
	defaultPort = "7638"
	defaultHost = "http://yanalytics.domain"
)

// JsTemplate is used to render the js file
// it will have a unique user identifier
// and the host where the analytics call will be sent
type JsTemplate struct {
	UserID string
	Host   string
}

// AnalyticsRequest contains json encoded data
// sent from y.js
type AnalyticsRequest struct {
	URL         string `json:"u"`
	UserAgent   string `json:"a,omitempty"`
	Source      string `json:"s,omitempty"`
	Referrer    string `json:"r,omitempty"`
	WindowWidth *int   `json:"w,omitempty"`
	UserID      string `json:"x,omitempty"`
}

// config handles runtime configuration vars
type config struct {
	Port, Host string
}

func main() {
	config := parseConfig()
	http.HandleFunc("/y", TrackPageView)
	http.HandleFunc("/y.js", SendJavaScriptTracker)
	http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil)
	fmt.Printf("listening on port %s\n", config.Port)
}

func parseConfig() *config {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	host := os.Getenv("HOST")
	if host == "" {
		host = defaultHost
	}
	return &config{port, host}
}

// SendJavaScriptTracker sends a javascript tracker with a unique user id
// force the file to be permanently cached in the browser cache
func SendJavaScriptTracker(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("If-Modified-Since") != "" {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("Last-Modified", time.Unix(0, 0).Format("Mon, 1 Jan 2006 15:04:05 MST"))
	w.Header().Set("Cache-Control", "private")
	w.Header().Set("Content-Type", "application/javascript")
	jsTempl := template.Must(template.ParseFiles("./y.js"))
	jsTempl.Execute(w, JsTemplate{uuid.New().String(), defaultHost})
}

// TrackPageView receives a track request from the y.js file tracker
// it decodes and collects all information
func TrackPageView(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	decodedStr, err := base64.StdEncoding.DecodeString(string(buf))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var aRq AnalyticsRequest
	json.Unmarshal([]byte(decodedStr), &aRq)
	fmt.Printf("storing %#v\n", aRq)
	w.WriteHeader(http.StatusOK)
	return
}
