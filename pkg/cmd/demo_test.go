package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"log"
	"math/big"
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	sum := add(2, 3)
	if sum != 5 {
		t.Errorf("Expected 5 but got %d", sum)
	}
}

func add(x, y int) int {
	return x + y
}

func TestConnect(t *testing.T) {
	client, err := ethclient.Dial("https://data-seed-prebsc-1-s3.binance.org:8545")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("we have a connection")
	_ = client // we'll use this in the upcoming sections

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(header.Number.String())
	}
}

// LogTransfer ..
type LogTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

// LogApproval ..
type LogApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
}

func TestReadingEventLogs(t *testing.T) {
	client, err := ethclient.Dial("https://bsc-dataseed4.ninicoin.io")
	if err != nil {
		log.Fatal(err)
	}

	// 0x Protocol (ZRX) token address
	contractAddress := common.HexToAddress("0x094B109d635c34f24b1dC784b2b09F30ccf6408C")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(28054627),
		ToBlock:   big.NewInt(28054627 + 5000),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(ERC20MetaData.ABI)))
	if err != nil {
		log.Fatal(err)
	}

	logTransferSig := []byte("Transfer(address,address,uint256)")
	LogApprovalSig := []byte("Approval(address,address,uint256)")
	logTransferSigHash := crypto.Keccak256Hash(logTransferSig)
	logApprovalSigHash := crypto.Keccak256Hash(LogApprovalSig)

	for _, vLog := range logs {
		fmt.Printf("Log Block Number: %d\n", vLog.BlockNumber)
		fmt.Printf("Log Index: %d\n", vLog.Index)
		fmt.Printf("Log TxHash: %s\n", vLog.TxHash)

		switch vLog.Topics[0].Hex() {
		case logTransferSigHash.Hex():
			fmt.Printf("Log Name: Transfer\n")

			var transferEvent LogTransfer

			err := contractAbi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			transferEvent.From = common.HexToAddress(vLog.Topics[1].Hex())
			transferEvent.To = common.HexToAddress(vLog.Topics[2].Hex())

			fmt.Printf("From: %s\n", transferEvent.From.Hex())
			fmt.Printf("To: %s\n", transferEvent.To.Hex())
			fmt.Printf("Value: %s\n", transferEvent.Value.String())

		case logApprovalSigHash.Hex():
			fmt.Printf("Log Name: Approval\n")

			var approvalEvent LogApproval

			err := contractAbi.UnpackIntoInterface(&approvalEvent, "Approval", vLog.Data)
			if err != nil {
				log.Fatal(err)
			}

			approvalEvent.Owner = common.HexToAddress(vLog.Topics[1].Hex())
			approvalEvent.Spender = common.HexToAddress(vLog.Topics[2].Hex())

			fmt.Printf("Token Owner: %s\n", approvalEvent.Owner.Hex())
			fmt.Printf("Spender: %s\n", approvalEvent.Spender.Hex())
			fmt.Printf("Value: %s\n", approvalEvent.Value.String())
		}

		fmt.Printf("\n\n")
	}
}

func TestSubscribeEvent(t *testing.T) {
	client, err := ethclient.Dial("wss://snowy-alien-liquid.bsc.discover.quiknode.pro/ad3f1d96711a722a587e4f52b4fdef9eefb2e1c3/")
	if err != nil {
		log.Fatal(err)
	}

	// 0x Protocol (ZRX) token address
	contractAddress := common.HexToAddress("0x17d70FF5670742d9bf7609AB3d1C790210ace997")
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println(vLog) // pointer to event log
		}
	}
}

const (
	DOMAIN_NAME    = "MyContract"
	DOMAIN_VERSION = "1"
	DOMAIN_CHAINID = 1
)

type Person struct {
	Name string `json:"name"`
	Age  uint   `json:"age"`
}

func TestEIP712Signature(t *testing.T) {

	signerData := apitypes.TypedData{
		Types: apitypes.Types{
			"Challenge": []apitypes.Type{
				{Name: "address", Type: "address"},
				{Name: "nonce", Type: "string"},
				{Name: "timestamp", Type: "string"},
			},
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "version", Type: "string"},
				{Name: "salt", Type: "string"},
			},
		},
		PrimaryType: "Challenge",
		Domain: apitypes.TypedDataDomain{
			Name:    "ETHChallenger",
			Version: "1",
			Salt:    salt,
			ChainId: math.NewHexOrDecimal256(1),
		},
		Message: apitypes.TypedDataMessage{
			"timestamp": timestamp,
			"address":   walletAddress,
			"nonce":     nonce,
		},
	}

	typedDataHash, _ := signerData.HashStruct(signerData.PrimaryType, signerData.Message)
	domainSeparator, _ := signerData.HashStruct("EIP712Domain", signerData.Domain.Map())

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	challengeHash := crypto.Keccak256Hash(rawData)
	fmt.Println(challengeHash)

	privateKey, err := crypto.HexToECDSA("")
	signature, err := crypto.Sign(challengeHash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hexutil.Encode(signature))

	uint256Ty, _ := abi.NewType("uint256", "", nil)
	bytesTy, _ := abi.NewType("bytes", "", nil)
	addressTy, _ := abi.NewType("address", "", nil)

	arguments := abi.Arguments{
		{
			Type: addressTy,
		},
		{
			Type: bytesTy,
		},
		{
			Type: uint256Ty,
		},
	}

	bytes, _ := arguments.Pack(
		common.HexToAddress("0x0000000000000000000000000000000000000000"),
		signature,
		big.NewInt(42),
	)

	crypto.Keccak256Hash(bytes)

	log.Println(crypto.Keccak256Hash(bytes))

}
