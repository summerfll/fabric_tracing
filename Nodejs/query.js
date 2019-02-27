
var fs = require('fs');
var Fabric_Client = require('fabric-client');

//创建一个Client
Fabric_Client.newDefaultKeyValueStore({ path: '/tmp/xx/' }).then((state_store) => {
    client=new Fabric_Client();
    client.setStateStore(state_store)

    //设置用户信息    
    var userOpt = {
        username: 'Admin@org1.example.com',
        mspid: 'Org1MSP',
        cryptoContent: { 
            privateKey: '/opt/gopath/src/github.com/hyperledger/fabric/examples/e2e_cli/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/5e69b75b50342be13b7cb34ae42befe6284a3994b93f437855e64a53369dddc9_sk',
            signedCert: '/opt/gopath/src/github.com/hyperledger/fabric/examples/e2e_cli/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem'
        }
    }

    return client.createUser(userOpt)

}).then((user)=>{

    //设置要连接的Channel
    var channel = client.newChannel('mychannel');

    //设置要连接的Peer
    var peer = client.newPeer( 'grpc://localhost:7051');

    channel.addPeer(peer);

    //调用chaincode
    const request = {
        chaincodeId: 'fabric_tracing',   //chaincode名称
        fcn: 'getHistoryForGoods',          //调用的函数名
        args: ['123']         //参数
    };

    // send the query proposal to the peer
    return channel.queryByChaincode(request);
	}).then((query_responses) => {
		console.log("Query has completed, checking results");
		// query_responses could have more than one  results if there multiple peers were used as targets
		if (query_responses && query_responses.length == 1) {
			if (query_responses[0] instanceof Error) {
				console.error("error from query = ", query_responses[0]);
			} else {
				console.log("Response is ", query_responses[0].toString());	
				var sendToMysql = JSON.parse(query_responses[0]);
				console.log("Response is ", sendToMysql.TxId);
				console.log("Response is ", sendToMysql.Value);
				console.log("Response is ", sendToMysql.Timestamp);
				response_data=query_responses[0].toString();
			}
		} else {
			console.log("No payloads were returned from query");
		}
	}).catch((err) => {
		console.error('Failed to query successfully :: ' + err);
	});