package CasualOrderBroadcast

import (
	"fmt"

	BEB "../BestEffortBroadcast"
	PP2PLink "../PP2PLink"
)

type CasualOrderBroadcast_Req_Message struct {
	Addresses []string
	Message   string
	Clock     map[string]int
}

type CasualOrderBroadcast_Ind_Message struct {
	From    string
	Message string
	Clock   map[string]int
}
type CasualOrderBroadcast_Module struct {
	Ind     chan CasualOrderBroadcast_Ind_Message
	Req     chan CasualOrderBroadcast_Req_Message
	Beb     BEB.BestEffortBroadcast_Module
	Clock   map[string]int
	Lsn     int
	Pending []CasualOrderBroadcast_Ind_Message
	ip      string
}

func (cob CasualOrderBroadcast_Module) Init(address []string) {

	fmt.Println("Init COB!")
	cob.Beb = BEB.BestEffortBroadcast_Module{
		Req: make(chan BEB.BestEffortBroadcast_Req_Message),
		Ind: make(chan BEB.BestEffortBroadcast_Ind_Message)}

	// initializing clocks (clock[ip]=0)
	cob.Clock = make(map[string]int, len(address))

	for i := 0; i < len(address); i++ {
		cob.Clock[address[i]] = 0
		fmt.Printf("Clock[%v] = %v \n", address[i], cob.Clock[address[i]])
	}

	// starts lsn at 0
	cob.Lsn = 0

	cob.ip = address[0]

	cob.Beb.Init(address[0])
	cob.Start()
}

func (cob CasualOrderBroadcast_Module) Start() {

	go func() {
		for {
			select {
			case req := <-cob.Req:
				cob.Beb.Req <- cob.processReq(req)

			case ind := <-cob.Beb.Ind:
				indComplete := BEBIndToCOBInd(ind)
				cob.processInd(indComplete)
			}
		}
	}()

}

func (cob CasualOrderBroadcast_Module) processInd(req CasualOrderBroadcast_Ind_Message) {

	// deliver
	cob.Ind <- req
}

func (cob *CasualOrderBroadcast_Module) processReq(req CasualOrderBroadcast_Req_Message) BEB.BestEffortBroadcast_Req_Message {
	req.Clock = cob.Clock
	req.Clock[cob.ip] = cob.Lsn
	cob.Lsn++
	return COBReqToBEBReq(req)
}

func COBReqToBEBReq(req CasualOrderBroadcast_Req_Message) BEB.BestEffortBroadcast_Req_Message {

	return BEB.BestEffortBroadcast_Req_Message{
		Addresses: req.Addresses,
		Data: PP2PLink.PP2LinkMessage{
			Message: req.Message,
			Clock:   req.Clock}}

}

func BEBIndToCOBInd(ind BEB.BestEffortBroadcast_Ind_Message) CasualOrderBroadcast_Ind_Message {

	return CasualOrderBroadcast_Ind_Message{
		From:    ind.From,
		Message: ind.Data.Message,
		Clock:   ind.Data.Clock}
}
