// coding:utf-8
// 存储合约
package main

import (
	"log"

	pb "chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type STORE struct{}

func (s *STORE) InitContract() pb.Response {
	return sdk.Success([]byte("success"))
}

// UpgradeContract use to upgrade contract
func (h *STORE) UpgradeContract() pb.Response {
	return sdk.Success([]byte("Upgrade success"))
}

// InvokeContract use to select specific method
func (h *STORE) InvokeContract(method string) pb.Response {
	// according method segment to select contract functions
	args := sdk.Instance.GetArgs()
	key := string(args["key"])
	val := args["value"]
	sdk.Instance.PutStateFromKeyByte(key, val)
	return sdk.Success([]byte("success" + ", key:" + key + ", val:" + string(val)))
}

// main
func main() {
	err := sandbox.Start(new(STORE))
	if err != nil {
		log.Fatal(err)
	}
}
