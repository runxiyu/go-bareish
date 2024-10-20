package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	bare "git.sr.ht/~runxiyu/go-bareish"
	"git.sr.ht/~runxiyu/go-bareish/example"
)

func main() {
	var person example.Person
	var err error
	for ; err == nil; err = bare.UnmarshalReader(os.Stdin, &person) {
		switch person := person.(type) {
		case *example.Customer:
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
		case *example.Employee:
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
		case *example.TerminatedEmployee:
			log.Println("Terminated employee (no data)")
		}
	}
	log.Printf("Error: %e", err)
}
