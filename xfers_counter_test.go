package bslc

import (
    "testing"
    "time"
)

func TestTransfersCounter(t *testing.T) {
    counter := transfersCounter{}
    counter.initializeCounter()

    if counter.ActiveTransfers() != 0 {
        t.Error("Expected zero counter")
    }

    for i := 0; i < 50; i++ {
        go func() {
            counter.transferStart()
        }()
        go func() {
            counter.transferEnd()
            counter.transferStart()
        }()
        go func() {
            counter.transferStart()
        }()
        go func() {
            counter.transferEnd()
        }()
    }

    <-time.After(time.Millisecond * 100)

    if counter.ActiveTransfers() != 50 {
        t.Error("Expected 50 active transfers")
    }

    if counter.TotalTransfers() != 150 {
        t.Error("Expected 150 total transfers")
    }
}
