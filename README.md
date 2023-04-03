### bytespool
The main idea of this repository is to implement a memory pool mechanism to manage the memory usage of data blocks and avoid the impact of frequent memory allocation and deallocation operations on system performance. The specific implementation method is to divide the data into blocks of a specific size and store them in a linked list. The underlying management of the blocks uses sync.Pool and limits the maximum usage of the memory pool and the maximum memory of each requested Bytes. This way, whenever a new data block is needed, it can be obtained from the memory pool without the need for memory allocation operations every time.

In this repository, Bytes implements the Read and Write methods, which can be used to read and write data blocks. By using the Bytes approach, frequent memory allocation and deallocation operations can be avoided.


### Installation
```bash
go get github.com/zjj/bytespool
```
### Example
```golang
package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
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

	rest, err := ioutil.ReadAll(bytes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("read %d bytes -> rest\n", len(rest))
}
```
### More
If the first argument of `bytespool.NewBytesPool` is 0, there is no limit to the memory usage of the pool. 

If the first argument of `NewBytes` is 0, there is no limit to the maximum length of the returned Bytes, but this length is subject to the maximum memory limit set by the `bytespool.NewBytesPool`'s first argment if that is not 0.

The `Read` and `Write` operations of `Bytes` may block until there are idle blocks available to return, similar to blocking I/O operations.

### License
https://github.com/zjj/bytespool/blob/master/LICENSE
