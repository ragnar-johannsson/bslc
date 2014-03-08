package main

import (
    "bufio"
    "github.com/ragnar-johannsson/bslc"
    "github.com/ragnar-johannsson/bslc/mimetypes"
    "io"
    "log"
    "os"
    "path"
    "sync"
)

func main() {
    // Initialize URL container with an IP prefix filter and seed URLs
    allowedNetworks := bslc.NewIPNetContainer(readFile("is-prefixes.txt"))
    seedURLs := readFile("is-url_seeds.txt")
    urls := bslc.NewLocalURLContainer(allowedNetworks, seedURLs)

    // Initialize crawler
    crawler := bslc.Crawler{
        URLs: urls,
        MaxConcurrentConnections: 20,
    }

    // Register mimetype handler channel with crawler
    ch := make(chan *bslc.Content)
    crawler.AddMimeTypes(mimetypes.Images, ch)

    // Start content handlers
    wg := sync.WaitGroup{}
    for i := 0; i < crawler.MaxConcurrentConnections; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for content := range ch {
                // Process received content
                log.Println("Received URL: ", content.URL.String())
                saveFile("saved", content)
                content.Done <- true // Signal when done with content
            }
        }()
    }

    // Start crawling and wait until done
    crawler.StartCrawling()
    wg.Wait()
}

func saveFile(prefix string, content *bslc.Content) {
    dir, filename := path.Split(content.URL.Path)
    dir = path.Join(prefix, content.URL.Host, dir)
    if content.Filename != "" {
        filename = content.Filename
    }

    filename = path.Join(dir, filename)
    os.MkdirAll(dir, 0755)
    file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
    if err != nil {
        log.Println("Unable to save file", filename, err.Error())
        return
    }
    defer file.Close()

    io.Copy(file, content.Reader())

}

func readFile(filename string) []string {
    file, err := os.Open(filename)
    if err != nil {
        return []string{}
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    slice := []string{}
    for scanner.Scan() {
       slice = append(slice, scanner.Text())
    }

    return slice
}

