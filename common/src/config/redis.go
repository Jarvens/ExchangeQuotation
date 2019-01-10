// auth: kunlun
// date: 2019-01-10
// description:
package config

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/configor"
	"time"
)

var conn redis

func init() {
	fmt.Printf("initial redis time: %v", time.Now().Unix())
}
