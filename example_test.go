package bslc_test

import (
    "fmt"
    "github.com/ragnar-johannsson/bslc"
    "github.com/ragnar-johannsson/bslc/mimetypes"
    "sync"
)

func Example() {
    // Initialize URL container with IPnets filter and add seed URLs
    allowedNetworks := bslc.NewIPNetContainer([]string{ "127.0.0.0/8" })
    seedUrls := []string{ "http://127.0.0.1/" }
    urls := bslc.NewLocalURLContainer(allowedNetworks, seedUrls)

    // Initialize crawler
    crawler := bslc.Crawler{
        URLs: urls,
        MaxConcurrentConnections: 5,
    }

    // Register mimetype handler channel with crawler
    ch := make(chan *bslc.Content)
    crawler.AddMimeTypes(mimetypes.Audio, ch)

    // Start content handlers
    wg := sync.WaitGroup{}
    for i := 0; i < crawler.MaxConcurrentConnections; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for content := range ch {
                // Process received content
                fmt.Println("Received content from URL: ", content.URL.String())

                // Signal when done with content
                content.Done <- true
            }
        }()
    }

    // Start crawling and wait until done
    crawler.StartCrawling()
    wg.Wait()
}
