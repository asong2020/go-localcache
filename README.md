# go-localcache


A high-performance local cache implemented in the Go language.


## example

```go
package main

import (
	"fmt"
	cache "github.com/asong2020/go-localcache"
)

func main(){
    c, err := cache.NewCache()
	if err != nil {
	    return	
    }
	key := "asong"
	value := []byte("公众号：Golang梦工厂")
    err = c.Set(key, value)
	if err != nil {
	    return
    }
    entry, err := c.Get(key)
	if err != nil {
	    return
    }  
	fmt.Printf("get value is %s\n", string(entry))
	
	err = c.Delete(key)
	if err != nil{
		return
    }
}
```