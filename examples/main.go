package main
//go:generate go run git.sr.ht/~sircmpwn/go-bare/cmd/gen -p main -s Time schema.bare schema.go

import (
	"log"
)

func main() {
	var dept = ACCOUNTING
	log.Printf("%s: %d", dept.String(), uint8(dept))
}
