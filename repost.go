package main

import (
  "net/http"
//  "encoding/json"
  "regexp"
  "github.com/hashicorp/golang-lru"
  "log"
  "time"
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
   urlRegex = regexp.MustCompile(`<(https?([^\||>]*))(\||>)`)
   urlCache, _  = lru.New(20)
}

func urlMatcher(s string) (urls []string) {
     matches := urlRegex.FindAllString(s, -1)
     
     for _, match := range matches {
	urls = append(urls, match[1:len(match)-1])
     }

     return
}

//func urlCheck(s string) (url string) {
//  response, err := http.Get(url)
//}

func webHookHandler(w http.ResponseWriter, r *http.Request) {
  incomingText := r.PostFormValue("text")
  log.Printf("Handling incoming request: %s", incomingText)

  urls := urlMatcher(incomingText)
  
  for _, url := range urls {
    log.Printf("Checking %s", url)

    raw, repost := urlCache.Get(url)
    if repost {
      t1, _  := time.Parse(time.RFC3339,raw.(string)) 
      log.Printf("[REPOST] Posted %s ago", time.Since(t1))
    }
    
    t := time.Now()
    urlCache.Add(url, t.Format(time.RFC3339))
  }
}

func main() {
  http.HandleFunc("/", webHookHandler)
  http.ListenAndServe(":8080", nil)
}
