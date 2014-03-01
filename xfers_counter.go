package bslc

type transfersCounter struct {
    counter chan int
    currentTransfers, totalTransfers int
}

// ActiveTransfers returns the current number of active transfers.
func (c *transfersCounter) ActiveTransfers() int {
   return c.currentTransfers
}

// TotalTransfers returns the total sum of transfers initiated.
func (c *transfersCounter) TotalTransfers() int {
   return c.totalTransfers
}

func (c *transfersCounter) transferStart() {
    c.counter <- 1
}

func (c *transfersCounter) transferEnd() {
    c.counter <- -1
}

func (c *transfersCounter) initializeCounter() {
    c.counter = make(chan int)
    go func() {
        for count := range c.counter {
            c.currentTransfers += count
            if count > 0 {
                c.totalTransfers += count
            }
        }
    }()
}
