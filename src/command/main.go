// auth: kunlun
// date: 2019-01-14
// description:
package main

import (
	"net/http"
	"websocket"
)

func main() {

	//server.Start()

	http.HandleFunc("/", websocket.Handle)
	http.ListenAndServe("0.0.0.0:1234", nil)

}
