package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	strconv "strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
/*
查询商品信息

*/
// ============================================================================================================================
func queryGoods(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("-------------------------start queryGoods---------------------------")
	var err error

	if len(args) != 1 {
		return shim.Error("信息输入错误！")
	}
	id := args[0]
	value, err := stub.GetState(id)
	if err != nil {
		fmt.Println("Failed to get state for " + id)
		return shim.Error(err.Error())
	}
	var goods Goods
	json.Unmarshal(value, &goods)
	fmt.Println(goods)
	fmt.Println("---------------------------end queryGoods-----------------------------------")
	return shim.Success(value)
}

// ============================================================================================================================
/*
查询商品历史信息

*/
// ============================================================================================================================
func getHistoryForGoods(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var buffer bytes.Buffer

	if len(args) != 1 {
		return shim.Error("信息输入错误！")
	}

	goods_id := args[0]
	fmt.Println("---------------start getHistoryForGoods ID is " + goods_id)
	//get history
	results, err := stub.GetHistoryForKey(goods_id)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer results.Close()
	isWrite := false
	buffer.WriteString("[")
	for results.HasNext() {
		response, err := results.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if isWrite == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}
		buffer.WriteString(",\"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		isWrite = true
	}
	buffer.WriteString("]")
	fmt.Printf("getHistoryForGoods returning : \n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}
