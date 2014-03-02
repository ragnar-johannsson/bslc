package bslc

import (
    "testing"
    "time"
)

func TestCrawler(t *testing.T) {
    crawler := Crawler{
        URLs: NewLocalURLContainer(NewIPNetContainer([]string{ "1.2.3.4/32"}), []string{}),
    }

    ch := make(chan *Content)
    crawler.AddMimeType("text/html", ch)

    closed := false
    go func() {
        for _ = range ch {}
        closed = true
    }()

    crawler.StartCrawling()
    <-time.After(time.Second * 3)

    if !closed {
        t.Error("Content channel not closed")
    }
}
