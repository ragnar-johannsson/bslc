package bslc

import (
    "net/http"
    "net/url"
    "testing"
    "time"
)

func TestResponseHandler(t *testing.T) {
    rh := responseHandler{}
    ch := make(chan *Content)

    rh.AddMimeTypes([]string{"text/plain", "text/html"}, ch)

    theUrl, _ := url.Parse("http://localhost/index.html")
    header := http.Header{}
    header.Set("Content-Type", "text/html")
    res := http.Response{
        Status: "OK",
        Header: header,
        Request: &http.Request{ URL: theUrl },
        Body: mockBody{},
    }

    go func() {
        rh.sendResponse(&res)
    }()

    if content := <-ch; content.Filename != "index.html" {
        t.Error("Wrong filename")
    }

    closed := false
    go func() {
        for _ = range ch {}
        closed = true
    }()

    rh.closeHandlers()
    <-time.After(time.Millisecond * 100)

    if !closed {
        t.Error("Content channel not closed")
    }
}
