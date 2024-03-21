package lib

type NetworkConfig struct {
	VethName     string `json:"veth_name"`
	VethAddr     string `json:"veth_addr"`
	BridgeName   string `json:"bridge_name"`
	NS           string `json:"ns"`
	DefaultRoute string `json:"default_route"`
}
