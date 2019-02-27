package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args) //利用shim.MockStub编写单元测试代码，直接在无网络的环境中debug
	if res.Status != shim.OK {
		fmt.Println("Invoke", "failed", string(res.Message))
		t.FailNow()
	}
}
func TestExample02_Invoke(t *testing.T) {
	scc := new(GoodsChaincode)
	stub := shim.NewMockStub("GoodsChaincode", scc)
	checkInvoke(t, stub, [][]byte{[]byte("initOwner"), []byte("13"), []byte("京东总仓"), []byte("京东物流")})
	checkInvoke(t, stub, [][]byte{[]byte("initOwner"), []byte("14"), []byte("菜鸟总仓"), []byte("菜鸟物流")})
	checkInvoke(t, stub, [][]byte{[]byte("initGoods"), []byte("132456"), []byte("华为"), []byte("20100501"), []byte("北京"), []byte("华为公司"), []byte("13")})
	checkInvoke(t, stub, [][]byte{[]byte("initGoods"), []byte("111111"), []byte("苹果"), []byte("20010504"), []byte("上海"), []byte("苹果公司"), []byte("13")})
	checkInvoke(t, stub, [][]byte{[]byte("queryGoods"), []byte("132456")})
	checkInvoke(t, stub, [][]byte{[]byte("transferGoods"), []byte("132456"), []byte("14")})
	checkInvoke(t, stub, [][]byte{[]byte("queryGoods"), []byte("132456")})
	checkInvoke(t, stub, [][]byte{[]byte("deleteGoods"), []byte("132456")})
	checkInvoke(t, stub, [][]byte{[]byte("queryGoods"), []byte("132456")})
	//checkInvoke(t, stub, [][]byte{[]byte("deleteGoods"), []byte("132456")})
	checkInvoke(t, stub, [][]byte{[]byte("getHistoryForGoods"), []byte("132456")})

}
