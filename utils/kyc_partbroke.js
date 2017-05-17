// ==================================
// KYC Part 1 - incoming messages, look for type
// ==================================
var ibc = {};
var chaincode = {};
var async = require('async');

module.exports.setup = function(sdk, cc){
	ibc = sdk;
	chaincode = cc;
};

module.exports.process_msg = function(ws, data){
	console.log('kyc part broke - process_msg');
	if(data.v === 'br'){
		if(data.type == 'getcus'){
			console.log('get cus msg');
			chaincode.query.read([data.cusid], cb_got_cus);
		}
		else if(data.type == 'getcustomerofbroke'){
			console.log('get customer of broker msg');
			chaincode.query.readbroker([data.brokeno], cb_got_broker);
		}
		else if(data.type == 'chainstats'){
			console.log('chainstats msg');
			ibc.chain_stats(cb_chainstats);
		}
	}

	function cb_got_broker(e, brokeinfo){
		if(e!=null) console.log('[ws error] did not get broke', e);
		else{
			try{
				console.log('cb got broker', brokeinfo);
				var broker = JSON.parse(brokeinfo);
				if (broker.allowcustomer){
					var keys = Object.keys(broker.allowcustomer);
					var concurrency = 1;

					async.eachLimit(keys, concurrency, function(key, cb) {
						console.log('!', broker.allowcustomer[key]);
						chaincode.query.readcustomergid([broker.allowcustomer[key]], function(e, customer) {
							if(e != null) console.log('[ws error] did not get customer:', e);
							else {
								console.log('read !!!! ', JSON.parse(customer));
								if (customer) sendMsg({msg: 'customerofbroke', e:e, customerofbroke: JSON.parse(customer)});
								cb(null);
							}
						});
					}, function() {
						sendMsg({msg: 'action', e:e, status: 'finished'});
					});
				}
				
				if (broker.pendingcustomer){
					var keys = Object.keys(broker.pendingcustomer);
					var concurrency = 1;

					async.eachLimit(keys, concurrency, function(key, cb) {
						console.log('!', broker.pendingcustomer[key]);
						chaincode.query.readcustomergid([broker.pendingcustomer[key]], function(e, customer) {
							if(e != null) console.log('[ws error] did not get customer:', e);
							else {
								console.log('read !!!! ', JSON.parse(customer));
								if (customer) sendMsg({msg: 'pendingcustomer', e:e, pendingcustomer: JSON.parse(customer)});
								cb(null);
							}
						});
					}, function() {
						sendMsg({msg: 'action', e:e, status: 'finished'});
					});
				}
			}
			catch(e){
				console.log('[ws error] could not parse response', e);
			}
		}
	}
	function cb_got_cus(e, customer){
		if(e!=null)console.log('[ws error] did not get customer:', e);
		else{
			try {
				console.log('customer', customer);
				var json = JSON.parse(customer);
				if (json.name){
					console.log('json.name', json.name);
					sendMsg({msg: 'cusname', e: e, cusname: json.name});
				}
				sendMsg({msg: 'action', e:e, status: 'finished'});
			}catch(e){
				console.log('[ws error] could not parse response', e);
			}
		}
	}
	
	function cb_invoked(e, a){
		console.log('response: ', e, a);
	}
	
	//call back for getting the blockchain stats, lets get the block stats now
	function cb_chainstats(e, chain_stats){
		if(chain_stats && chain_stats.height){
			chain_stats.height = chain_stats.height - 1;								//its 1 higher than actual height
			var list = [];
			for(var i = chain_stats.height; i >= 1; i--){								//create a list of heights we need
				list.push(i);
				if(list.length >= 8) break;
			}
			list.reverse();																//flip it so order is correct in UI
			async.eachLimit(list, 1, function(block_height, cb) {						//iter through each one, and send it
				ibc.block_stats(block_height, function(e, stats){
					if(e == null){
						stats.height = block_height;
						sendMsg({msg: 'chainstats', e: e, chainstats: chain_stats, blockstats: stats});
					}
					cb(null);
				});
			}, function() {
			});
		}
	}
	
	//send a message, socket might be closed...
	function sendMsg(json){
		if(ws){
			try{
				ws.send(JSON.stringify(json));
			}
			catch(e){
				console.log('[ws error] could not send msg', e);
			}
		}
	}
};