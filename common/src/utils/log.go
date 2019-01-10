// auth: kunlun
// date: 2019-01-10
// description:
package utils

import "github.com/alecthomas/log4go"

func Debug(arg0 interface{}, args ...interface{}) {
	log4go.Debug(arg0, args)
}

func Info(arg0 interface{}, args ...interface{}) {
	log4go.Info(arg0, args)
}

func Warn(arg0 interface{}, args ...interface{}) {
	log4go.Warn(arg0, args)
}

func Error(arg0 interface{}, args ...interface{}) {
	log4go.Error(arg0, args)
}

func init() {
	log4go.LoadConfiguration("log4go.xml")
}
