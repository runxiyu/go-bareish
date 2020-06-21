package main
//go:generate go run git.sr.ht/~sircmpwn/go-bare/cmd/gen -p main -s Time schema.bare schema.go

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("read stdin: %e", err)
	}
	var customer Customer
	err = customer.Decode(data)
	if err != nil {
		log.Fatalf("decode: %e", err)
	}
	var addrs []string
	for _, addr := range customer.Address.Address {
		if addr != "" {
			addrs = append(addrs, addr)
		}
	}
	fmt.Printf(`Customer details:
Name: %s
Email: %s
Address:
	%s
	%s, %s
	%s
Orders:
`, customer.Name, customer.Email, strings.Join(addrs, "\n"),
	customer.Address.City, customer.Address.State,
	customer.Address.Country)
	for _, order := range customer.Orders {
		fmt.Printf("- Order ID: %d\n  Quantity: %d\n",
			order.OrderId, order.Quantity)
	}
}
