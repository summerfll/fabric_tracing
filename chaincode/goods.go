package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type GoodsChaincode struct {
}

//货物的基本信息
type Goods struct {
	ObjectType   string        `json:"docType"`      //field for couchdb
	Id           string        `json:"id"`           //商品id
	Name         string        `json:"name"`         //货物名称
	BornTime     string        `json:"borntime"`     //生产时间
	BornSpace    string        `json:"bornspace"`    //生产地
	Manufacturer string        `json:"manufacturer"` //生产厂商
	Owner        OwnerRelation `json:"owner"`        //拥有者：生产商，商1，商2
}

// ----- Owners ----- //
type Owner struct {
	ObjectType string `json:"docType"`   //field for couchdb
	Id         string `json:"id"`        //仓库id
	Storename  string `json:"storename"` //仓库名称
	//InputStore  string `json:"inputstore"`//入库时间
	//OutputStore string `json:"outputstore"`//出库时间
	Merchant string `json:"merchant"` //商家
}

type OwnerRelation struct {
	Id        string `json:"id"`
	Storename string `json:"storename"` //仓库名称
	Merchant  string `json:"merchant"`  //商家
}

func (t *GoodsChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil) //nil一般为NULL的意思
}

func (t *GoodsChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)
	if function == "initGoods" { //创建商品
		return initGoods(stub, args)
	} else if function == "initOwner" { //创建仓储
		return initOwner(stub, args)
	} else if function == "transferGoods" { //转移商品
		return transferGoods(stub, args)
	} else if function == "queryGoods" { //查询商品信息
		return queryGoods(stub, args)
	} else if function == "deleteGoods" { //删除商品
		return deleteGoods(stub, args)
	} else if function == "getHistoryForGoods" { //查询商品历史记录
		return getHistoryForGoods(stub, args)
	}
	return shim.Error("没有相应的方法！")
}

// ============================================================================================================================
// 查询
// ============================================================================================================================
func (t *GoodsChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Unknown supported call - Query()")
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(GoodsChaincode))
	if err != nil {
		fmt.Printf("Error starting goods chaincode:%s", err)
	}
}
