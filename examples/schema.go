package main

import (
	"errors"

	"git.sr.ht/~sircmpwn/go-bare"
)

type PublicKey [128]byte

func (t *PublicKey) Decode(data []byte) error {
	return bare.Unmarshal(data, t)
}

func (t *PublicKey) Encode() ([]byte, error) {
	return bare.Marshal(t)
}

type Department uint8

const (
	ACCOUNTING Department = 0
	ADMINISTRATION = 1
	CUSTOMER_SERVICE = 2
	DEVELOPMENT = 3
	JSMITH = 99
)

func (t Department) String() string {
	switch (t) {
	case ACCOUNTING:
		return "ACCOUNTING"
	case ADMINISTRATION:
		return "ADMINISTRATION"
	case CUSTOMER_SERVICE:
		return "CUSTOMER_SERVICE"
	case DEVELOPMENT:
		return "DEVELOPMENT"
	case JSMITH:
		return "JSMITH"
	}
	panic(errors.New("Invalid Department value"))
}

type Customer struct {
	Name     string `bare:"name"`
	Email    string `bare:"email"`
	Address  Address `bare:"address"`
	Orders   []struct {
		OrderId  int64 `bare:"orderId"`
		Quantity int32 `bare:"quantity"`
	} `bare:"orders"`
	Metadata map[string][]byte `bare:"metadata"`
}

func (t *Customer) Decode(data []byte) error {
	return bare.Unmarshal(data, t)
}

func (t *Customer) Encode() ([]byte, error) {
	return bare.Marshal(t)
}

type Employee struct {
	Name       string `bare:"name"`
	Email      string `bare:"email"`
	Address    Address `bare:"address"`
	Department Department `bare:"department"`
	HireDate   Time `bare:"hireDate"`
	PublicKey  *PublicKey `bare:"publicKey"`
	Metadata   map[string][]byte `bare:"metadata"`
}

func (t *Employee) Decode(data []byte) error {
	return bare.Unmarshal(data, t)
}

func (t *Employee) Encode() ([]byte, error) {
	return bare.Marshal(t)
}

type Address struct {
	Address [4]string `bare:"address"`
	City    string `bare:"city"`
	State   string `bare:"state"`
	Country string `bare:"country"`
}

func (t *Address) Decode(data []byte) error {
	return bare.Unmarshal(data, t)
}

func (t *Address) Encode() ([]byte, error) {
	return bare.Marshal(t)
}
