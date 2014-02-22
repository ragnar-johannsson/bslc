package bslc

// URLContainer is a container for URLs encountered during crawling. Call AddURL to
// add new URLs to the container. If the same URL has evern been added to the container,
// nothing is added. NextURL returns the next URL in line and removes it from the container, 
// or an empty string and an error if the container is empty. Len returns the number of URLs 
// in the container.
type URLContainer interface {
    AddURL(u string)
    NextURL() (u string, err error)
    Len() int
}
