// auth: kunlun
// date: 2019-01-10
// description:
package config

type Config struct {
	Redis    *Redis
	RabbitMQ *RabbitMQ
	Mysql    *Mysql
}

type Redis struct {
	Host     string
	Port     int
	Password string
	DataBase int
}

type RabbitMQ struct {
	Host     string
	Port     string
	Vhost    string
	UserName string
	Password string
	Exchange []*Exchange
}

type Mysql struct {
	UserName string
	Password string
	Url      string
}

type Exchange struct {
	Name    string
	Type    string
	Durable bool
	Queue   []*Queue
}

type Queue struct {
	Name    string
	Key     string
	Durable bool
}
