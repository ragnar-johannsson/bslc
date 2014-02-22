package bslc

import (
    "container/list"
    "errors"
)

type urlContainerLocal struct {
    allowedNetworks IPNetContainer
    urlMap map[string]bool
    urlLst *list.List
    in chan string
    out chan string
}

func (ucl *urlContainerLocal) AddURL(u string) {
    ucl.in <- u
}

func (ucl *urlContainerLocal) NextURL() (string, error) {
    defer func () {
        if len(ucl.out) < cap(ucl.out) / 2 && ucl.urlLst.Len() > 0 {
            ucl.in <- ucl.urlLst.Back().Value.(string)
        }
    }()

    if len(ucl.out) == 0 {
        return "", errors.New("No URLs")
    }

    if u, ok := <-ucl.out; !ok {
        return "", errors.New("No URLs")
    } else {
        return u, nil
    }
}

func (ucl *urlContainerLocal) Len() int {
    return len(ucl.out) + ucl.urlLst.Len()
}

func (ucl *urlContainerLocal) initialize() {
    go func() {
        for u := range ucl.in {
            if ucl.filter(u) {
                if len(ucl.out) < cap(ucl.out) {
                    ucl.out <- u
                } else {
                    ucl.urlLst.PushBack(u)
                }
            }

            if ucl.urlLst.Len() > 0 && len(ucl.out) < cap(ucl.out) {
                outLen := cap(ucl.out) - len(ucl.out)
                if outLen > ucl.urlLst.Len() {
                    outLen = ucl.urlLst.Len()
                }

                for i := 0; i < outLen; i++ {
                    ucl.out <-ucl.urlLst.Remove(ucl.urlLst.Front()).(string)
                }
            }
        }
    }()
}

func (ucl *urlContainerLocal) filter(u string) bool {
    if !ucl.allowedNetworks.Contains(getHostFromURL(u)) {
        return false
    }

    if _, exists := ucl.urlMap[u]; exists {
        return false
    }

    ucl.urlMap[u] = true
    return true
}

// NewLocalURLContainer returns a new URLContainer, stored locally in memory without
// any persistence between sessions and without the ability to share state between
// remote crawler instances. The allowedNetworks IP container handles the IP filtering
// and seedUrls are the URLs too bootstrap with.
func NewLocalURLContainer(allowedNetworks IPNetContainer, seedUrls []string) URLContainer {
    ucl := &urlContainerLocal{
        allowedNetworks: allowedNetworks,
        urlMap: make(map[string]bool, 10000),
        urlLst: list.New(),
        in: make(chan string, 10000),
        out: make(chan string, 1000),
    }
    ucl.initialize()

    for _, seedUrl := range seedUrls {
        ucl.AddURL(seedUrl)
    }

    return ucl
}
