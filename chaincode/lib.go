package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ============================================================================================================================
// Get Owner - get the owner asset from ledger
// ============================================================================================================================
func get_owner(stub shim.ChaincodeStubInterface, id string) (Owner, error) {
	var owner Owner
	ownerAsBytes, err := stub.GetState(id)
	if err != nil {
		return owner, errors.New("Failed to get owner - " + id)
	}
	json.Unmarshal(ownerAsBytes, &owner)

	if len(owner.Storename) == 0 { //test if owner is actually here or just nil
		return owner, errors.New("Owner does not exist - " + id + ", '" + owner.Storename + "' '" + owner.Merchant + "'")
	}

	return owner, nil
}

// ============================================================================================================================
// Get goods - get a goods asset from ledger
// ============================================================================================================================
func get_goods(stub shim.ChaincodeStubInterface, id string) (Goods, error) {
	var goods Goods
	goodsAsBytes, err := stub.GetState(id) //从coutchdb数据库中获取信息
	if err != nil {
		return goods, errors.New("Failed to find goods - " + id)
	}
	json.Unmarshal(goodsAsBytes, &goods) //将一个 JSON 字符串转换到相应的数据结构

	if goods.Id != id { //test if goods is actually here or just nil
		return goods, errors.New("Goods does not exist - " + id)
	}

	return goods, nil
}
