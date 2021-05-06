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
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	. "./COB"
)

func process(cob CasualOrderBroadcast_Module, address []string) {

	// envio de mensagens
	go func() {
		// indica que o processo foi incializado
		println(cob.Ip, " STARTED")
		i := 0
		for {

			// variacao de tempo para envio de mensagens
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

			// estrutura de mensagem a ser enviada
			// ip do processo substituido pela porta (ex: 127.0.0.1:5001 => 5001)
			msg := " Msg " + strconv.Itoa(i) + " " + strings.Split(cob.Ip, ":")[1] + "§" + cob.Ip
			i++
			req := CasualOrderBroadcast_Req_Message{
				Addresses: address,
				Message:   msg}
			cob.Req <- req

			// envio aleatorio de delay do processo 127.0.0.1:5001 ao 127.0.0.1:7001
			if cob.Ip == "127.0.0.1:5001" && rand.Intn(10) <= 7 {
				msg := " Msg " + strconv.Itoa(i) + " 5001 delay 7001 §" + cob.Ip
				i++
				req := CasualOrderBroadcast_Req_Message{
					Addresses: address,
					Message:   msg}
				cob.Req <- req
			}

		}

	}()

	// recebimento de mensagens
	go func() {
		// criacao de arquivo com os resultados de cada processo
		file, err := os.Create(strings.Split(cob.Ip, ":")[1] + "_output.txt")
		if err != nil {
			return
		}
		defer file.Close()

		// escrita no arquivo
		strAddresses := address[0]
		for i := 1; i < len(address); i++ {
			strAddresses += "|" + address[i]
		}
		file.WriteString(strAddresses + "\n")
		for {
			in := <-cob.Ind

			str := "-----------------------------------------------------\nDelivered " + cob.Ip + "\nMessage from " + in.From + ": " + in.Message + "\n"
			c := "Clock: map[ "
			for k, v := range in.Clock {
				c += strings.Split(k, ":")[1] + ":" + strconv.Itoa(v) + " "
			}
			c += "]\n"

			file.WriteString(str + c)
			if err != nil {
				panic(err)
			}

			// escrita no console
			fmt.Printf("-----------------------------------------------------\nDelivered %v\nMessage from %v: %v\nClock: %v\n", cob.Ip, in.From, in.Message, in.Clock)

		}
	}()

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify at least one address:port!")
		fmt.Println("go run app.go 127.0.0.1:5001  127.0.0.1:6001    ...")
		fmt.Println("go run app.go 127.0.0.1:6001  127.0.0.1:5001    ...")
		fmt.Println("go run app.go ...  127.0.0.1:6001  127.0.0.1:5001")
		return
	}

	addresses := os.Args[1:]

	// criacao de um array de processos
	cobs := make([]CasualOrderBroadcast_Module, len(addresses))

	for i, ip := range addresses {
		cob := CasualOrderBroadcast_Module{
			Req: make(chan CasualOrderBroadcast_Req_Message),
			Ind: make(chan CasualOrderBroadcast_Ind_Message)}
		cob.Init(ip, addresses)
		cobs[i] = cob
	}

	// inicializacao dos processos
	for _, c := range cobs {
		go process(c, addresses)
	}

	// bloqueio do programa
	blq := make(chan int)
	<-blq

}
