package bslc

import (
    "io"
    "net/http"
    "net/url"
    "testing"
)

type mockBody struct {
    readFunc func([]byte) (int, error)
    counterFunc func()
}

func (m mockBody) Read(p []byte) (n int, err error) {
    if m.readFunc != nil {
        return m.readFunc(p)
    }

    if m.counterFunc != nil {
        m.counterFunc()
    }

    return 0, io.EOF
}

func (m mockBody) Close() error {
    if m.counterFunc != nil {
        m.counterFunc()
    }

    return nil
}

func TestContent(t *testing.T) {
    done := make(chan bool)
    theUrl, _ := url.Parse("http://localhost/simple.txt")
    header := http.Header{}
    header.Set("Content-Type", "text/plain; charset=utf-8")
    res := http.Response{
        Status: "OK",
        Header: header,
        Request: &http.Request{ URL: theUrl },
    }

    content := newContent(&res, done)

    if content.ContentType != "text/plain" {
        t.Error("Wrong content-type")
    }

    if content.Filename != "simple.txt" {
        t.Error("Wrong filename")
    }

    header.Set("Content-Disposition", "attachment; filename=indeed.txt")
    res.Header = header

    content = newContent(&res, done)

    if content.Filename != "indeed.txt" {
        t.Error("Wrong filename")
    }

    counter := 0
    res.Body = mockBody{ counterFunc: func () { counter++ }}
    content = newContent(&res, done)
    content.readBody()

    if counter != 2 {
        t.Error("Error calling read or close")
    }
}
