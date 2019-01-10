// auth: kunlun
// date: 2019-01-10
// description:
package main

import (
	"config"
	"fmt"
	"utils"
)

func main() {

	config.InitRedisPool()
	//boll := utils.RedisSetWithExpire("pool", "message",20)
	//fmt.Println(boll)
	str := "{\"name\":\"张三\"}"
	utils.RedisSet("json", str)
	v, err := utils.RedisGetString("json")
	if err != nil {
		fmt.Println("error  ", err)
	}
	fmt.Println(v)
}
