package scan

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strings"
	"template/pkg/data"
	"time"

	"github.com/go-errors/errors"
	"github.com/golang/glog"
)

type Engine struct {
	ctx context.Context
	//infoDTO    dtos.InfoDTO
	//httpClient *http.Client
	bscCient *ethclient.Client
}

func NewEngine(ctx context.Context) (*Engine, error) {
	return &Engine{
		ctx: ctx,
		/*infoDTO: dtos.InfoDTO{
			//TotalIssuance:    0,
			//TotalCirculation: 0,
			TotalAccounts: 0,
			TotalTokens:   0,
			TotalBlocks:   0,
			TotalTxs:      0,
		},
		httpClient: data.MustGetHttp("node"),*/
		bscCient: data.MustGetBscClient("mainnet"),
	}, nil
}

func (engine *Engine) Run(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			goErr := errors.Wrap(err, 1)
			glog.Errorf("%s\n stack: %s", goErr.Error(), goErr.Stack())
		}
	}()

	syncTicker := time.NewTicker(time.Second * 3)
	defer func() {
		syncTicker.Stop()
	}()

	processCountTicker := time.NewTicker(time.Minute * 1)
	defer func() {
		processCountTicker.Stop()
	}()

	for {
		select {
		case <-syncTicker.C:
			//engine.sync()
			fmt.Println(">>>>>>>>>>>>>>>>>>>>定时任务111111<<<<<<<<<<<<<<<<<")
			engine.readingEventLogs()
			fmt.Println(">>>>>>>>>>>>>>>>>>>>定时任务111111 任务执行结束<<<<<<<<<<<<<<<<<")
		case <-processCountTicker.C:
			//engine.processCount()
			//fmt.Println(">>>>>>>>>>>>>>>>>>>>定时任务222222<<<<<<<<<<<<<<<<<")
		case <-ctx.Done():
			break
		}
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

func (engine *Engine) readingEventLogs() {
	// 0x Protocol (ZRX) token address
	contractAddress := common.HexToAddress("0x094B109d635c34f24b1dC784b2b09F30ccf6408C")
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(28054627),
		ToBlock:   big.NewInt(28054627 + 5000),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := engine.bscCient.FilterLogs(context.Background(), query)
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

func (engine *Engine) Run2(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			goErr := errors.Wrap(err, 1)
			glog.Errorf("%s\n stack: %s", goErr.Error(), goErr.Stack())
		}
	}()

	processInfoTicker := time.NewTicker(time.Second * 30)
	defer func() {
		processInfoTicker.Stop()
	}()

	//engine.processInfo()
	/*
		for {
			select {
			case <-processInfoTicker.C:
				engine.processInfo()
			case <-ctx.Done():
				break
			}
		}*/
}

func (engine *Engine) sync() {
	/*getInfoReq := &helpers.CommonReq{JsonRpc: "1.0", Id: "curltest", Method: "getinfo"}
	getInfoRsp, err := helpers.SendGetInfo(getInfoReq)
	if err != nil {
		glog.Errorf("call getinfo rpc failed: %s", err.Error())
		return
	}

	infoDTOInternal := getInfoRsp.Result
	targetSyncBlockHeight := infoDTOInternal.SynBlockHeight
	currentSyncBlockHeight := helpers.GetLatestBlockHeight()

	// nothing to do here
	if currentSyncBlockHeight == targetSyncBlockHeight {
		return
	}

	for currentSyncBlockHeight++; currentSyncBlockHeight <= targetSyncBlockHeight; currentSyncBlockHeight++ {
		// 1. fetch data from node
		getBlockReq := &helpers.CommonReq{JsonRpc: "1.0", Id: "curltest", Method: "getblock", Params: []interface{}{currentSyncBlockHeight}}
		getBlockRsp, err := helpers.SendGetBlock(getBlockReq)
		if err != nil {
			glog.Errorf("call getblock rpc failed: %s", err.Error())
			return
		}

		// 2. write block into db
		blockDTO := getBlockRsp.Result
		block, _ := helpers.BlockDTO2Model(blockDTO)
		dao.SetBlock(block)
		dao.UpdatePreviousBlock(block)

		if targetSyncBlockHeight-currentSyncBlockHeight < util.DefaultDisableWebSocketBlockHeight {
			ws.BroadcastChannel.Broadcast(ws.Message{
				Type: util.MessageTypeBlock,
				Data: blockDTO,
			})
		}

		// 3. write tx into db
		for _, txid := range blockDTO.Tx {
			getTxReq := &helpers.CommonReq{JsonRpc: "1.0", Id: "curltest", Method: "gettxdetail", Params: []interface{}{txid}}
			getTxRsp, err := helpers.SendGetTx(getTxReq)
			if err != nil {
				glog.Errorf("call gettxdetail rpc failed: %s", err.Error())
				return
			}

			txDTO := getTxRsp.Result
			tx, _ := helpers.TxDTO2Model(txDTO)
			dao.WriteTx(tx)

			accounts := helpers.FindRelatedAccounts(txDTO)
			for addr, _ := range accounts {
				accountTx := models.AccountTx{
					//ID:   0,
					Addr: addr,
					TxID: txid,
				}
				dao.SetAccountTx(accountTx)

				getAccountReq := &helpers.CommonReq{JsonRpc: "1.0", Id: "curltest", Method: "getaccountinfo", Params: []interface{}{addr}}
				getAccountRsp, err := helpers.SendGetAccount(getAccountReq)
				if err != nil {
					return
				}

				accountDTO := getAccountRsp.Result
				account, _ := helpers.AccountDTO2Model(accountDTO)
				if oldAccount, ret := dao.AccountExist(accountDTO.Address); ret {
					account.ID = oldAccount.ID
					account.CreateTime = oldAccount.CreateTime
				} else {
					account.CreateTime = tx.ConfirmedTime
				}
				dao.SetAccount(account) // update or create
			}

			if targetSyncBlockHeight-currentSyncBlockHeight < util.DefaultDisableWebSocketBlockHeight {
				ws.BroadcastChannel.Broadcast(ws.Message{
					Type: util.MessageTypeTx,
					Data: txDTO,
				})
			}
		}

		glog.Infof("currentSyncBlockHeight: %d, targetSyncBlockHeight: %d", currentSyncBlockHeight, targetSyncBlockHeight)
	}

	engine.processInfo()

	ws.BroadcastChannel.Broadcast(ws.Message{
		Type: util.MessageTypeInfo,
		Data: engine.infoDTO,
	})*/
}
