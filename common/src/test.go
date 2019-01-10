// auth: kunlun
// date: 2019-01-10
// description:
package main

import "config"

func main() {

	pool := config.InitRedisPool()
	conn := pool.Get()
	conn.Do("set", "pool", "test")
}
