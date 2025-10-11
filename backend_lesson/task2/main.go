package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	INFURA_URL  = "https://eth-sepolia.g.alchemy.com/v2/8a3Cl1-qJs6wZGnTr0WDE"
	PRIVATE_KEY = "8523fe792478014f750f7910375b3a98d70ab6c0b04a14e4f051da05f7d85f2d"
)

func main() {

	client, err := ethclient.Dial(INFURA_URL)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(PRIVATE_KEY)
	if err != nil {
		log.Fatal(err)
	}

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal(err)
	}

	// 部署合约
	input := "1.0"
	_, tx, instance, err := DeployCounter(auth, client, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("部署交易 hash:", tx.Hash().Hex())
	fmt.Println("等待交易确认...")

	// 等待交易确认
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status == 0 {
		log.Fatal("交易执行失败")
	}

	contractAddress := receipt.ContractAddress
	fmt.Println("合约部署成功，地址：", contractAddress.Hex())

	// 用新部署的合约地址创建实例
	instance, err = NewCounter(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// 调用合约方法
	tx, err = instance.Count(auth)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("调用 Count 交易 hash:", tx.Hash().Hex())

	// 等待 Count() 交易确认
	countReceipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}
	if countReceipt.Status == 0 {
		log.Fatal("Count 交易执行失败")
	}
	fmt.Println("Count 交易已确认！")

	// 读取计数器
	callOpt := &bind.CallOpts{Context: context.Background()}
	counter, err := instance.GetCounter(callOpt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("当前 counter:", counter)
}
