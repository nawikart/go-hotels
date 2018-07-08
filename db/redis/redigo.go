package redis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	)


func ZAdd(k string, s int, v string) {
	c, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
	    fmt.Println("connect to redis err", err.Error())
	    return
	}
	defer c.Close()

	_, err = c.Do("zadd",k, s, v)  //写
	// fmt.Println(n)
}

// func ZScan(k string, cs int, p string) string{
// 	c, err := redis.Dial("tcp", "localhost:6379")
// 	if err != nil {
// 	    fmt.Println("connect to redis err", err.Error())
// 	    return "error"
// 	}
// 	// defer c.Close()

// 	res, _ := c.Do("zscan", k, cs, p)//读
// 	fmt.Println("res: ")
// 	fmt.Println(res)
// 	return k
// }