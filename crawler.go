// Package bslc provides an IP bound crawler with channel based delivery of
// crawled content.
package bslc

import (
    "net/http"
    "time"
)

// Crawler does the heavy lifting of the actual crawling. URLs is the URLContainer
// used for bookkeeping and must be initialized. MaxConcurrentConnections is optional;
// the default is 5 concurrent transfers.
type Crawler struct {
    URLs URLContainer
    MaxConcurrentConnections int

    responseHandler
    transfersCounter
    isCrawling bool
}

// StartCrawling starts the crawling process.
func (c *Crawler) StartCrawling() {
    if c.MaxConcurrentConnections == 0 {
        c.MaxConcurrentConnections = 5
    }

    c.initializeCounter()
    c.registerHTMLHandler()
    c.isCrawling = true
    urls := make(chan string)

    // URL processors
    for i := 0; i < c.MaxConcurrentConnections; i++ {
        go func() {
            for uri := range urls {
                c.transferStart()
                c.processURL(uri)
                c.transferEnd()
            }
        }()
    }

    // URL dispatcher
    go func() {
        for c.isCrawling {
            uri, err := c.URLs.NextURL()
            if err != nil {
                <-time.After(time.Second * 2)
                if c.ActiveTransfers() != 0 || c.URLs.Len() != 0 {
                    continue
                } else {
                    break
                }
            }

            urls <- uri
        }

        close(urls)
        c.closeHandlers()
    }()
}

// StopCrawling stops the crawling process. Transfers in progress will be completed.
func (c *Crawler) StopCrawling() {
    c.isCrawling = false
}

func (c *Crawler) registerHTMLHandler() {
    ch := make(chan *Content)
    c.AddMimeType("text/html", ch)

    for i := 0; i < c.MaxConcurrentConnections; i++ {
        go func() {
            for content := range ch {
                parseHTML(content, c.URLs)
                content.Done <- true
            }
        }()
    }
}

func (c *Crawler) processURL(uri string) {
    res, err := http.Get(uri)
    if err != nil {
        return
    }
    defer res.Body.Close()

    c.sendResponse(res)
}
