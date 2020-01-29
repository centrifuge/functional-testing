package utils

import (
	"os"
	"strings"
	"testing"

	"github.com/gavv/httpexpect"
)

const (
	NODE1 = "node1"
	NODE2 = "node2"
)

var Nodes map[string]node
var Network string
var RegistryAddress string
var AssetAddress string

type node struct {
	ID   string
	HOST string
}

func SetupEnvironment() {
	nodesEnv := os.Getenv("NODES")
	idsEnv := os.Getenv("IDS")
	nodesSlice := SplitString(nodesEnv)
	idsSlice := SplitString(idsEnv)

	if len(nodesSlice) == 0 {
		nodesSlice = append(nodesSlice, "http://localhost:8082", "http://localhost:8083")
	}

	if len(idsSlice) == 0 {
		idsSlice = append(idsSlice, "0xF4F58f2e86C8546d7BE1ED43b347f09a91f85be8", "0x331db0dCDd37ceAD608416df2697c0F28c246f8a")
	}

	Nodes = map[string]node{
		NODE1: {
			idsSlice[0],
			nodesSlice[0],
		},
		NODE2: {
			idsSlice[1],
			nodesSlice[1],
		},
	}

	Network = os.Getenv("NETWORK")
	if Network == "" {
		Network = "fulvous"
	}

	RegistryAddress = os.Getenv("NFT_REGISTRY")
	AssetAddress = os.Getenv("NFT_ASSET")
	if RegistryAddress == "" || AssetAddress == "" {
		RegistryAddress = "0xb56ac47948157a7259ac2b72b950193a7fa40f81"
		AssetAddress = "0x445982361d26f22f1e3c0385eb6c31b3f423acb1"
	}
}

func GetInsecureClient(t *testing.T, nodeId string) *httpexpect.Expect {
	SetupEnvironment()
	return CreateInsecureClient(t, Nodes[nodeId].HOST)
}

func SplitString(data string) []string {
	result := strings.Split(data, ",")
	if result[0] == "" {
		return []string{}
	}

	return result
}
