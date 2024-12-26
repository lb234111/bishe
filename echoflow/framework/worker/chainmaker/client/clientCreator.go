// coding:utf-8
// 创建chainmaker客户端
package client

import (
	"fmt"
	"framework/utils"
	"strings"
	"sync"

	sdk "chainmaker.org/chainmaker/sdk-go/v2"
)

// 生成chainmaker证书配置文件
// @param clientNum 客户端数量
func generateChainmakerConfigYaml(clientNum int) {
	// 读取文件
	ls_content := utils.ReadFileWithBufio("./worker/chainmaker/config_files/chainmaker/chainmaker_user_template.yml")
	for chain_id := 0; chain_id < 1; chain_id++ {
		curr_chain_id := fmt.Sprintf("chain%d", chain_id+1)
		for org_id := 0; org_id < 11; org_id++ {
			curr_org_id := fmt.Sprintf("org%d", org_id+1)
			var wg sync.WaitGroup
			wg.Add(clientNum)
			// 生成每一个客户端文件
			for i := 0; i < clientNum; i++ {
				curr_client_id := fmt.Sprintf("client%d", i+1)
				curr_port := fmt.Sprintf("%d", 12301+org_id)
				go func(chain_id, org_id, client_id, curr_port string) {
					// 写入的文件路径
					save_path := fmt.Sprintf("./worker/chainmaker/config_files/chainmaker/tmpUserConfig/%s_%s_%s.yml", chain_id, org_id, client_id)

					// 文件存在就不生成
					if utils.FileExist(save_path) {
						wg.Done()
						return
					}

					// 生成文件内容
					ls_result := []string{}
					for _, line := range ls_content {
						// 替换chain_id
						curr_line := strings.ReplaceAll(line, "{chain_id}", chain_id)
						// 替换org_id
						curr_line = strings.ReplaceAll(curr_line, "{org_id}", org_id)
						// 替换client_id
						curr_line = strings.ReplaceAll(curr_line, "{client_id}", client_id)
						// 替换port
						curr_line = strings.ReplaceAll(curr_line, "{port}", curr_port)

						ls_result = append(ls_result, curr_line)
					}

					// 写入文件
					utils.RewriteFile(save_path, ls_result)

					fmt.Println(save_path)

					wg.Done()
				}(curr_chain_id, curr_org_id, curr_client_id, curr_port)
			}
			wg.Wait()
		}
	}

}

// 返回实例化完成的客户端
// @param clientNumber 生成客户端个数
// @param chain_name 链名, chain1 or chain2
// @param org_name 组织名, org1 or org2
func CreateChainmakerClient(clientNumber, orgNum int, chain_name string) (aaa []*sdk.ChainClient, ee error) {
	fmt.Println("生成证书")
	// 生成证书
	generateChainmakerConfigYaml(clientNumber)

	// 生成客户端
	ls_client := make([]*sdk.ChainClient, 0)

	// 锁
	var mutex sync.Mutex
	// 等待组
	var wg sync.WaitGroup

	fmt.Println("加入配置文件")
	// 加入配置文件
	chainmaker_config_path_pattern := "./worker/chainmaker/config_files/chainmaker/tmpUserConfig/%s_%s_%s.yml"

	// 开始生成客户端
	for org_id := 0; org_id < orgNum; org_id += 1 {
		curr_org_name := fmt.Sprintf("org%d", org_id+1)
		for user_id := 0; user_id < clientNumber; user_id += 1 {
			curr_user_name := fmt.Sprintf("client%d", user_id+1)
			val := fmt.Sprintf(chainmaker_config_path_pattern, chain_name, curr_org_name, curr_user_name)
			currClient, err := NewClient(val)
			if err != nil {
				fmt.Println("生成客户端错误: ", err)
				return nil, err
			}
			mutex.Lock()
			ls_client = append(ls_client, currClient)
			mutex.Unlock()
			fmt.Println("加入完成: ", val)
		}
	}

	utils.PrintInfo("等待客户端启动完成")
	wg.Wait()
	return ls_client, nil
}
