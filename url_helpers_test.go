package bslc

import "testing"

func TestCompleteURL(t *testing.T) {
    var uri, path, out string

    uri = "http://test.uri"
    path = "http://full.uri/path"
    out = "http://full.uri/path"
    if completeURL(uri, path) != out { t.Fail() }

    uri = "http://test.uri/dir/file.ext"
    path = "/path.ext"
    out = "http://test.uri/path.ext"
    if completeURL(uri, path) != out { t.Fail() }

    uri = "http://test.uri/dir/file.ext"
    path = "path.ext"
    out = "http://test.uri/dir/path.ext"
    if completeURL(uri, path) != out { t.Fail() }

    uri = "http://test.uri"
    path = "main.html"
    out = "http://test.uri/main.html"
    if completeURL(uri, path) != out { t.Fail() }

    uri = "http://test.uri/dir1/dir2/index.html"
    path = ".././dir2/../../index.html"
    out = "http://test.uri/index.html"
    if completeURL(uri, path) != out { t.Fail() }
}

func TestGetHostFromURL(t *testing.T) {
    var uri, out string

    uri = "http://host.tld/smu"
    out = "host.tld"
    if getHostFromURL(uri) != out { t.Fail() }

    uri = "http://host:123/smu"
    out = "host"
    if getHostFromURL(uri) != out { t.Fail() }

    uri = "http://[FEDC:BA98:7654:3210:FEDC:BA98:7654:3210]:80/index.html"
    out = "FEDC:BA98:7654:3210:FEDC:BA98:7654:3210"
    if getHostFromURL(uri) != out { t.Fail() }

    uri = "http://[1080:0:0:0:8:800:200C:417A]/index.html"
    out = "1080:0:0:0:8:800:200C:417A"
    if getHostFromURL(uri) != out { t.Fail() }

    uri = "invalidurl"
    out = ""
    if getHostFromURL(uri) != out { t.Fail() }
}
