// auth: kunlun
// date: 2019-01-23
// description:
package entity

type PushData struct {
	Cmd    string
	Symbol string
	Data   DepthItem
}

type DepthItem struct {
	Asks [][]float32
	Bids [][]float32
}
