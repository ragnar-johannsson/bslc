package bslc

import (
    "net/http"
    "sync"
)

type responseHandler struct {
    handlers map[string][]chan *Content
}

// AddMimeType registers channel ch to receive content of the specified mimeType.
func (c *responseHandler) AddMimeType(mimeType string, ch chan *Content) {
    if c.handlers == nil {
        c.handlers = make(map[string][]chan *Content)
    }
    c.handlers[mimeType] = append(c.handlers[mimeType], ch)
}

// AddMimeTypes registers channel ch to receive content of the specified mimeTypes.
func (c *responseHandler) AddMimeTypes(mimeTypes []string, ch chan *Content) {
    for _, mimeType := range mimeTypes {
        c.AddMimeType(mimeType, ch)
    }
}

func (c *responseHandler) sendResponse(res *http.Response) {
    done := make(chan bool)
    content := newContent(res, done)

    wg := sync.WaitGroup{}
    channels := c.handlers[content.ContentType]
    if len(channels) > 0 {
        content.readBody()
    }

    for _, ch := range channels {
        wg.Add(1)
        go func(channel chan *Content) {
            defer wg.Done()
            channel <- content
            <-done
        }(ch)
    }
    wg.Wait()
}

func (c *responseHandler) closeHandlers() {
    alreadyClosed := make(map[chan *Content]bool)
    for key, _ := range c.handlers {
        for _, ch := range c.handlers[key] {
            if _, ok := alreadyClosed[ch]; ok {
                continue
            }

            close(ch)
            alreadyClosed[ch] = true
        }

        delete(c.handlers, key)
    }
}
