package call

import (
	"fmt"
	"math/rand"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
)

// 调用合约方法
// @param client: 客户端
// @param contractName: 合约名称
// @param method: 调用或者查询，这里一般是设置成invoke_contract
// @param kvs: 参数列表
func InvokeContract(
	client *sdk.ChainClient,
	contractName string,
	kvs []*common.KeyValuePair,
	withSyncResult bool,
) (*common.TxResponse, error) {
	resp, err := client.InvokeContract(contractName, "invoke_contract", "", kvs, -1, withSyncResult)

	if err != nil {
		fmt.Println("##################################################")
		fmt.Println(err)
		return resp, fmt.Errorf(err.Error())
	}
	return resp, nil
}

// 调用合约方法
// @param c: 客户端
// @param contractName: 合约名称
// @param looptime 循环调用合约次数
func CallContract(c *sdk.ChainClient, contractName string, looptime int, withSyncResult bool) error {
	args := []*common.KeyValuePair{
		{
			// 调用哪个方法
			Key:   "method",
			Value: []byte(""),
		},
	}
	// 执行方法
	for i := 0; i < looptime; i++ {
		_, err := InvokeContract(c, contractName, args, withSyncResult)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		time.Sleep(time.Millisecond)
	}
	return nil
}

func CallStoreContrace(c *sdk.ChainClient, looptime int, withSyncResult bool) error {
	for i := 0; i < looptime; i++ {
		randKey := rand.Int()
		args := []*common.KeyValuePair{
			{
				// 调用哪个方法
				Key:   "method",
				Value: []byte(""),
			},
			{
				Key:   "key",
				Value: []byte(fmt.Sprintf("%d", randKey)),
			},
			{
				Key:   "value",
				Value: []byte(""),
			},
		}
		_, err := InvokeContract(c, "store", args, withSyncResult)
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		time.Sleep(500 * time.Microsecond)

	}
	return nil
}
