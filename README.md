# go-bare [![GoDoc](https://godoc.org/git.sr.ht/~sircmpwn/go-bare?status.svg)](https://godoc.org/git.sr.ht/~sircmpwn/go-bare) [![builds.sr.ht status](https://builds.sr.ht/~sircmpwn/go-bare.svg)](https://builds.sr.ht/~sircmpwn/go-bare?)

An implementation of the [BARE](https://git.sr.ht/~sircmpwn/bare) message format
for Go.

## Code generation

An example is provided in the `examples` directory. Here is a basic
introduction:

```
$ cat schema.bare
type Address {
	address: [4]string
	city: string
	state: string
	country: string
}
$ go run git.sr.ht/~sircmpwn/go-bare/cmd/gen -p models schema.bare models/gen.go
```

Then you can write something like the following:

```go
import "models"

/* ... */

bytes := []byte{ /* ... */ }
var addr Address
err := addr.Decode(bytes)
```

You can also add custom types and skip generating them by passing the `-s
TypeName` flag to gen, then providing your own implementation. For example, to
rig up time.Time with a custom "Time" BARE type, add this to your BARE schema:

```
type Time string
```

Then pass `-s Time` to gen, and provide your own implementation of Time in the
same package. See `examples/time.go` for an example of such an implementation.

## Marshal usage

For many use-cases, it may be more convenient to write your types manually and
use Marshal and Unmarshal directly. If you choose this approach, you may also
use `git.sr.ht/~sircmpwn/go-bare/schema.SchemaFor` to generate a BARE schema
langauge document describing your structs.

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

### Unions

To use union types, you need to define an interface to represent the union of
possible values, and this interface needs to implement `bare.Union`:

```go
type Person interface {
	Union
}
```

Then, for each possible union type, implement the interface:

```go
type Employee struct { /* ... */ }
func (e Employee) IsUnion() {}

type Customer struct { /* ... */ }
func (c Customer) IsUnion() {}
```

The IsUnion function is necessary to make the type compatible with the Union
interface. Then, to marshal and unmarshal using this union type, you need to
tell go-bare about your union:

```go
func init() {
    // The first argument is a pointer of the union interface, and the
    // subsequent arguments are values of each possible subtype, in ascending
    // order of union tag:
    bare.RegisterUnion((*Person)(nil), *new(Employee), *new(Customer))
}
```

This is all done for you if you use code generation.
