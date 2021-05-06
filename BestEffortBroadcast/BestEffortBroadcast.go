package BestEffortBroadcast

/*
  Construido como parte da disciplina de Sistemas Distribuidos
  Semestre 2018/2  -  PUCRS - Escola Politecnica
  Estudantes:  Andre Antonitsch e Rafael Copstein
  Professor: Fernando Dotti  (www.inf.pucrs.br/~fldotti)
  Algoritmo baseado no livro:
  Introduction to Reliable and Secure Distributed Programming
  Christian Cachin, Rachid Gerraoui, Luis Rodrigues

  Para uso vide ao final do arquivo, ou aplicacao chat.go que usa este
*/

import (
	"fmt"
	"strconv"
	"strings"

	PP2PLink "../PP2PLink"
)

// estrutura request
type BestEffortBroadcast_Req_Message struct {
	Addresses []string
	Message   string
	Clock     map[string]int
}

// estrutura ind
type BestEffortBroadcast_Ind_Message struct {
	From    string
	Message string
	Clock   map[string]int
}

// estrutura modulo
type BestEffortBroadcast_Module struct {
	Ind      chan BestEffortBroadcast_Ind_Message
	Req      chan BestEffortBroadcast_Req_Message
	Pp2plink PP2PLink.PP2PLink
}

// incializa modulo
func (module BestEffortBroadcast_Module) Init(address string) {
	fmt.Println("Init BEB!")
	module.Pp2plink = PP2PLink.PP2PLink{
		Req: make(chan PP2PLink.PP2PLink_Req_Message),
		Ind: make(chan PP2PLink.PP2PLink_Ind_Message),
	}
	module.Pp2plink.Init(address)
	module.Start()
}

func (module BestEffortBroadcast_Module) Start() {

	go func() {
		for {
			select {
			// caso modulo receba request da camada superior (COB) => realizar broadcast
			case y := <-module.Req:
				module.Broadcast(y)

				// caso modulo receba ind da camada inferior (PP2Link) => entregar mensagem
			case y := <-module.Pp2plink.Ind:
				module.Deliver(y)
			}
		}
	}()

}

func (module BestEffortBroadcast_Module) Broadcast(message BestEffortBroadcast_Req_Message) {
	// enviar a mensaegem para cada um dos enderecos no request
	for i := 0; i < len(message.Addresses); i++ {
		msg := BEB2PP2PLink(message, message.Addresses[i])
		module.Pp2plink.Req <- msg
	}
}

func (module BestEffortBroadcast_Module) Deliver(message PP2PLink.PP2PLink_Ind_Message) {
	// fazer deliver da mensagem acionanando a camada superior (COB)
	msg := PP2PLink2BEB(message)
	module.Ind <- msg
}

// conversor de Req BEB => Req PP2Link
func BEB2PP2PLink(req BestEffortBroadcast_Req_Message, to string) PP2PLink.PP2PLink_Req_Message {

	// encodificacao do relogio vetorial na mensagem
	res := ""
	var clock = req.Clock

	for k, v := range clock {
		res += "&" + k + "|" + strconv.Itoa(v)
	}

	return PP2PLink.PP2PLink_Req_Message{
		To:      to,
		Message: req.Message + "@" + res}

}

// conversor de Ind PP2Link => Ind BEB
func PP2PLink2BEB(ind PP2PLink.PP2PLink_Ind_Message) BestEffortBroadcast_Ind_Message {

	// decodificacao do relogio vetorial na mensagem
	dataMessageFrom := strings.Split(ind.Message, "§")
	message := dataMessageFrom[0]
	dataFromClock := strings.Split(dataMessageFrom[1], "@")
	from := dataFromClock[0]
	clocks := strings.Split(dataFromClock[1], "&")
	clock := make(map[string]int, len(clocks))

	for i := 1; i < len(clocks); i++ {
		dataIpValue := strings.Split(clocks[i], "|")
		if value, err := strconv.Atoi(dataIpValue[1]); err == nil {
			clock[dataIpValue[0]] = value
		}

	}
	return BestEffortBroadcast_Ind_Message{
		From:    from,
		Message: message,
		Clock:   clock,
	}

}
