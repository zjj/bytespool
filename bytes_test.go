package bytespool

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	bb, err := ioutil.ReadAll(b)
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

	bb1, err := ioutil.ReadAll(b)
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

func TestJSONEncoder(t *testing.T) {
	p := NewBytesPool(1024*100, 1024*4)
	b := p.NewBytes(1024 * 100)
	enc := json.NewEncoder(b)

	i := make([]byte, 1024*2)
	_, err := rand.Read(i)
	if err != nil {
		log.Fatal(err)
	}

	err = enc.Encode(map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
		"key3": i,
	})
	if err != nil {
		log.Fatal(err)
	}
	/*
		x, err := b.ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(x))
	*/

	h, err := ioutil.ReadAll(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("hi", string(h))
}
