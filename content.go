package bslc

import (
    "bytes"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "path"
    "strings"
)

// Content represents the returned content from a crawled URL. The content's
// origin URL is saved in URL. ContentType holds the MIME type as specified by
// the remote server. Filename is the content's filename or an empty string if
// one cannot be determined. Done is channel that must be signalled after all 
// processing of the content is done.
type Content struct {
    URL url.URL
    ContentType string
    Filename string
    Done chan bool

    body io.ReadCloser
    data []byte
}

// Reader returns a new io.Reader for the crawled content.
func (c *Content) Reader() io.Reader {
    return bytes.NewBuffer(c.data)
}

func newContent(res *http.Response, ch chan bool) *Content {
    mimeType := res.Header.Get("Content-Type")
    if strings.Contains(mimeType, ";") {
        mimeType = strings.Split(mimeType, ";")[0]
    }

    _, filename := path.Split(res.Request.URL.Path)
    if cd := res.Header.Get("Content-Disposition"); cd != "" {
        for _, scd := range strings.Split(cd, "; ") {
            if strings.HasPrefix(scd, "filename=") {
                filename = strings.TrimPrefix(scd, "filename=")
                break
            }
        }
    }

    return &Content{
        URL: *res.Request.URL,
        ContentType: mimeType,
        Filename: filename,
        Done: ch,
        body: res.Body,
    }
}

func (c *Content) readBody() {
    c.data, _ = ioutil.ReadAll(c.body)
    c.body.Close()
}
