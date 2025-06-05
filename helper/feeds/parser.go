package feeds

import (
    "errors"
    "net/http"
    "time"
    "github.com/mmcdole/gofeed"
)

var httpClient = &http.Client{Timeout: 8 * time.Second}

func ParseRSS(url string, limit int) ([]*gofeed.Item, error) {
    fp := gofeed.NewParser()
    fp.Client = httpClient

    feed, err := fp.ParseURL(url)
    if err != nil {
        return nil, err
    }
    if limit <= 0 || limit > len(feed.Items) {
        limit = len(feed.Items)
    }
    if limit == 0 {
        return nil, errors.New("rss: feed contains no items")
    }
    return feed.Items[:limit], nil
}