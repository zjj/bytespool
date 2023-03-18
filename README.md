Example
### 
```golang
package main

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/zjj/bytespool"
)

func main() {
	i := make([]byte, 1024*1024*2)
	n, err := rand.Read(i)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("generate %d bytes\n", n)

	// create a pool of max 100M, and each size is 4K
	p := bytespool.NewBytesPool(1024*1024*100, 1024*4)

	// get a max 10M bytes from pool
	bytes := p.NewBytes(1024 * 1024 * 10)
	n, err = bytes.Write(i)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %d bytes\n", n)

	read10bytes := make([]byte, 10)
	n, err = bytes.Read(read10bytes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("read %d bytes -> read10byte\n", n)

	read10bytes = make([]byte, 10)
	n, err = bytes.Read(read10bytes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("read %d bytes -> read10byte\n", n)

	rest, err := bytes.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("read %d bytes -> rest\n", len(rest))
}
```

### License
https://github.com/zjj/bytespool/blob/master/LICENSE
