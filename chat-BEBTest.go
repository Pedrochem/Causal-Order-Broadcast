// Construido como parte da disciplina de Sistemas Distribuidos
// PUCRS - Escola Politecnica
// Professor: Fernando Dotti  (www.inf.pucrs.br/~fldotti)

/*
LANCAR N PROCESSOS EM SHELL's DIFERENTES, PARA CADA PROCESSO, O SEU PROPRIO ENDERECO EE O PRIMEIRO DA LISTA
go run chat.go 127.0.0.1:5001  127.0.0.1:6001    ...
go run chat.go 127.0.0.1:6001  127.0.0.1:5001    ...
go run chat.go ...  127.0.0.1:6001  127.0.0.1:5001
*/

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	. "./COB"
)

func quit(messageLog []string) {
	for i := range messageLog {
		println(messageLog[i])
	}
	os.Exit(3)
}

func process(cob CasualOrderBroadcast_Module, address []string) {
	// envia broadcast
	go func() {
		println(cob.Ip, " STARTED")
		i := 0
		for {
			time.Sleep(2 * time.Second)
			msg := " Msg " + strconv.Itoa(i) + "§" + cob.Ip
			i++
			req := CasualOrderBroadcast_Req_Message{
				Addresses: address,
				Message:   msg}
			cob.Req <- req

		}

	}()

	// receptor de broadcasts
	go func() {
		for {
			// in := <-cob.Ind

			// imprime a mensagem recebida na tela
			// fmt.Println("-----------------------------------------------------")
			// fmt.Printf("Message from %v: %v\n", in.From, in.Message)
			// fmt.Printf("Clock: %v\n", in.Clock)
			// messageLog = append(messageLog, in.Message)
			// clockLog = append(clockLog, in.Clock)

		}
	}()

}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify at least one address:port!")
		fmt.Println("go run chat.go 127.0.0.1:5001  127.0.0.1:6001    ...")
		fmt.Println("go run chat.go 127.0.0.1:6001  127.0.0.1:5001    ...")
		fmt.Println("go run chat.go ...  127.0.0.1:6001  127.0.0.1:5001")
		return
	}

	addresses := os.Args[1:]
	// messageLog := make([]string, 10000)
	// clockLog := make([]map[string]int, 1000)

	cobs := make([]CasualOrderBroadcast_Module, len(addresses))

	println("Addresses size :", len(addresses))
	println("Cobs size :", len(cobs))

	for i, ip := range addresses {
		cob := CasualOrderBroadcast_Module{
			Req: make(chan CasualOrderBroadcast_Req_Message),
			Ind: make(chan CasualOrderBroadcast_Ind_Message)}
		cob.Init(ip, addresses)
		cobs[i] = cob
	}

	for _, c := range cobs {
		go process(c, addresses)
	}

	blq := make(chan int)
	<-blq

}
