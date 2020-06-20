# go-bare [![GoDoc](https://godoc.org/git.sr.ht/~sircmpwn/go-bare?status.svg)](https://godoc.org/git.sr.ht/~sircmpwn/go-bare) [![builds.sr.ht status](https://builds.sr.ht/~sircmpwn/go-bare.svg)](https://builds.sr.ht/~sircmpwn/go-bare?)

An implementation of the [BARE](https://git.sr.ht/~sircmpwn/bare) message format
for Go.

## Marshal usage

```go
package main

import (
    "fmt"
    "git.sr.ht/~sircmpwn/go-bare"
)

// type Coordinates {
//    x: int
//    y: int
//    z: int
//    q: optional<int>
// }
type Coordinates struct {
    X int
    Y int
    Z int
    Q *int
}

func main() {
    var coords Coordinates
    payload := []byte{
        0x01, 0x00, 0x00, 0x00,
        0x02, 0x00, 0x00, 0x00,
        0x03, 0x00, 0x00, 0x00,
        0x01, 0x04, 0x00, 0x00, 0x00,
    }
    err := bare.Unmarshal(payload, &coords)
    if err != nil {
        panic(err)
    }
    fmt.Printf("coords: %d, %d, %d (%d)",
        coords.X, coords.Y, coords.Z, *coords.Q)
}
```

## Code generation

TODO
