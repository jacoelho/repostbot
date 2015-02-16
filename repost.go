package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "net/url"
  "encoding/json"
  "regexp"
  "github.com/hashicorp/golang-lru"
  "log"
)

type WebhookResponse struct {
  Username string `json:"username"`
  Text     string `json:"text"`
}

type UrlWatcher struct {
  urlRegex *regexp.Regexp
  urlCache *lru.Cache
}

func NewUrlWatcher() (*UrlWatcher)
{

}

func init() {
    urlRegex = regexp.MustCompile(`<(https?([^\||>]*))>`)
    urlCache, _  = lru.New(20)
}

func webHookHandler(w http.ResponseWriter, r *http.Request) {
  incomingText := r.PostFormValue("text")
  log.Printf("Handling incoming request: %s", incomingText)

}

func main() {
  http.HandleFunc("/", webHookHandler)
  http.ListenAndServe(":8080", nil)
}
