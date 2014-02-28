package bslc

import (
    "net/url"
    "strings"
    pathUtils "path"
)

func isValidPath(path string) bool {
    return !strings.HasPrefix(path, "#") &&
        !strings.HasPrefix(path, "data:") &&
        !strings.HasPrefix(path, "mailto:") &&
        !strings.HasPrefix(path, "javascript:") &&
        path != ""
}

func completeURL(uri string, path string) string {
    if strings.HasPrefix(strings.ToLower(path), "http://") || strings.HasPrefix(strings.ToLower(path), "https://") {
        return path
    }

    if strings.HasPrefix(path, "/") {
        p, _ := url.Parse(uri)
        return strings.Join([]string{ p.Scheme, "://", p.Host, pathUtils.Clean(path) }, "")
    }

    p, _ := url.Parse(uri)
    uri = strings.Join([]string{ p.Scheme, "://", p.Host }, "")

    if p.Path != "" && strings.Contains(p.Path, "/") {
        path = strings.Join([]string{ p.Path[0:strings.LastIndex(p.Path,"/")], path }, "/")
    } else {
        path = strings.Join([]string{ "/", path }, "")
    }

    return strings.Join([]string{ uri, pathUtils.Clean(path) }, "")
}

func getHostFromURL(uri string) string {
    u, err := url.Parse(uri)
    if err != nil {
        return ""
    }

    if count := strings.Count(u.Host, ":"); count == 1 {
        return strings.Split(u.Host, ":")[0]
    } else if count > 1 {
        if strings.Contains(u.Host, "]") {
            return strings.TrimPrefix(strings.Split(u.Host, "]")[0], "[")
        }
    }

    return u.Host
}
