package anonlib

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sync"

	"../common"
)

type anonConnection struct {
	sync.Mutex
}

const (
	TYPE_ACK byte = iota
	TYPE_DATA
	TYPE_BEAT
)

//Some Global variables
var (
	BindAddr string
	BindPort string
)

// GetTheClientData builds the frame-id:bool value byte string
func BuildClientFrameData(framid string, enable bool) []byte {

	enableStr := "0"
	if enable == true {
		enableStr = "1"
	}
	tempStr := framid + ":" + enableStr
	return []byte(tempStr)

}

// BuildTheClientDataFromMap locks and read the map and build calls the BuildClientFrameData
func BuildTheClientDataFromMap(M map[string]bool) []byte {

	common.ToAnon.Lck.Lock()
	defer common.ToAnon.Lck.Unlock()

	result := []byte{}
	for key, value := range M {
		result = append(result, BuildClientFrameData(key, value)...)
	}

	return result
}

func GenerateDataMsg(msgType byte, data []byte, totalMsgCount int) []byte {

	//calculate the total byte lenght
	totalLength := int(len(data) + 4 + 4 + 1)
	//create a dummy byte array
	resultByteArray := make([]byte, totalLength)

	//encode msg type in the first byte
	resultByteArray[0] = msgType
	//encode the total msg length
	binary.BigEndian.PutUint32(resultByteArray[1:5], uint32(len(data)+(totalMsgCount*2)))
	//encode the total frames id sent
	binary.BigEndian.PutUint32(resultByteArray[5:9], uint32(totalMsgCount))
	//next start from 9th index
	copy(resultByteArray[9:], data)
	return resultByteArray
}

func NewaAnonConnection() *anonConnection {
	return &anonConnection{}
}

func BindAndStartListener() {
	//Bind to all the network inteface int the system on the same port
	tcpListener, err := net.Listen("tcp", ":"+BindPort)
	if err != nil {
		fmt.Println(" BindAndStartListener", err)
	}
	HandleClientConnection(tcpListener)
}

//listen for incommng connection
func HandleClientConnection(tcpListener net.Listener) {

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			fmt.Println("Unable to accept connection")
		}
		newAnonClientConn := NewaAnonConnection()
		go newAnonClientConn.SendMsg(conn)
		go newAnonClientConn.RecvMsg(conn)
	}
}

// SendMsg will wait on the Common.ToAnon
func (annonClient *anonConnection) SendMsg(conn net.Conn) {
	for {
		<-common.ToAnon.Ch
		result := BuildTheClientDataFromMap(common.ToAnon.M)
		result = GenerateDataMsg(TYPE_DATA, result, len(common.ToAnon.M))
		fmt.Println(result, "sending to cleint")
		annonClient.Lock()
		n, err := conn.Write(result)
		if err != nil {
			fmt.Println("Unable to send data to client", err)
			annonClient.Unlock()
			return
		}
		log.Printf("Sent %d bytes", n)
		annonClient.Unlock()
	}

}

func (annonClient *anonConnection) RecvMsg(conn net.Conn) {

	for {
		localClientBuf := make([]byte, 4096)
		//defer annonClient.Unlock()
		_, err := conn.Read(localClientBuf)
		if err != nil {
			fmt.Println("Unable to read from the client", err)
			return
		}
		msgType := localClientBuf[0:1]
		if msgType[0] == TYPE_BEAT {
			annonClient.Lock()
			annonClient.SendAck(conn)
			annonClient.Unlock()
		} else {
			fmt.Println("Unknown msg from client", msgType)
		}
	}

}

func (annonClient *anonConnection) SendAck(conn net.Conn) {
	_, err := conn.Write([]byte{TYPE_ACK})
	if err != nil {
		fmt.Println("Failed to send ack", err)
	}

}

func Run(inBindAddress string, inBindPort string) {
	BindAddr = inBindAddress
	BindPort = inBindPort

	go BindAndStartListener()
}
