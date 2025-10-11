package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/8a3Cl1-qJs6wZGnTr0WDE")
	if err != nil {
		log.Fatal(err)
	}

	// 查询区块相关信息
	//blockNumber := big.NewInt(9387431)
	//block, err := client.BlockByNumber(context.Background(), blockNumber)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("block number: ", block.Number().Uint64())
	//fmt.Println("block timestamp: ", time.Unix(int64(block.Time()), 0).Format("2006-01-02 15:04:05"))
	//fmt.Println("block difficulty: ", block.Difficulty().Uint64())
	//fmt.Println("block hash:", block.Hash().Hex())
	//count, err := client.TransactionCount(context.Background(), block.Hash())
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("block transaction count: ", count)

	// 转账交易
	privateKey, err := crypto.HexToECDSA("8523fe792478014f750f7910375b3a98d70ab6c0b04a14e4f051da05f7d85f2d")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(10000000000000000) // in wei (0.01 eth)
	gasLimit := uint64(21000)              // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0x46B848559414E3b80142FD0D34A8818b7D540bef")
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	// 交易哈希：0x52578b9591f6fe37c7b6ede31a3c09718f815f41ed398d8d801f1132c092bee3
	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}
