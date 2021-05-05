package PP2PLink

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"
)

type PP2LinkMessage struct {
	Message string
	Clock   map[string]int
}

type PP2PLink_Req_Message struct {
	To   string
	Data PP2LinkMessage
}

type PP2PLink_Ind_Message struct {
	From string
	Data PP2LinkMessage
}

type PP2PLink struct {
	Ind   chan PP2PLink_Ind_Message
	Req   chan PP2PLink_Req_Message
	Run   bool
	Cache map[string]net.Conn
}

func (module *PP2PLink) Init(address string) {

	fmt.Println("[PP LINK] - Init PP2PLink!")
	if module.Run {
		return
	}

	module.Cache = make(map[string]net.Conn)
	module.Run = true
	module.Start(address)
}

func (module *PP2PLink) Start(address string) {
	go func() {
		listen, _ := net.Listen("tcp4", address)
		for {
			conn, err := listen.Accept()
			go func() {
				for {
					if err != nil {
						fmt.Println(err)
						continue
					}
					bufTam := make([]byte, 4)
					_, err := io.ReadFull(conn, bufTam)
					if err != nil {
						fmt.Println(err)
						continue
					}
					tam, err := strconv.Atoi(string(bufTam))
					bufMsg := make([]byte, tam)
					_, err = io.ReadFull(conn, bufMsg)
					if err != nil {
						fmt.Println(err)
						continue
					}
					bufMsgAsStruct := StringToMessage(string(bufMsg))
					cont := strings.Split(bufMsgAsStruct.Message, "§")
					bufMsgAsStruct.Message = cont[0]
					msg := PP2PLink_Ind_Message{
						From: cont[1],
						Data: bufMsgAsStruct}
					module.Ind <- msg
				}
			}()
		}
	}()

	go func() {
		for {
			message := <-module.Req
			go module.Send(message)
		}
	}()

}

func (module *PP2PLink) Send(message PP2PLink_Req_Message) {
	value := message.Data.Message
	if strings.Contains(value, "delay") {
		time.Sleep(3 * time.Second)
	}

	var conn net.Conn
	var ok bool
	var err error

	if conn, ok = module.Cache[message.To]; ok {

	} else {
		conn, err = net.Dial("tcp", message.To)
		if err != nil {
			fmt.Println(err)
			return
		}
		module.Cache[message.To] = conn
	}
	messageAsString := MessageToString(message.Data)
	strSize := strconv.Itoa(len(messageAsString))
	for len(strSize) < 4 {
		strSize = "0" + strSize
	}
	if !(len(strSize) == 4) {
		fmt.Println("ERROR AT PPLINK MESSAGE SIZE CALCULATION - INVALID MESSAGES MAY BE IN TRANSIT")
	}
	fmt.Fprintf(conn, strSize)
	fmt.Fprintf(conn, messageAsString)
}

func StringToMessage(s string) PP2LinkMessage {
	rawIn := json.RawMessage(s)
	bytes, err := rawIn.MarshalJSON()
	if err != nil {
		panic(err)
	}

	var message PP2LinkMessage
	err = json.Unmarshal(bytes, &message)
	if err != nil {
		panic(err)
	}
	return message
}

func MessageToString(message PP2LinkMessage) string {
	bytes, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
