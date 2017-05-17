/* global new_block,formatDate, randStr, bag, $, clear_blocks, document, WebSocket, escapeHtml, window */
var ws = {};
var bgcolors = ['whitebg', 'blackbg', 'redbg', 'greenbg', 'bluebg', 'purplebg', 'pinkbg', 'orangebg', 'yellowbg'];

// =================================================================================
// On Load
// =================================================================================
$(document).on('ready', function () {
	connect_to_server();
	//$('input[name="name"]').val('cus' + randStr(6));

	// =================================================================================
	// jQuery UI Events
	// =================================================================================
	$('#submitrequestcustomer').click(function () {
		console.log('creating customer');
		var obj = {
			type: 'requestcustomer',
			cardid: $('input[name="cusid"]').val(),
			broke: bag.setup.BROKER1NO,
			v: 'br'
		};
		if (obj.name && obj.telno) {
			console.log('creating customer, sending', obj);
			ws.send(JSON.stringify(obj));
			showHomePanel();
			// $('.colorValue').html('Color');											//reset
			// for(var i in bgcolors) $('.createball').removeClass(bgcolors[i]);		//reset
			// $('.createball').css('border', '2px dashed #fff');						//reset
		}
		return false;
	});

	$('#customerofbrokerLink').click(function () {
		showHomePanel();
	});

	$('#requestcustomerLink').click(function () {
		showRequestCustomerPanel();
	});

	//marble color picker
	$(document).on('click', '.colorInput', function () {
		$('.colorOptionsWrap').hide();											//hide any others
		$(this).parent().find('.colorOptionsWrap').show();
	});
	$(document).on('click', '.colorOption', function () {
		var color = $(this).attr('color');
		var html = '<span class="fa fa-circle colorSelected ' + color + '" color="' + color + '"></span>';

		$(this).parent().parent().find('.colorValue').html(html);
		$(this).parent().hide();

		for (var i in bgcolors) $('.createball').removeClass(bgcolors[i]);			//remove prev color
		$('.createball').css('border', '0').addClass(color + 'bg');				//set new color
	});

	//$('#cusid').on('keydown', function(e){
	$('#getinfobutton').click(function () {
		console.log('getinfobutton click');
		//if (e.which == 13) {
		ws.send(JSON.stringify({ type: 'getcus', cusid: $('input[name="cusid"]').val(), v: 'br' }));
		//}
	})


	//drag and drop marble
	$('#customerwrap, #editpane').sortable({ connectWith: '.sortable' }).disableSelection();
	$('#customerwrap').droppable({
		drop:
		function (event, ui) {
			var user = $(ui.draggable).attr('user');
			if (user.toLowerCase() != bag.setup.CUSTOMER) {
				$(ui.draggable).addClass('invalid');
				transfer($(ui.draggable).attr('id'), bag.setup.CUSTOMER);
			}
		}
	});
	// $('#brokerwrap').droppable({drop:
	// 	function( event, ui ) {
	// 		var user = $(ui.draggable).attr('user');
	// 		if(user.toLowerCase() != bag.setup.USER1){
	// 			$(ui.draggable).addClass('invalid');
	// 			transfer($(ui.draggable).attr('id'), bag.setup.USER1);
	// 		}
	// 	}
	// });
	$('#editpane').droppable({
		drop:
		function (event, ui) {
			var cardid = $(ui.draggable).attr('cardid');
			console.log('edit customer ', cardid);
			if (cardid) {
				console.log('editing customer', cardid);
				// var obj = 	{
				// 				type: 'remove',
				// 				name: cardid,
				// 				v: 'br'
				// 			};
				// ws.send(JSON.stringify(obj));
				// $(ui.draggable).fadeOut();
				// setTimeout(function(){
				// 	$(ui.draggable).remove();
				// }, 300);
				openEditCustomer(cardid);
			}
		}
	});


	function openEditCustomer(cardid) {
		//$('')
		console.log('edit ', cardid);
	}

	// =================================================================================
	// Helper Fun
	// ================================================================================
	//show admin panel page
	function showHomePanel() {
		$('#customerofbrokerPanel').fadeIn(300);
		$('#requestcustomerPanel').hide();
		//$('#createbrokerPanel').hide();

		var part = window.location.pathname.substring(0, 3);
		console.log('kycp1 - part ', part);
		window.history.pushState({}, '', part + '/customerofbroker');						//put it in url so we can f5

		console.log('getting new balls');
		setTimeout(function () {
			$('#customerofbrokewrap').html('');											//reset the panel
			$('#pendingcustomerwrap').html('');
			//ws.send(JSON.stringify({type: 'get', v: 'br'}));						//need to wait a bit
			ws.send(JSON.stringify({ type: 'getcustomerofbroke', v: 'br' }));						//need to wait a bit
			ws.send(JSON.stringify({ type: 'chainstats', v: 'br' }));
		}, 1000);
	}

	function showRequestCustomerPanel() {
		//$('#customerPanel').fadeIn(300);
		//$('#createcustomerPanel').hide();
		//$('#createbrokerPanel').hide();

		var part = window.location.pathname.substring(0, 3);
		window.history.pushState({}, '', part + '/requestcustomer');			//put it in url so we can f5

		// console.log('getting new balls');
		// setTimeout(function(){
		// 	$('#customerwrap').html('');											//reset the panel
		// 	$('#brokerwrap').html('');
		// 	ws.send(JSON.stringify({type: 'get', v: 'br'}));						//need to wait a bit
		// 	ws.send(JSON.stringify({type: 'chainstats', v: 'br'}));
		// }, 1000);
	}
});


// =================================================================================
// Socket Stuff
// =================================================================================
function connect_to_server() {
	var connected = false;

	// Redirect https requests to http so the server can handle them
	if (this.location.href.indexOf("https://") > -1) {
		this.location.href = this.location.href.replace("https://", "http://");
	}

	connect();

	function connect() {
		var wsUri = 'ws://' + document.location.hostname + ':' + document.location.port;
		console.log('Connectiong to websocket', wsUri);

		ws = new WebSocket(wsUri);
		ws.onopen = function (evt) { onOpen(evt); };
		ws.onclose = function (evt) { onClose(evt); };
		ws.onmessage = function (evt) { onMessage(evt); };
		ws.onerror = function (evt) { onError(evt); };
	}

	function onOpen(evt) {
		console.log('WS CONNECTED');
		connected = true;
		clear_blocks();
		$('#errorNotificationPanel').fadeOut();
		//ws.send(JSON.stringify({type: 'get', v: 'br'}));
		ws.send(JSON.stringify({ type: 'getcustomerofbroke', v: 'br', brokeno: bag.setup.BROKER1NO }));
		ws.send(JSON.stringify({ type: 'chainstats', v: 'br' }));
	}

	function onClose(evt) {
		console.log('WS DISCONNECTED', evt);
		connected = false;
		setTimeout(function () { connect(); }, 5000);					//try again one more time, server restarts are quick
	}

	function onMessage(msg) {
		try {
			console.log('onMessage - ', msg.data);
			var msgObj = JSON.parse(msg.data);
			if (msgObj.customer) {
				console.log('rec - cus', msgObj.msg, msgObj);
				build_customer(msgObj.customer);
			}
			else if (msgObj.broker) {
				console.log('rec- broke', msgObj.msg, msgObj);
				build_broker(msgObj.broker);
			}
			else if (msgObj.customerofbroke) {
				console.log('rec - cus of broke', msg.data)
				build_customer(msgObj.customerofbroke)
			}
			else if (msgObj.pendingcustomer) {
				console.log('rec - pending customer', msg.data)
				build_pendingcustomer(msgObj.pendingcustomer)
			}
			else if (msgObj.cusname) {
				console.log('rec -cusname', msgObj.cusname, msgObj);
				var name = msgObj.cusname;
				//console.log();
				$('#customernametext').html(name);
			}
			else if (msgObj.msg === 'chainstats') {
				console.log('rec', msgObj.msg, ': ledger blockheight', msgObj.chainstats.height, 'block', msgObj.blockstats.height);
				if (msgObj.blockstats && msgObj.blockstats.transactions) {
					var e = formatDate(msgObj.blockstats.transactions[0].timestamp.seconds * 1000, '%M/%d/%Y &nbsp;%I:%m%P');
					$('#blockdate').html('<span style="color:#fff">TIME</span>&nbsp;&nbsp;' + e + ' UTC');
					var temp = {
						id: msgObj.blockstats.height,
						blockstats: msgObj.blockstats
					};
					new_block(temp);								//send to blockchain.js
				}
			}
			else if (msgObj.msg ==='reset') {
				ws.send(JSON.stringify({ type: 'getcustomerofbroke', v: 'br', brokeno: bag.setup.BROKER1NO }));
			}
			else console.log('rec', msgObj.msg, msgObj);
		}
		catch (e) {
			console.log('ERROR', e);
		}
	}

	function onError(evt) {
		console.log('ERROR ', evt);
		if (!connected && bag.e == null) {											//don't overwrite an error message
			$('#errorName').html('Warning');
			$('#errorNoticeText').html('Waiting on the node server to open up so we can talk to the blockchain. ');
			$('#errorNoticeText').append('This app is likely still starting up. ');
			$('#errorNoticeText').append('Check the server logs if this message does not go away in 1 minute. ');
			$('#errorNotificationPanel').fadeIn();
		}
	}
}


// =================================================================================
//	UI Building
// =================================================================================
function build_customer(data) {
	var html = '';
	var colorClass = '';
	var size = 'fa-5x';

	console.log('data', data);
	data.customerid = escapeHtml(data.customerid);
	data.color = escapeHtml('blue');
	//data.user = escapeHtml(data.user);

	console.log('got a customer: ', data.color);
	if ($('#cus_' + data.customerid).length) {
		if ($('#cus_' + data.customerid).parent == 'pending') {
			$('#cus_' + data.customerid).html('')
		}
	}

	if (!$('#cus_' + data.customerid).length) {								//only populate if it doesn't exists
		//if(data.size == 16) size = 'fa-3x';
		if (data.color) colorClass = data.color.toLowerCase();

		//html += '<span id="' + data.name + '" class="fa fa-circle ' + size + ' ball ' + colorClass + ' title="' + data.name + '" user="customer">'+data.name+'</span>';
		html += '<span style="font-size: 200%" parent="allow" id="cus_' + data.customerid + '" cardid="' + data.customerid + '" class="fa fa-square fa-1x ball blue title="' + data.customerid + '">  ' + data.customerid + '  </span>'
		//if(data.user && data.user.toLowerCase() == bag.setup.USER1){
		$('#customerofbrokewrap').append(html);
		//}
		// else{
		// 	$('#user2wrap').append(html);
		// }
	}
	console.log('html after build - ', html);
	return html;
}

function build_pendingcustomer(data) {
	var html = '';
	var colorClass = '';
	var size = 'fa-5x';

	console.log('data', data);
	data.customerid = escapeHtml(data.customerid);
	data.color = escapeHtml('teal');
	//data.user = escapeHtml(data.user);

	console.log('got a pendingcustomer: ', data.color);
	if (!$('#cus_' + data.customerid).length) {								//only populate if it doesn't exists
		//if(data.size == 16) size = 'fa-3x';
		if (data.color) colorClass = data.color.toLowerCase();

		//html += '<span id="' + data.name + '" class="fa fa-circle ' + size + ' ball ' + colorClass + ' title="' + data.name + '" user="customer">'+data.name+'</span>';
		html += '<span style="font-size: 200%" parent="pending" id="cus_' + data.customerid + '" cardid="' + data.customerid + '" class="fa fa-asterisk fa-1x ball teal title="' + data.customerid + '">  ' + data.customerid + '  </span>'
		//if(data.user && data.user.toLowerCase() == bag.setup.USER1){
		$('#pendingcustomerwrap').append(html);
		//}
		// else{
		// 	$('#user2wrap').append(html);
		// }
	}
	console.log('html after build - ', html);
	return html;
}