'use strict';
var net =require("net")//引入socket
var Fabric_Client = require('fabric-client');//fabric-sdk-node 中的依赖包之一
var path = require('path');//Node.js path 模块提供了一些用于处理文件路径的小工具
var util = require('util');//util 是一个Node.js 核心模块，提供常用函数的集合，用于弥补核心JavaScript 的功能 过于精简的不足
var os = require('os');//Node.js os 模块提供了一些基本的系统操作函数
var iconv=require("iconv-lite");//引入字符编码包
var fs = require('fs');
//定义收到的function和data
var func;
var recvdata;
//定义发送的数据
var response_data;

//设置channel
var channelName = 'mychannel';
//配置fabric基础
var Fabric_Client = require('fabric-client');

var client=new Fabric_Client();

//设置要连接的Channel
var channel = client.newChannel(channelName);

//设置要连接的Peer
var peer = client.newPeer( 'grpc://localhost:7051');
channel.addPeer(peer);
//设置的order
var order = client.newOrderer('grpc://192.168.242.218:7050');
channel.addOrderer(order);




//------------------------socket
//-----------server------------
var listenPort=8888;//监听端口
var server=net.createServer(function(socket){//创建socket服务器
    console.log('connect: '+socket.remoteAddress + ':'+ socket.remotePort);
    socket.setEncoding("UTF8");
    //接收数据
    socket.on('data',function(data){
      
       var recv = iconv.decode(data,'utf8');
	   recvdata =JSON.parse(recv);
	   func=recvdata.function;
	   console.log('client send: '+recv);
	   if(func=="initGoods"||func=="initOwner"||func=="transferGoods"||func=="deleteGoods"){
		    invoke();	   
	   }
	   else if(func=="queryGoods"||func=="getHistoryForGoods"){
			query();	
            setTimeout(function(){socket.write(response_data);},100);//延迟1s发送数据给客户端
       }
	  
    });
    
    //数据错误事件
    socket.on('error',function(err){
        console.log('socket error: '+err);
        socket.end();//发送FIN给客户端
    });
    //客户端关闭事件
    socket.on('close',function(data){
        console.log('client closed!');
    });

}).listen(listenPort);
//服务器监听事件
server.on('listening',function(){
    console.log('server listening: '+ server.address().port);
});
//服务器错误事件
server.on('error',function(err){
    console.log('server error: '+err);
});

//invoke操作链码
var invoke = function invoke(){
	//创建一个Client
	Fabric_Client.newDefaultKeyValueStore({ path: '/tmp/xx/' }).then((state_store) => {
	    
	    client.setStateStore(state_store);

	    //设置用户信息    
	    var userOpt = {
	        username: 'Admin@org1.example.com',
	        mspid: 'Org1MSP',
	        cryptoContent: { 
	            privateKey: '/opt/gopath/src/github.com/hyperledger/fabric/examples/e2e_cli/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/5e69b75b50342be13b7cb34ae42befe6284a3994b93f437855e64a53369dddc9_sk',
	            signedCert: '/opt/gopath/src/github.com/hyperledger/fabric/examples/e2e_cli/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem'
	        }
	    }

	    return client.createUser(userOpt);

	}).then((user)=>{


	    // get a transaction id object based on the current user assigned to fabric client
		tx_id = client.newTransactionID();
		console.log("Assigning transaction_id: ", tx_id._transaction_id);
    	/*
		request内容为cli命令行中的内容
	
		*/                                                                                               
		// createCar chaincode function - requires 5 args, ex: args: ['CAR12', 'Honda', 'Accord', 'Black', 'Tom'],
		// changeCarOwner chaincode function - requires 2 args , ex: args: ['CAR10', 'Dave'],
		// must send the proposal to endorsing peers
		if(func=="initGoods"){
			var request = {
				//targets: let default to the peer assigned to the client
				chaincodeId: 'fabric_tracing',
				fcn: func,
				args: [recvdata.goods_id,recvdata.goods_name,recvdata.borntime,recvdata.bornspace,recvdata.manufacturer,recvdata.owner_id],
				chainId: channelName,
				txId: tx_id
			};
		}
		else if(func=="initOwner"){
			var request = {
				//targets: let default to the peer assigned to the client
				chaincodeId: 'fabric_tracing',
				fcn: func,
				args: [recvdata.store_id,recvdata.store_name,recvdata.merchant],
				chainId: channelName,
				txId: tx_id
			};
		}
		else if(func=="transferGoods"){
			var request = {
				//targets: let default to the peer assigned to the client
				chaincodeId: 'fabric_tracing',
				fcn: func,
				args: [recvdata.goods_id,recvdata.new_owner_id],
				chainId: channelName,
				txId: tx_id
			};
		}
		else if(func=="deleteGoods"){
			var request = {
				//targets: let default to the peer assigned to the client
				chaincodeId: 'fabric_tracing',
				fcn: func,
				args: [recvdata.goods_id],
				chainId: channelName,
				txId: tx_id
			};
		}
		else{
			console.log("----没有此方法");
		}

		// send the transaction proposal to the peers
		return channel.sendTransactionProposal(request);//交易流程第二步---proposal
		}).then((results) => {
			var proposalResponses = results[0];
			var proposal = results[1];
			let isProposalGood = false;
			if (proposalResponses && proposalResponses[0].response &&
				proposalResponses[0].response.status === 200) {
					isProposalGood = true;
					console.log('Transaction proposal was good');
				} else {
					console.error('Transaction proposal was bad');
				}
			if (isProposalGood) {
				console.log(util.format(
					'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s"',
					proposalResponses[0].response.status, proposalResponses[0].response.message));

				// build up the request for the orderer to have the transaction committed
				var request = {
					proposalResponses: proposalResponses,
					proposal: proposal
				};

				// set the transaction listener and set a timeout of 30 sec
				// if the transaction did not get committed within the timeout period,
				// report a TIMEOUT status
				var transaction_id_string = tx_id.getTransactionID(); //Get the transaction ID string to be used by the event processing
				var promises = [];

				var sendPromise = channel.sendTransaction(request);
				promises.push(sendPromise); //we want the send transaction first, so that we know where to check status

				// get an eventhub once the fabric client has a user assigned. The user
				// is required bacause the event registration must be signed
				let event_hub = channel.newChannelEventHub(peer);

				// using resolve the promise so that result status may be processed
				// under the then clause rather than having the catch clause process
				// the status
				let txPromise = new Promise((resolve, reject) => {
					let handle = setTimeout(() => {
						event_hub.unregisterTxEvent(transaction_id_string);
						event_hub.disconnect();
						resolve({event_status : 'TIMEOUT'}); //we could use reject(new Error('Trnasaction did not complete within 30 seconds'));
					}, 3000);
					event_hub.registerTxEvent(transaction_id_string, (tx, code) => {
						// this is the callback for transaction event status
						// first some clean up of event listener
						clearTimeout(handle);

						// now let the application know what happened
						var return_status = {event_status : code, tx_id : transaction_id_string};
						if (code !== 'VALID') {
							console.error('The transaction was invalid, code = ' + code);
							resolve(return_status); // we could use reject(new Error('Problem with the tranaction, event status ::'+code));
						} else {
							console.log('The transaction has been committed on peer ' + event_hub.getPeerAddr());
							resolve(return_status);
						}
					}, (err) => {
						//this is the callback if something goes wrong with the event registration or processing
						reject(new Error('There was a problem with the eventhub ::'+err));
					},
						{disconnect: true} //disconnect when complete
					);
					event_hub.connect();

				});
				promises.push(txPromise);

				return Promise.all(promises);
			} else {
				console.error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
				throw new Error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
			}
		}).then((results) => {
			console.log('Send transaction promise and event listener promise have completed');
			// check the results in the order the promises were added to the promise all list
			if (results && results[0] && results[0].status === 'SUCCESS') {
				console.log('Successfully sent transaction to the orderer.');
			} else {
				console.error('Failed to order the transaction. Error code: ' + results[0].status);
			}

			if(results && results[1] && results[1].event_status === 'VALID') {
				console.log('Successfully committed the change to the ledger by the peer');
			} else {
				console.log('Transaction failed to be committed to the ledger due to ::'+results[1].event_status);
			}
		}).catch((err) => {
			console.error('Failed to invoke successfully :: ' + err);
		});
}



var query= function query(){	
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

	    // //设置要连接的Channel
	    // var channel = client.newChannel(channelName);

	    // //设置要连接的Peer
	    // var peer = client.newPeer( 'grpc://localhost:7051');

	    // channel.addPeer(peer);

	    //调用chaincode
	    const request = {
	        chaincodeId: 'fabric_tracing',   //chaincode名称
	        fcn: func,          //调用的函数名
	        args: [recvdata.goods_id]         //参数
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
	}