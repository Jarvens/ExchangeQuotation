// auth: kunlun
// date: 2019-01-16
// description:
package handle

import (
	"encoding/json"
	"reflect"
	"request"
	"strings"
)

// value example:
func TickHandle(value string) {

}

/*

{
  "sub": "market.btcusdt.kline.1min",
  "id": "id1"
}
*/
func SubHandle(value string) {
	var slice []string = strings.Split(value, ".")
	//  market  btcusdt kline  1min
	if slice[0] == "" {

	}
}

func MessageHandle(data []byte) {
	var sub request.SubRequest
	err := Decoder(data, sub)
	if err != nil {
		//fmt.Println(err)
	}
	//fmt.Println("Print analysis Result: ", sub)
}

func Decoder(data []byte, result interface{}) (err error) {
	result = reflect.Interface
	err = json.Unmarshal(data, &result)
	return err
}
