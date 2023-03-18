package bytespool

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"testing"
)

func TestNewBytes0(t *testing.T) {
	i := make([]byte, 1024*1024*2)
	n, err := rand.Read(i)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("n:", n)

	p := NewBytesPool(1024*1024*100, 1024*4)
	b := p.NewBytes(1024 * 1024 * 100)
	n, err = b.Write(i)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("n:", n)

	bb, err := b.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	if !bytes.Equal(i, bb) {
		t.Fatal("not equal")
	} else {
		t.Log("TestNewBytesBuffer passed")
	}
}

func TestNewBytes1(t *testing.T) {
	i := make([]byte, 1024*1024*2)
	n, err := rand.Read(i)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("n:", n)

	p := NewBytesPool(1024*1024*100, 1024*4)
	b := p.NewBytes(1024 * 1024 * 100)
	n, err = b.Write(i)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("n:", n)

	bb := make([]byte, 10)
	n, err = b.Read(bb)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("n:", n)

	bb1, err := b.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	if !bytes.Equal(i[10:], bb1) {
		t.Fatal("not equal")
	}

	if !bytes.Equal(i, append(bb, bb1...)) {
		t.Fatal("not equal")
	}
}
