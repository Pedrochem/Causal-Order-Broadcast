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
	"bufio"
	"fmt"
	"os"

	. "./COB"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify at least one address:port!")
		fmt.Println("go run chat.go 127.0.0.1:5001  127.0.0.1:6001    ...")
		fmt.Println("go run chat.go 127.0.0.1:6001  127.0.0.1:5001    ...")
		fmt.Println("go run chat.go ...  127.0.0.1:6001  127.0.0.1:5001")
		return
	}

	// inicializacao do processo
	addresses := os.Args[1:]
	fmt.Println(addresses)

	cob := CasualOrderBroadcast_Module{
		Req: make(chan CasualOrderBroadcast_Req_Message),
		Ind: make(chan CasualOrderBroadcast_Ind_Message)}

	cob.Init(addresses[0], addresses)

	// enviador de broadcasts
	go func() {

		scanner := bufio.NewScanner(os.Stdin)
		var msg string

		for {
			if scanner.Scan() {
				msg = scanner.Text()
				msg += "§" + addresses[0]
			}
			req := CasualOrderBroadcast_Req_Message{
				Addresses: addresses[0:],
				Message:   msg}
			cob.Req <- req
		}
	}()

	// receptor de broadcasts
	go func() {
		for {
			in := <-cob.Ind

			// imprime a mensagem recebida na tela
			fmt.Printf("Message from %v: %v\n", in.From, in.Message)
			fmt.Printf("Clock: %v\n", in.Clock)
		}
	}()

	blq := make(chan int)
	<-blq
}
