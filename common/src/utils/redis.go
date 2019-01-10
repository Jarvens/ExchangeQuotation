// auth: kunlun
// date: 2019-01-10
// description: redis操作工具类
package utils

import (
	"config"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

//
func RedisSet(key, value string) bool {
	conn := config.RedisPool.Get()
	//释放资源
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	if err != nil {
		fmt.Println("set  error : ", err)
		return false
	}
	return true
}

func RedisSetWithExpire(key, value string, num int) bool {
	conn := config.RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	if err != nil {
		fmt.Println("set error: ", err)
		return false
	}
	v1, err1 := conn.Do("EXPIRE", key, num*60)
	fmt.Println(v1)
	if err1 != nil {
		fmt.Println("set with expire error: ", err1)
		return false
	}
	return true
}

func RedisGetString(key string) (value string, err error) {
	conn := config.RedisPool.Get()
	defer conn.Close()
	v, err := redis.String(conn.Do("GET", key))
	if err != nil {
		fmt.Printf("get %s error: %v", key, err)
		return "", err
	}
	return v, nil
}

func RedisGetStringMap(key string) (value map[string]string, err error) {
	conn := config.RedisPool.Get()
	defer conn.Close()
	v, err := redis.StringMap(conn.Do("GET", key))
	if err != nil {
		fmt.Println("get string map  error: ", err)
		return nil, err
	}
	return v, nil
}
