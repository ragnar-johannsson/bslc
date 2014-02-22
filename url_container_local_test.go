package bslc

import (
    "strconv"
    "testing"
    "time"
)

func TestURLContainerLocal(t *testing.T) {
    nets := NewIPNetContainer([]string{ "127.0.0.0/8" })
    uc := NewLocalURLContainer(nets, []string{})
    tick := time.Tick(time.Millisecond * 250)

    uc.AddURL("http://127.0.0.1/1")
    uc.AddURL("http://127.0.0.1/1")

    <-tick
    if u, err := uc.NextURL(); u != "http://127.0.0.1/1" || err != nil {
        t.Error("Error calling NextURL()")
    }

    if _, err := uc.NextURL(); err == nil {
        t.Error("Expecting error, got none")
    }

    uc.AddURL("http://127.0.0.1/1")

    <-tick
    if _, err := uc.NextURL(); err == nil {
        t.Error("Expecting error, got none")
    }

    uc.AddURL("http://127.0.0.1/2")
    uc.AddURL("http://127.0.0.1/3")
    uc.AddURL("http://1.2.3.4/")

    <-tick
    if u, err := uc.NextURL(); u != "http://127.0.0.1/2" || err != nil {
        t.Error("Error calling NextURL()")
    }

    if u, err := uc.NextURL(); u != "http://127.0.0.1/3" || err != nil {
        t.Error("Error calling NextURL()")
    }

    if uc.Len() != 0 {
        t.Error("Expected empty container")
    }

    for i := 0; i < 100000; i++ {
        uc.AddURL("http://127.0.0.1/x/" + strconv.Itoa(i))
    }
    <-tick

    if uc.Len() != 100000 {
        t.Log(uc.Len())
        t.Error("Invalid container length")
    }

    for i := 0; i < 5000; i++ {
        if _, err := uc.NextURL(); err != nil {
            <-tick
            if _, err2 := uc.NextURL(); err2 != nil {
               t.Error("Error getting URL")
            }
        }
    }

    if uc.Len() != 95000 {
        t.Log(uc.Len())
        t.Error("Invalid container length")
    }

}
