package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================================================================================
/*
初始化商品

*/
// ============================================================================================================================
func initGoods(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("-------------------------start initGoods---------------------------")
	var err error

	if len(args) != 6 {
		return shim.Error("信息输入错误！")
	}
	id := args[0]
	name := args[1]
	borntime := args[2]
	bornspace := args[3]
	manufacturer := args[4]
	owner_id := args[5]
	//检查owner是否存在
	owner, err := get_owner(stub, owner_id)
	if err != nil {
		fmt.Println("Failed to find owner - " + owner_id)
		return shim.Error(err.Error())
	}
	//检查商品id是否已经存在
	_, err = get_goods(stub, id)
	if err == nil {
		fmt.Println("This goods already exists - " + id)
		return shim.Error("This goods already exists - " + id)
	}
	str := `{
		"docType":"goods", 
		"id": "` + id + `", 
		"name": "` + name + `", 
		"borntime": "` + borntime + `",
		"bornspace": "` + bornspace + `",
		"manufacturer": "` + manufacturer + `",
		"owner": {
			"id": "` + owner_id + `", 
			"storename": "` + owner.Storename + `", 
			"merchant": "` + owner.Merchant + `"
		}
	}`
	err = stub.PutState(id, []byte(str)) //将商品信息存入couthDB数据库
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("----------------------------end initGoods------------------------")
	return shim.Success(nil)

}

// =============================================================================================================
/*
初始化owner
*/
// =============================================================================================================
func initOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("---------------starting initOwner---------------")

	if len(args) != 3 {
		return shim.Error("信息输入错误！")
	}

	var owner Owner
	owner.ObjectType = "goods_owner"
	owner.Id = args[0]
	owner.Storename = strings.ToLower(args[1]) //字符串转换，转换为小写
	owner.Merchant = args[2]
	fmt.Println(owner)
	//检查owner是否存在，存在即err为nil
	_, err = get_owner(stub, owner.Id)
	if err == nil {
		fmt.Println("This owner already exists - " + owner.Id)
		return shim.Error("This owner already exists - " + owner.Id)
	}
	//存储owner信息
	ownerAsBytes, _ := json.Marshal(owner) //转为json格式
	err = stub.PutState(owner.Id, ownerAsBytes)
	if err != nil {
		fmt.Println("Could not store user")
		return shim.Error(err.Error())
	}
	fmt.Println("-----------------end initOwner---------------")
	return shim.Success(nil)
}

// =============================================================================================================
/*
转移商品，改变商品的owner
1.检查owner是否存在，检查商品是否存在
2.从数据库中获得商品信息，更新owner的id及相关信息
3.将改变后的商品信息存入数据库
*/
// =============================================================================================================
func transferGoods(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("-------------------starting transferGoods------------------------")

	if len(args) != 2 {
		return shim.Error("信息输入错误！")
	}
	goods_id := args[0]
	new_owner_id := args[1]

	//检查owner是否存在
	owner, err := get_owner(stub, new_owner_id)
	if err != nil {
		return shim.Error("This owner does not exist - " + new_owner_id)
	}
	//检查商品是否存在
	_, err = get_goods(stub, goods_id)
	if err != nil {
		return shim.Error("This goods does not exist - " + goods_id)
	}
	// 在数据库中获取商品信息
	goodsAsBytes, err := stub.GetState(goods_id)
	if err != nil {
		return shim.Error("Failed to get goods")
	}
	res := Goods{}
	json.Unmarshal(goodsAsBytes, &res)
	// 转换商品信息
	res.Owner.Id = new_owner_id //change the owner
	res.Owner.Storename = owner.Storename
	res.Owner.Merchant = owner.Merchant
	jsonAsBytes, _ := json.Marshal(res)
	err = stub.PutState(goods_id, jsonAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("将id为" + goods_id + "的商品转移到" + owner.Storename)

	fmt.Println("--------------------end transferGoods-----------------")
	return shim.Success(nil)
}

// =============================================================================================================
/*
删除商品

*/
// =============================================================================================================
func deleteGoods(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("-------------------starting deleteGoods------------------------")

	if len(args) != 1 {
		return shim.Error("信息输入错误！")
	}

	id := args[0]
	//获取goods的信息
	_, err = get_goods(stub, id)
	if err != nil {
		fmt.Println("Failed to find goods by id " + id)
		return shim.Error(err.Error())
	}
	err = stub.DelState(id) //从coutchDB中从删除
	if err != nil {
		return shim.Error("Failed to delete ")
	}
	fmt.Println("---------------------end deleteGoods----------------------------------")
	return shim.Success(nil)
}
