# Pool [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/lucklrj/pool)  [![Build Status](http://img.shields.io/travis/fatih/pool.svg?style=flat-square)](https://travis-ci.org/fatih/pool)


Pool is a very simple connection pool for connections interface. It can be used to
manage and reuse connections.


## Install and Usage

Install the package with:

```bash
go get github.com/lucklrj/pool
```

## Example

```go
package main

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"os"
	"time"
	 "github.com/lucklrj/pool"
)


func main() {
	d := make(chan int)
	pool := new(pool.Pool)
	pool.MaxOpenConns = 100
	pool.ConnMaxLifeTime = 2000
	pool.ConnTimeOut = 60
	//set client construct method
	pool.CreateClient = func() interface{} {
		return memcache.New("127.0.0.1:11211")
	}
	//set client destroy method
	pool.DestroyClient = func(c interface{}) {
		c.(*memcache.Client).Set(&memcache.Item{Key: "lrj", Value: []byte("my 333")})
	}
	err := pool.Init()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	//monitor clients in pool
	go func() {
		for {
			fmt.Print("pools num:")
			fmt.Println(len(pool.Pools))
			time.Sleep(time.Second)
		}
		d <- 1
	}()

	mc, err := pool.Get()
	if err != nil {
		fmt.Println(err.Error())
	}
	 mc.Client.(*memcache.Client).Set(&memcache.Item{Key: "test", Value: []byte("hello world")})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	//put used client into pool
	pool.Put(mc)
	pool.Close()
	fmt.Print("remain:")
	fmt.Println(mc.LeftTime)
	<-d
}
```

## Credits

 * [Lucklrj](https://github.com/lucklrj)

## License

The MIT License (MIT) - see LICENSE for more details


