package main

import (
  "net/http"
  "encoding/json"
  "regexp"
  "log"
  "time"
  "os"
  "strconv"
  "github.com/hashicorp/golang-lru"
)

type WebhookResponse struct {
  Username string `json:"username"`
  Text     string `json:"text"`
}

const (
  CacheDefaultSize = 20
  PortDefault = "8080"
)

var (
  urlRegex *regexp.Regexp
  urlCache *lru.Cache
)

func init() {
  urlRegex = regexp.MustCompile(`<(https?([^\||>]*))(\||>)`)
}

func urlMatcher(s string) (urls []string) {
  matches := urlRegex.FindAllString(s, -1)

  for _, match := range matches {
    urls = append(urls, match[1:len(match)-1])
  }

  return
}

func webHookHandler(w http.ResponseWriter, r *http.Request) {
  _ = r.ParseForm()
  log.Printf("%s", r.Form)

  incomingText := r.PostFormValue("text")

  // avoid feedback
  if r.PostFormValue("user_id") == "RepostBOT" {
    return
  }
  
  log.Printf("Handling incoming request: %s", incomingText)

  urls := urlMatcher(incomingText)

  for _, url := range urls {
    log.Printf("Checking %s", url)

    raw, repost := urlCache.Get(url)
    if repost {
      t1, _  := time.Parse(time.RFC3339,raw.(string))
      log.Printf("[REPOST] Posted %s ago", time.Since(t1))

      var response WebhookResponse
      response.Username = "RepostBOT"
      response.Text = "http://i.imgur.com/YndMSe5.png"
      b, _ := json.Marshal(response)
      w.Write(b)
    }

    t := time.Now()
    urlCache.Add(url, t.Format(time.RFC3339))
  }
}

func main() {
  port := os.Getenv("PORT")
  if port == "" {
    port = PortDefault
  }
  
  cacheSize, err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
  if err != nil || cacheSize <= 0 {
    cacheSize = CacheDefaultSize
  }

  urlCache, _  = lru.New(cacheSize)

  http.HandleFunc("/", webHookHandler)
  err = http.ListenAndServe(":" + port , nil)
  if err != nil {
    log.Fatal(err)
  }
}
