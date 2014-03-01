package bslc

import "code.google.com/p/go.net/html"

func parseHTML(content *Content, container URLContainer) {
    doc, err := html.Parse(content.Reader())
    if err != nil {
        return
    }

    var parse func(*html.Node)
    parse = func(n *html.Node) {
        if n.Type == html.ElementNode {
            switch n.Data {
            case "a", "link":
                for _, attr := range n.Attr {
                    if attr.Key == "href" {
                        if path := attr.Val; isValidPath(path) {
                            container.AddURL(completeURL(content.URL.String(), path))
                        }
                        break
                    }
                }
            case "embed", "input", "img", "script", "source":
                for _, attr := range n.Attr {
                    if attr.Key == "src" {
                        container.AddURL(completeURL(content.URL.String(), attr.Val))
                        break
                    }
                }
            }
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            parse(c)
        }
    }

    parse(doc)
}
