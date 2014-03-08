package main

import (
    "io"
    "github.com/ragnar-johannsson/bslc"
    "log"
    "os"
    "path"
    "strings"
    "sync"
)

// IP prefixes associated with Microsoft 
var prefixes = []string{
    "64.4.0.0/18",
    "65.52.0.0/14",
    "70.37.0.0/17",
    "70.37.128.0/18",
    "94.245.64.0/18",
    "111.221.64.0/18",
    "157.54.0.0/15",
    "157.56.0.0/14",
    "157.60.0.0/16",
    "168.61.0.0/16",
    "168.62.0.0/15",
}

// Seed URLs to use
var seedUrls = []string{
    "http://microsoft.com/",
}

// EXE MIME types
var exeTypes = []string{
     "application/dos-exe",
     "application/exe",
     "application/msdos-windows",
     "application/octet-stream",
     "application/x-msdownload",
     "application/x-exe",
     "application/x-winexe",
     "application/x-msdos-program",
     "vms/exe",
}

func main() {
    // Initialize crawler
    crawler := bslc.Crawler{
        URLs: bslc.NewLocalURLContainer(bslc.NewIPNetContainer(prefixes), seedUrls),
        MaxConcurrentConnections: 5,
    }

    // Register mimetype handler channel with crawler
    ch := make(chan *bslc.Content)
    crawler.AddMimeTypes(exeTypes, ch)

    // Start content handlers
    wg := sync.WaitGroup{}
    for i := 0; i < crawler.MaxConcurrentConnections; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for content := range ch {
                // Process received content
                if strings.HasPrefix(strings.ToLower(content.Filename), ".exe") {
                    saveFile("saved", content)
                }

                // Signal when done
                content.Done <- true
            }
        }()
    }

    // Start crawling and wait until done
    crawler.StartCrawling()
    wg.Wait()
}

func saveFile(prefix string, content *bslc.Content) {
    dir := path.Join(prefix, content.URL.Host, path.Dir(content.URL.Path))
    filename := path.Join(dir, content.Filename)

    os.MkdirAll(dir, 0755)
    file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
    if err != nil {
        log.Println("Unable to save file", filename, err.Error())
        return
    }
    defer file.Close()

    io.Copy(file, content.Reader())
}
