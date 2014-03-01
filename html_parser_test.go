package bslc

import (
    "bytes"
    "testing"
    "net/http"
    "net/url"
    "time"
)

func TestParseHTML(t *testing.T) {
    html := `<html>
                <body>
                    <a href="http://1.2.3.4/1">link</a>
                    <img src="http://2.3.4.5/image">
                </body>
            </html>`
    container := NewLocalURLContainer(NewIPNetContainer([]string{"0.0.0.0/0"}), []string{})
    buf := bytes.NewBuffer([]byte(html))
    body := mockBody{ readFunc: func(p []byte) (int, error) {
        return buf.Read(p)
    }}

    theUrl, _ := url.Parse("http://localhost/index.html")
    header := http.Header{}
    header.Set("Content-Type", "text/html; charset=utf-8")
    res := http.Response{
        Status: "OK",
        Header: header,
        Request: &http.Request{ URL: theUrl },
        Body: body,
    }

    content := newContent(&res, make(chan bool))
    content.readBody()

    parseHTML(content, container)
    <-time.After(time.Millisecond)

    if u, _ := container.NextURL(); u != "http://1.2.3.4/1" {
        t.Error("Invalid url received")
    }

    if u, _ := container.NextURL(); u != "http://2.3.4.5/image" {
        t.Error("Invalid url received")
    }
}
