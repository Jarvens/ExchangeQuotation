// auth: kunlun
// date: 2019-01-10
// description:
package codec

import "utils"

type BaseProto struct {
	Magic string
	Len   uint16
	Data  string
	Crc32 uint32
}

func (proto *BaseProto) Encoder() []byte {

	return nil
}

// 默认协议构造
func NewDefaultProto(message string) *BaseProto {
	crc32Val := utils.Crc32([]byte(message))
	return new(BaseProto{Magic: QuoteHeader, Len: uint16(len(message)), Data: message, Crc32: crc32Val})
}

//指定协议头构造
func NewProto(magic, message string) *BaseProto {
	crc32Val := utils.Crc32([]byte(message))
}
