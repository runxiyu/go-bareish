package example

import (
	"io/ioutil"
	"testing"

	"git.sr.ht/~sircmpwn/go-bare"
	"github.com/stretchr/testify/assert"
)

func BenchmarkMarshal(b *testing.B) {
	person, _ := makeCustomer(b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := bare.Marshal(person)
		if err != nil {
			panic(err)
		}
	}

}

func BenchmarkUnmarshal(b *testing.B) {
	_, buf := makeCustomer(b)
	var person Person
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := bare.Unmarshal(buf, &person)
		if err != nil {
			panic(err)
		}
	}

}

func makeCustomer(b *testing.B) (Person, []byte) {
	buf, err := ioutil.ReadFile("customer.bin")
	assert.Nil(b, err)

	b.SetBytes(int64(len(buf)))

	var person Person
	err = bare.Unmarshal(buf, &person)
	assert.Nil(b, err)

	return person, buf
}
