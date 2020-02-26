package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/simplechain-org/go-simplechain/accounts/abi"
	"github.com/simplechain-org/go-simplechain/common"
	"github.com/simplechain-org/go-simplechain/common/hexutil"
	"github.com/simplechain-org/go-simplechain/rpc"
	"io/ioutil"
	"log"
	"math/big"
)

var rawurlVar *string =flag.String("rawurl", "http://127.0.0.1:8556", "rpc url")

var abiPath *string =flag.String("abi", "../../contracts/crossdemo/crossdemo.abi", "abi文件路径")

var contract *string =flag.String("contract", "0x8eefA4bFeA64F2A89f3064D48646415168662a1e", "合约地址")

var signConfirm *uint64=flag.Uint64("signConfirm", 2, "最小锚定节点数")

var chainId *uint64=flag.Uint64("chainId", 1, "目的链id")

var fromVar *string=flag.String("from", "0xb9d7df1a34a28c7b82acc841c12959ba00b51131", "发起人地址")

var gaslimitVar *uint64=flag.Uint64("gaslimit", 2000000, "gas最大值")

var anchor1 *string = flag.String("anchor1","0x6051De4667626B97af2b81A392ad228e0fF58002","锚定节点名单")
var anchor2 *string = flag.String("anchor2","0x8e422d5Aff496974f7FaE17F6848a40C59F8b2E9","锚定节点名单")
var anchor3 *string = flag.String("anchor3","0x935d0d6851c8db45C75D2DD66A630db22A1a918A","锚定节点名单")


type SendTxArgs struct {
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`
}

func Register(client *rpc.Client) {

	data,err:=ioutil.ReadFile(*abiPath)

	if err!=nil{
		fmt.Println(err)
		return
	}

	from := common.HexToAddress(*fromVar)

	to := common.HexToAddress(*contract)

	gas := hexutil.Uint64(*gaslimitVar)

	value := hexutil.Big(*big.NewInt(0))

	abi, err := abi.JSON(bytes.NewReader(data))

	if err != nil {
		log.Fatalln(err)
	}
	////想要随机，请设置随机种子
	//rand.Seed(time.Now().UnixNano())
	//
	////这个值必须改变，因为已经限定合约中，计算出来的txid不能相同
	//nonce:=big.NewInt(0).SetUint64(rand.Uint64())
	//nonce:=big.NewInt(0).SetUint64(9536605289005490782)

	remoteChainId:=big.NewInt(0).SetUint64(*chainId)

	signConfirmCount:=uint8(*signConfirm)

	anchorA := common.HexToAddress(*anchor1)

	anchorB := common.HexToAddress(*anchor2)

	anchorC := common.HexToAddress(*anchor3)

	var anchors []common.Address
	anchors = append(anchors,anchorA)
	anchors = append(anchors,anchorB)
	anchors = append(anchors,anchorC)

	out, err := abi.Pack("chainRegister",remoteChainId ,signConfirmCount,anchors)

	input := hexutil.Bytes(out)

	fmt.Println("input=",input)

	args := &SendTxArgs{
		From:  from,
		To:    &to,
		Gas:   &gas,
		Value: &value,
		Input: &input,
	}

	var result common.Hash

	err = client.CallContext(context.Background(), &result, "eth_sendTransaction", args)

	if err != nil {
		fmt.Println("CallContext", "err", err)
		return
	}

	fmt.Println("result=", result.Hex())
}


//跨链交易发起人
func main() {
	flag.Parse()
	client, err := rpc.Dial(*rawurlVar)
	if err != nil {
		fmt.Println("dial", "err", err)
		return
	}

	Register(client)
}

