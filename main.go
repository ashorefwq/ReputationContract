package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 定义智能合约结构
type ReputationContract struct {
	contractapi.Contract
}

// 定义参与者的信誉度量结构
type Reputation struct {
	Participant string  // 参与者
	Score       float64 // 信誉度量分数
}

// 初始化智能合约
func (rc *ReputationContract) Init(ctx contractapi.TransactionContextInterface) error {
	fmt.Println("Reputation contract initialized")
	return nil
}

// 评估参与者的信誉度量，计算分数
func (rc *ReputationContract) EvaluateReputation(ctx contractapi.TransactionContextInterface, participant string) error {
	reputation := &Reputation{
		Participant: participant,
		Score:       0.0,
	}

	// 先随机生成一个分数
	reputation.Score = generateRandomScore()

	// 将信誉度量记录到分类账
	reputationJSON, _ := json.Marshal(reputation)
	err := ctx.GetStub().PutState(participant, reputationJSON)
	if err != nil {
		return fmt.Errorf("failed to put reputation data for participant %s: %w", participant, err)
	}

	return nil
}

// 获取参与者的信誉度量
func (rc *ReputationContract) GetReputation(ctx contractapi.TransactionContextInterface, participant string) (*Reputation, error) {
	reputationJSON, err := ctx.GetStub().GetState(participant)
	if err != nil {
		return nil, fmt.Errorf("failed to read reputation data for participant %s: %w", participant, err)
	}
	if reputationJSON == nil {
		return nil, fmt.Errorf("reputation data does not exist for participant %s", participant)
	}

	var reputation Reputation
	err = json.Unmarshal(reputationJSON, reputation)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal reputation data for participant %s: %w", participant, err)
	}

	return &reputation, nil
}

// 生成随机分数
func generateRandomScore() float64 {
	//此处可根据规则计算信誉度量分数
	rand.Seed(time.Now().UnixNano())
	return rand.Float64() * 100
}

func main() {
	// 创建智能合约实例并调用启动函数
	reputationContract := new(ReputationContract)

	cc, err := contractapi.NewChaincode(reputationContract)
	if err != nil {
		fmt.Printf("Error creating reputation smart contract: %s", err.Error())
		return
	}

	if err := cc.Start(); err != nil {
		fmt.Printf("Error starting reputation smart contract: %s", err.Error())
	}
}
