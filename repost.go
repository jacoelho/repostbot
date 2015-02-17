package main

import (
//  "fmt"
//  "io/ioutil"
  "net/http"
//  "net/url"
//  "encoding/json"
  "regexp"
  "github.com/hashicorp/golang-lru"
  "log"
)

type WebhookResponse struct {
  Username string `json:"username"`
  Text     string `json:"text"`
}

var (
  urlRegex *regexp.Regexp
  urlCache *lru.Cache
)

func init() {
    urlRegex = regexp.MustCompile(`<(https?([^\||>]*))>`)
    urlCache, _  = lru.New(20)
}

func urlMatcher(s string) (urls []string) {
     matches := urlRegex.FindAllStringSubmatch(s, -1)
     
     for _, match := range matches {
         if match[1] != "" {
	    urls = append(urls, match[1])
	 }
     }

     return
}

func webHookHandler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  log.Println(r.Form)

  for key, _ := range r.Form {
    log.Println(key)
  }

  incomingText := r.FormValue("text")
  log.Printf("cenas: %d", len(incomingText))
  log.Printf("Handling incoming request: %s", incomingText)

  urls := urlMatcher(incomingText)
  
  for _, url := range urls {
    _, ok := urlCache.Get(url)
    if ok {
      log.Printf("cenas")
    } else {
      log.Printf("cenas2")
      urlCache.Add(url, "2")
    }
  }
}

func main() {
  http.HandleFunc("/", webHookHandler)
  http.ListenAndServe(":8080", nil)
}
