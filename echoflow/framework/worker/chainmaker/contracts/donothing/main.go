// coding:utf-8
// 存储合约
package main

import (
	"log"

	pb "chainmaker.org/chainmaker/contract-sdk-go/v2/pb/protogo"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sandbox"
	"chainmaker.org/chainmaker/contract-sdk-go/v2/sdk"
)

type DONOTHING struct{}

func (s *DONOTHING) InitContract() pb.Response {
	return sdk.Success([]byte("success"))
}

// UpgradeContract use to upgrade contract
func (h *DONOTHING) UpgradeContract() pb.Response {
	return sdk.Success([]byte("Upgrade success"))
}

// InvokeContract use to select specific method
func (h *DONOTHING) InvokeContract(method string) pb.Response {
	// according method segment to select contract functions
	return sdk.Success([]byte("success"))
}

// main
func main() {
	err := sandbox.Start(new(DONOTHING))
	if err != nil {
		log.Fatal(err)
	}
}
