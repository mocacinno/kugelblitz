package lightningrpc

import (
	"fmt"
	"net"

	"github.com/powerman/rpc-codec/jsonrpc2"
	log "github.com/sirupsen/logrus"
)

type LightningRpc struct {
	socketPath string
	conn       net.Conn
}

type Empty struct{}

type NewAddressResponse struct {
	Address string `json:"address"`
}

type GetInfoResponse struct {
	Id          string `json:"id"`
	Port        uint   `json:"port"`
	Testnet     bool   `json:"testnet"`
	Version     string `json:"version"`
	BlockHeight uint   `json:"blockheight"`
}

type GetPeersResult struct {
	Peers []Peer
}

type Peer struct {
	State       string `json:"state"`
	PeerId      string `json:"peerid"`
	Connected   bool   `json:"connected"`
	OurAmount   int    `json:"our_amount"`
	TheirAmount int    `json:"their_amount"`
	OurFee      int    `json:"our_fee"`
	TheirFee    int    `json:"their_fee"`
}

type Channel struct {
	From            string `json: "from"`
	To              string `json: "to"`
	BaseFee         uint   `json:"base_fee"`
	ProportionalFee uint   `json:"proportional_fee"`
}

type GetChannelsResponse struct {
	Channels []Channel `json:"channels"`
}

func (lr *LightningRpc) call(method string, req interface{}, res interface{}) error {
	log.Debugf("Calling lightning.%s with args %v", method, req)

	clientTCP, err := jsonrpc2.Dial("unix", lr.socketPath)
	if err != nil {
		return err
	}
	defer clientTCP.Close()
	err = clientTCP.Call(method, req, res)
	if err != nil {
		log.Debugf("error calling %s: %v", method, err)
		return errors.Wrap(err, fmt.Sprintf("error calling %s", method))
	} else {
		log.Debugf("method %s returned %v", method, err)
		return nil
	}
}

func (lr *LightningRpc) NewAddress() (NewAddressResponse, error) {
	res := NewAddressResponse{}
	err := lr.call("newaddr", &Empty{}, &res)
	return res, err
}

func (lr *LightningRpc) GetInfo() (GetInfoResponse, error) {
	res := GetInfoResponse{}
	err := lr.call("getinfo", &Empty{}, &res)
	return res, err
}

func (lr *LightningRpc) GetChannels() (GetChannelsResponse, error) {
	res := GetChannelsResponse{}
	err := lr.call("getchannels", &Empty{}, &res)
	return res, err
}

type GetPeersResponse struct {
	Peers []Peer `json:"peers"`
}

func (lr *LightningRpc) GetPeers() (GetPeersResponse, error) {
	res := GetPeersResponse{}
	err := lr.call("getpeers", &Empty{}, &res)
	return res, err
}

func (lr *LightningRpc) Connect(host string, port uint, fundingTx string) error {
	var params []interface{}
	params = append(params, host)
	params = append(params, port)
	params = append(params, fundingTx)
	return lr.call("connect", params, &Empty{})
}

type PeerReference struct {
	PeerId string `json:"peerid"`
}

func (lr *LightningRpc) Close(peerId string) error {
	var params []interface{}
	params = append(params, peerId)
	return lr.call("close", params, &Empty{})
}

type RouteHop struct {
	NodeId  string `json:"id"`
	Amount  uint64 `json:"msatoshi"`
	Delay   uint32 `json:"delay"`
	Channel string `json:"channel"`
}

type Route struct {
	Hops []RouteHop `json:"route"`
}

type GetRouteRequest struct {
	Destination string  `json:"destination"`
	Amount      uint64  `json:"amount"`
	RiskFactor  float32 `json:"risk"`
}

func (lr *LightningRpc) GetRoute(destination string, amount uint64, riskfactor float32) (Route, error) {
	var params []interface{}
	params = append(params, destination)
	params = append(params, amount)
	params = append(params, riskfactor)
	res := Route{}
	err := lr.call("getroute", params, &res)
	return res, err
}

type SendPaymentRequest struct {
	Route       []RouteHop `json:"route"`
	PaymentHash string     `json:"paymenthash"`
}

type SendPaymentResponse struct {
	PaymentKey string `json:"preimage"`
}

func (lr *LightningRpc) SendPayment(route []RouteHop, paymentHash string) (SendPaymentResponse, error) {
	var params []interface{}
	params = append(params, route)
	params = append(params, paymentHash)
	res := SendPaymentResponse{}
	err := lr.call("sendpay", params, &res)
	return res, err
}

type NodeAddress struct {
	Type    string `json:"type"`
	Address string `json:"address"`
	Port    uint16 `json:"port"`
}

type Node struct {
	Id        string        `json:"nodeid"`
	Addresses []NodeAddress `json:"addresses"`
}

type GetNodesResponse struct {
	Nodes []Node `json:"nodes"`
}

func (lr *LightningRpc) GetNodes() (GetNodesResponse, error) {
	res := GetNodesResponse{}
	err := lr.call("getnodes", &Empty{}, &res)
	return res, err
}

type ConnectRequest struct {
	Host         string `json:"host"`
	Port         uint   `json:"port"`
	FundingTxHex string `json:"tx"`
}

type Invoice struct {
	PaymentHash string `json:"rhash"`
	PaymentKey  string `json:"paymentKey"`
}

func (lr *LightningRpc) Invoice(amount uint64, label string) (Invoice, error) {
	var params []interface{}
	params = append(params, amount)
	params = append(params, label)
	res := Invoice{}
	err := lr.call("invoice", params, &res)
	return res, err
}

func NewLightningRpc(socketPath string) *LightningRpc {
	return &LightningRpc{
		socketPath: socketPath,
	}
}

func (lr *LightningRpc) Stop() error {
	var params []interface{}
	res := Empty{}
	return lr.call("stop", params, &res)
}
