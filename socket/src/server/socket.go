// auth: kunlun
// date: 2019-01-10
// description:
package server

import (
	"codec"
	"encoding/json"
	"net"
	"utils"
)

func Start() {
	address := "localhost:12345"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		utils.Debug("Analysis address error: %v", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		utils.Debug("Listener error: ", err)
		return
	}

	for {
		conn, err := listener.AcceptTCP()
		conn.SetNoDelay(true)
		if err != nil {
			utils.Debug("Accept error: %v", err)
			return
		}

		utils.Debug("Client address is: %s", conn.RemoteAddr())

		go loopHandler(conn)
	}

}

func loopHandler(conn net.Conn) {
	defer conn.Close()
	tmpBuffer := make([]byte, 0)
	//创建带缓冲的 chan
	ch := make(chan []byte, 16)
	go read(ch)
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			utils.Debug("Read message error: %v", err)
			return
		}
		tmpBuffer = codec.Decoder(buffer[:n], ch)
		response := codec.ResponseData{}
		_ = json.Unmarshal(tmpBuffer, &response)
		utils.Debug("Analysis success: %v", response)
	}
}

func read(ch chan []byte) {
	select {
	case data := <-ch:
		utils.Debug("Wtite data to chan: %v", data)
	}
}
