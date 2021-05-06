package CasualOrderBroadcast

import (
	"fmt"

	BEB "../BestEffortBroadcast"
)

// estrutura request
type CasualOrderBroadcast_Req_Message struct {
	Addresses []string
	Message   string
	Clock     map[string]int
}

// estrutura ind
type CasualOrderBroadcast_Ind_Message struct {
	From    string
	Message string
	Clock   map[string]int
}

// estrutura modulo
type CasualOrderBroadcast_Module struct {
	Ind     chan CasualOrderBroadcast_Ind_Message
	Req     chan CasualOrderBroadcast_Req_Message
	Beb     BEB.BestEffortBroadcast_Module
	Clock   map[string]int
	Lsn     int
	Pending []CasualOrderBroadcast_Ind_Message
	Ip      string
}

func (cob *CasualOrderBroadcast_Module) Init(ip string, address []string) {

	cob.Beb = BEB.BestEffortBroadcast_Module{
		Req: make(chan BEB.BestEffortBroadcast_Req_Message),
		Ind: make(chan BEB.BestEffortBroadcast_Ind_Message)}

	// iniciando relogio vetorial do processo (clock[ip]=0)
	cob.Clock = make(map[string]int, len(address))

	// garantindo que o relogio foi corretamente inicializado
	for i := 0; i < len(address); i++ {
		cob.Clock[address[i]] = 0
		fmt.Printf("%v - Clock[%v] = %v \n", ip, address[i], cob.Clock[address[i]])
	}

	cob.Lsn = 0
	cob.Ip = ip

	cob.Beb.Init(ip)
	cob.Start()
}

func (cob CasualOrderBroadcast_Module) Start() {

	go func() {
		for {
			select {

			// caso request => processar request => enviar para camada inferior (BEB)
			case req := <-cob.Req:
				x := cob.processReq(req)
				cob.Beb.Req <- x

			// caso ind => processar ind
			case ind := <-cob.Beb.Ind:
				indComplete := BEBIndToCOBInd(ind)
				cob.processInd(indComplete)
			}
		}
	}()

}

func (cob *CasualOrderBroadcast_Module) processInd(ind CasualOrderBroadcast_Ind_Message) {
	// adiciona ind a lista de pendentes => processa pendentes
	cob.Pending = append(cob.Pending, ind)
	cob.processPendings()
}

func (cob *CasualOrderBroadcast_Module) processPendings() {
	tmp := make([]CasualOrderBroadcast_Ind_Message, len(cob.Pending))
	copy(tmp, cob.Pending)

	for i, ind := range tmp {
		higher := false
		for k := range ind.Clock {
			if ind.Clock[k] > cob.Clock[k] {
				higher = true
			}
		}
		// caso o relogio vetorial da ind nao seja maior que o do modulo => processar ind
		if !higher {
			cob.Pending = append(cob.Pending[:i], cob.Pending[i+1:]...)
			cob.Clock[ind.From] = ind.Clock[ind.From] + 1
			cob.Ind <- ind
			cob.processPendings()
			break
		}
	}
}

func (cob *CasualOrderBroadcast_Module) processReq(req CasualOrderBroadcast_Req_Message) BEB.BestEffortBroadcast_Req_Message {

	// mapear relogio da aplicacao para relogio do request
	clock := make(map[string]int)
	for k, v := range cob.Clock {
		clock[k] = v
	}
	tmp := cob.Lsn
	clock[cob.Ip] = tmp

	req.Clock = clock

	// incrementar lsn
	cob.Lsn = cob.Lsn + 1

	return COBReqToBEBReq(req)
}

// conversao de Request COB => Request BEB
func COBReqToBEBReq(req CasualOrderBroadcast_Req_Message) BEB.BestEffortBroadcast_Req_Message {

	return BEB.BestEffortBroadcast_Req_Message{
		Addresses: req.Addresses,
		Message:   req.Message,
		Clock:     req.Clock}

}

// conversao de Ind Beb => Ind COB
func BEBIndToCOBInd(ind BEB.BestEffortBroadcast_Ind_Message) CasualOrderBroadcast_Ind_Message {

	return CasualOrderBroadcast_Ind_Message{
		From:    ind.From,
		Message: ind.Message,
		Clock:   ind.Clock}
}
