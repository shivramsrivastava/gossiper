package anonlib

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"../config"
)

type anonConnection struct {
	//*net.TCPConn
	//*net.TCPAddr
}

//mesg Format
// type single byte
// total bytes
// total frameids
// data
//ex

//type PacketType byte

const (
	TYPE_ACK byte = iota
	TYPE_DATA
	TYPE_BEAT
)

func GetDumyBase64FrameIds(framid string) string {
	return base64.StdEncoding.EncodeToString([]byte(framid))
}

func GetTheClientData(framid string, enable int) []byte {
	tempStr := GetDumyBase64FrameIds(framid) + ":" + strconv.FormatInt(int64(enable), 10)
	return []byte(tempStr)

}

func GenerateRandomData() []byte {

	localFrameId := []string{"REDIS", "CASANDRA", "MONGO-DB", "MSSQL", "IN-FLUXDB"}
	rand.Seed(int64(time.Now().Nanosecond()))

	str := localFrameId[rand.Intn(len(localFrameId))]

	return GetTheClientData(str, rand.Intn(1))

}

func GetNFrameIds(nFrame int) []byte {
	result := []byte{}
	for i := 0; i < nFrame; i++ {

		if i != 0 {
			//skip comma
			result = append(result, []byte(" ")...)
		}
		time.Sleep(time.Second * 1)
		result = append(result, GenerateRandomData()...)

	}

	return result
}

func GenerateMsg(msgType byte, data []byte, totalMsgCount int) []byte {

	//calculate the total byte lenght
	totalLength := int(len(data) + 4 + 4 + 1)
	//create a dummy byte array
	resultByteArray := make([]byte, totalLength)

	//encode msg type in the first byte
	resultByteArray[0] = msgType
	//encode the total msg length
	binary.BigEndian.PutUint32(resultByteArray[1:], uint32(totalLength))
	//encode the total frames id sent
	binary.BigEndian.PutUint32(resultByteArray[5:], uint32(totalMsgCount))
	//next start from 9th index
	resultByteArray = append(resultByteArray[9:], data...)
	return resultByteArray
}

func NewaAnonConnection(tcp *net.TCPConn) *anonConnection {

	return nil
	//&anonConnection{tcp}

}

func BindAndStartListener() {
	tcpAddr := &net.TCPAddr{

		IP:   net.ParseIP(config.BindAddr),
		Port: config.BindPort,
	}
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println(" BindAndStartListener", err)
	}
	HandleClientConnection(tcpListener)
}

//listen for incommng connection
func HandleClientConnection(tcpListener *net.TCPListener) {

	for {
		conn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println("Unable to accept connection")
		}
		go SendMsg(conn)
		go RecvMsg(conn)

	}

}

func SendMsg(conn *net.TCPConn) {
	var MaxFrameLen int = 10
	count := 1
	for {

		rand.Seed(int64(time.Now().Nanosecond()))
		dateBytes := GetNFrameIds(count)
		conn.Write(GenerateMsg(byte(rand.Intn(1)), dateBytes, count))
		count++
		if count >= MaxFrameLen {
			//reset count
			count = 1
		}
	}

}

func RecvMsg(conn *net.TCPConn) {

	localClientBuf := make([]byte, 4096)
	conn.Read(localClientBuf)
	msgType := localClientBuf[0:1]

	if msgType[0] == TYPE_BEAT {
		//fmt.Printf("its a one heart beat")
	}

}

func SimpleClient() {

	clientConnection, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 5050,
	})

	if err != nil {
		fmt.Println("something got screwed in cleint", err)
		panic("gone")
	}

	fmt.Println("SimpleClient called")
	for {
		for i := 0; i < 10; i++ {

			str := "Client Msg" + strconv.FormatInt(int64(i), 10)
			fmt.Println("sent from clinet ", str)
			clientConnection.Write([]byte(str))
			time.Sleep(2 * time.Second)

			var locaClientBuf []byte
			locaClientBuf = make([]byte, 4096)
			readin, err := clientConnection.Read(locaClientBuf)
			if err != nil {
				fmt.Println("Unable to read in the bytes on the client side", err)
				return
			}

			exacttoread := int(binary.BigEndian.Uint32(locaClientBuf[:4]))

			fmt.Println(exacttoread)

			fmt.Println("data from the server", exacttoread, readin, ":", string(locaClientBuf[4:exacttoread]))

		}

	}
}
