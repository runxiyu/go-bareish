package main
//go:generate go run git.sr.ht/~sircmpwn/go-bare/cmd/gen -p main -s Time schema.bare schema.go

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"git.sr.ht/~sircmpwn/go-bare"
)

func main() {
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("read stdin: %e", err)
	}
	var person Person
	err = bare.Unmarshal(data, &person)
	if err != nil {
		log.Fatalf("decode: %e", err)
	}
	switch person := person.(type) {
	case *Customer:
		var addrs []string
		for _, addr := range person.Address.Address {
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
`, person.Name, person.Email, strings.Join(addrs, "\n"),
		person.Address.City, person.Address.State,
		person.Address.Country)
		for _, order := range person.Orders {
			fmt.Printf("- Order ID: %d\n  Quantity: %d\n",
				order.OrderId, order.Quantity)
		}
	case *Employee:
		var addrs []string
		for _, addr := range person.Address.Address {
			if addr != "" {
				addrs = append(addrs, addr)
			}
		}
		fmt.Printf(`Employee details:
Name: %s
Email: %s
Address:
	%s
	%s, %s
	%s
Department: %s
Hire date: %s
`, person.Name, person.Email, strings.Join(addrs, "\n"),
		person.Address.City, person.Address.State,
		person.Address.Country, person.Department.String(),
		time.Time(person.HireDate).Format(time.RFC3339))
	case *TerminatedEmployee:
		log.Println("Terminated employee (no data)")
	}
}
