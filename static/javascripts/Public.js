
let ws_protocol = document.location.protocol == "https:" ? "wss" : "ws"

const websocketHeartbeatJsOptions = {
	url: ws_protocol + "://"+ window.location.host +"/ws",
	pingTimeout: 15000,
	pongTimeout: 10000,
	reconnectTimeout: 2000,
	pingMsg: "heartbeat"
}

let websocketHeartbeatJs = new WebsocketHeartbeatJs(websocketHeartbeatJsOptions);

let ws = websocketHeartbeatJs;
// let ws = new WebSocket("ws://"+ window.location.host +"/ws");

// 存储发送的消息ID映射（用于撤回功能）
let sentMessages = {};

function _time(time = +new Date()) {
	var date = new Date(time + 8 * 3600 * 1000); // 增加8小时
	return date.toJSON().substr(0, 19).replace('T', ' ');
	//return date.toJSON().substr(0, 19).replace('T', ' ').replace(/-/g, '/');
}

// 检查消息是否可以撤回（2分钟内）
function canRecallMessage(msgTime) {
	let now = new Date().getTime();
	let diff = now - msgTime;
	return diff <= 2 * 60 * 1000; // 2分钟内
}

// 发送撤回消息请求
function recallMessage(msgId) {
	let send_data = JSON.stringify({
		"status": 6, // msgTypeRecall
		"data": {
			"uid": $('.room').attr('data-uid').toString(),
			"username": $('.room').attr('data-username'),
			"avatar_id": $('.room').attr('data-avatar_id'),
			"room_id": $('.room').attr('data-room_id'),
			"msg_id": msgId
		}
	});
	ws.send(send_data);
}

// 显示撤回提示
function showRecallNotification(username) {
	layer.msg(username + ' 撤回了一条消息');
}

function WebSocketConnect(userInfo,toUserInfo = null) {
	if ("WebSocket" in window) {
		//console.log("您的浏览器支持 WebSocket!");

		if ( userInfo.uid <= 0 )
		{
			alert('参数错误，请刷新页面重试');return false;
		}

		// 打开一个 web socket
		// let ws = new WebSocket("ws://127.0.0.1:8322/ws");

		let send_data = JSON.stringify({
			"status": toUserInfo ? 5 : 1,
			"data": {
				"uid": userInfo.uid.toString(),
				"room_id": userInfo.room_id,
				"avatar_id": userInfo.avatar_id,
				"username": userInfo.username,
				"to_user": toUserInfo
			}
		})

		ws.onopen = function () {
			// layer.msg("websocket 连接已建立");
			chat_info.html(chat_info.html() +
				'<li class="systeminfo" > <span>' +
				"✅ websocket 连接已建立 " +
				'</span></li>');
			ws.send(send_data);
			//console.log("send_data 发送数据", send_data)
			toLow();
		};

		// if ( toUserInfo )
		// {
		// 	let to_user_send_data = JSON.stringify({
		// 		"status": toUserInfo ? 5 : 1,
		// 		"data": {
		// 			"uid": toUserInfo.uid,
		// 			"room_id": toUserInfo.room_id,
		// 			"avatar_id": toUserInfo.avatar_id,
		// 			"username": toUserInfo.username,
		// 			"to_user": toUserInfo
		// 		}
		// 	})
		// 	ws.onopen = function () {
		// 		ws.send(to_user_send_data);
		// 		console.log("to_user_send_data 发送数据", to_user_send_data)
		// 	};
		// }


		let chat_info = $('.main .chat_info')
		let isServeClose = 0;

		ws.onmessage = function (evt) {
			var received_msg = JSON.parse(evt.data);

			// let myDate = new Date();
			// let time = myDate.toLocaleDateString() + myDate.toLocaleTimeString()
			let time = _time(received_msg.data.time)

			switch(received_msg.status)
			{
				case 1:
					chat_info.html(chat_info.html() +
						'<li class="systeminfo"> <span>' +
						"【" +
						received_msg.data.username +
						"】" +
						time +
						" 加入了房间" +
						'</span></li>');
					break;
				case 2:
					chat_info.html(chat_info.html() +
						'<li class="systeminfo"> <span>' +
						"【" +
						received_msg.data.username +
						"】" +
						time +
						" 离开了房间" +
						'</span></li>');
					break;
				case 3:
				if ( received_msg.data.uid != userInfo.uid && !isPrivateChat())
				{
					let msgHtml = '<li class="left" data-msg-id="' + received_msg.data.msg_id + '"><img src="/static/images/user/' +
						received_msg.data.avatar_id +
						'.png" alt=""><b>' +
						received_msg.data.username +
						'</b><i>' +
						time +
						'</i><div class="aaa">' +
						received_msg.data.content +
						'</div></li>';
					chat_info.html(chat_info.html() + msgHtml);
				} else if ( received_msg.data.uid == userInfo.uid && !isPrivateChat() ) {
					// 自己发送的消息，关联服务器返回的消息ID
					if (received_msg.data.msg_id && window.lastTempMsgId) {
						$('#' + window.lastTempMsgId).attr('data-msg-id', received_msg.data.msg_id);
						sentMessages[received_msg.data.msg_id] = {
							time: received_msg.data.time,
							content: received_msg.data.content,
							elementId: window.lastTempMsgId
						};
						window.lastTempMsgId = null;
					}
				}
				break;
			case 6:
				// 撤回消息通知
				if (received_msg.data.recall_msg) {
					// 撤回失败，显示错误信息
					layer.msg(received_msg.data.recall_msg);
				} else {
					// 撤回成功，更新界面
					let msgId = received_msg.data.msg_id;
					let msgElement = $('li[data-msg-id="' + msgId + '"]');
					if (msgElement.length > 0) {
						msgElement.find('.aaa').html('<span style="color: #999; font-style: italic;">消息已被撤回</span>');
						msgElement.find('.aaa').addClass('recalled');
						// 移除右键菜单（如果存在）
						msgElement.removeClass('has-recall-menu');
					}
					// 如果是别人撤回的消息，显示提示
					if (received_msg.data.uid != userInfo.uid.toString()) {
						showRecallNotification(received_msg.data.username);
					}
				}
				break;
				case -1:
					ws.close() // 主动close掉
					isServeClose = 1
					console.log("client 连接已关闭...");
					break;
				case 4:
					$('.popover-title').html('在线用户 '+ received_msg.data.count +' 人')

					$.each(received_msg.data.list,function (index, value) {

						if ( received_msg.data.uid == value.uid )
						{
							// 禁止点击
							$('.ul-user-list').html($('.ul-user-list').html() +
								'<li  style="pointer-events: none;" class="li-user-item" data-uid='+ value.uid +' data-username='+ value.username +' data-room_id='+ value.room_id +' data-avatar_id='+ value.avatar_id +'  ><img src="/static/images/user/' +
								value.avatar_id +
								'.png" alt=""><b>' + " " +
								value.username +
								'</b>' +
								'</li>'
							)
						}else{
							$('.ul-user-list').html($('.ul-user-list').html() +
								'<li  class="li-user-item" data-uid='+ value.uid +' data-username='+ value.username +' data-room_id='+ value.room_id +' data-avatar_id='+ value.avatar_id +'  ><img src="/static/images/user/' +
								value.avatar_id +
								'.png" alt=""><b>' + " " +
								value.username +
								'</b>' +
								'</li>'
							)
						}

					})
					//console.log("在线用户",received_msg);
					break;
				case 5:
				// 私聊通知
				if (!isPrivateChat())
				{
					layer.msg(received_msg.data.username+'：'+ received_msg.data.content);
				} else {
					// 在私聊窗口中，保存消息ID
					if (received_msg.data.uid == userInfo.uid.toString() && received_msg.data.msg_id) {
						sentMessages[received_msg.data.msg_id] = {
							time: received_msg.data.time,
							content: received_msg.data.content
						};
					}
				}
				break;
				default:
			}
			// console.log("数据已接收...", received_msg);

            if ( !(received_msg.data === "heartbeat ok") ){
                // 滚动条滚到最下面
                toLow();
            }

		};

		ws.onclose = function (evt) {
			// 关闭 websocket
			if ( isServeClose === 1 ){
				chat_info.html(chat_info.html() +
					'<li class="systeminfo"> <span>' +
					"❌ 与服务器连接断开，请检查是否在浏览器中打开了多个聊天界面" +
					'</span></li>');
			}else{
				chat_info.html(chat_info.html() +
					'<li class="systeminfo"> <span>' +
					"❌ 与服务器连接断开，正在尝试重新连接，请稍后..." +
					'</span></li>');
			}
			// let c = ws.close() // 主动close掉
			console.log("serve 连接已关闭... " + _time(),evt);
			// console.log(c);
			toLow();
		};
		
		ws.onerror = function (evt) {
			// ws.close()
			console.log("触发 onerror",evt)
		}

		ws.onreconnect = (e) => {
			console.log('reconnecting...');
		}
		
	} else {
		// 浏览器不支持 WebSocket
		alert("您的浏览器不支持 WebSocket!");
	}
}

$(document).ready(function(){
// ------------------------选择聊天室页面-----------------------------------------------

	// 在页面即将卸载之前关闭WebSocket连接
	window.addEventListener("beforeunload", function() {
		console.log("beforeunload close");
		ws.close();
	});
	// 用户信息提交

	$('#userinfo_sub').click(function(event) {
		var userName = $('.rooms .user_name input').val(); // 用户昵称
		var userPortrait = $('.rooms .user_portrait img').attr('portrait_id'); // 用户头像id
		if(userName=='') { // 如果不填用户昵称，就是以前的昵称
			userName = $('.rooms .user_name input').attr('placeholder');
		}


		// 下面是测试用的代码


		$('.userinfo a b').text(userName); // 修改标题栏的用户昵称
		$('.rooms .user_name input').val(''); // 昵称输入框清空
		$('.rooms .user_name input').attr('placeholder', userName); // 昵称输入框默认显示用户昵称
		$('.topnavlist .popover').not($(this).next('.popover')).removeClass('show'); // 关掉用户面板
		$('.clapboard').addClass('hidden'); // 关掉模糊背景
	});

	// 设置主题


	$('.theme img').click(function(event) {
		var theme_id = $(this).attr('theme_id');
		$('.clapboard').click(); // 关掉用户模糊背景




		// 下面是测试用的代码


		$('body').css('background-image', 'url(images/theme/' + theme_id + '_bg.jpg)'); // 设置背景
	});






















// --------------------聊天室内页面----------------------------------------------------

	// 获取在线用户列表
	$(document).on('click', '.a-user-list', function(e) {
		$('.ul-user-list').html('')
		let send_data = JSON.stringify({
			"status": 4,
			"data": {
				"uid": $('.room').attr('data-uid').toString(),
				"username": $('.room').attr('data-username'),
				"avatar_id": $('.room').attr('data-avatar_id'),
				"room_id": $('.room').attr('data-room_id'),
			}
		})
		ws.send(send_data);
	})

	// 发送图片

	$('.imgFileBtn').change(function(event) {

		var formData = new FormData();
		formData.append('file', $(this)[0].files[0]);
		$.ajax({
			url: '/img-kr-upload',
			type: 'POST',
			beforeSend: function (xhr) {
				// 在请求发送之前执行的代码
				console.log('请求即将发送');

				// 在请求发送之前调用 layer 的加载动画
				var index = layer.load(1, { // 1 是加载动画的样式，layer 提供了多种样式
					shade: [0.5, '#000'], // 遮罩层颜色和透明度
					time: 25000, // 最大显示时间（毫秒），超过此时间自动关闭
					success: function(layero, index) {
						// 加载动画加载完成时的回调
						console.log('加载动画已显示');
					}
				});
				// 将加载动画的索引存储到全局变量或闭包中，方便后续关闭
				window.layerIndex = index;

			},
			cache: false,
			data: formData,
			processData: false,
			contentType: false
		}).done(function(res) {
			console.log(res)

			var str = '<img src="' + res.data.url +'" />'

			let to_uid = "0"
			let status = 3
			if (isPrivateChat()) {
				// 私聊
				to_uid = getQueryVariable("uid")
				status = 5
			}

			let tempMsgId = sends_message($('.room').attr('data-username'), $('.room').attr('data-avatar_id'), str); // sends_message(昵称,头像id,聊天内容);

			let send_data = JSON.stringify({
				"status": status,
				"data": {
					"uid": $('.room').attr('data-uid').toString(),
					"username": $('.room').attr('data-username'),
					"avatar_id": $('.room').attr('data-avatar_id'),
					"room_id": $('.room').attr('data-room_id'),
					"image_url": res.data.url,
					"content": str,
					"to_uid" : to_uid,
				}
			})

			// 保存临时ID到全局变量，等待服务器返回消息ID
			if (tempMsgId) {
				window.lastTempMsgId = tempMsgId;
			}

			console.log("send_data",send_data)
			ws.send(send_data);


			// 滚动条滚到最下面
			toLow();

			// 解决input上传文件选择同一文件change事件不生效
			event.target.value=''

			layer.close(window.layerIndex);
		}).fail(function(res) {});



	});

	// 发送消息
	
	//$('.text input').focus();
	$("#emojionearea2")[0].emojioneArea.setFocus()
	$('#subxx').click(function(event) {
		//var str = $('.text input').val(); // 获取聊天内容
		var str = $("#emojionearea2")[0].emojioneArea.getText() // 获取聊天内容
		str = str.replace(/\</g,'&lt;');
		str = str.replace(/\>/g,'&gt;');
		str = str.replace(/\n/g,'<br/>');
		str = str.replace(/\[em_([0-9]*)\]/g,'<img src="images/face/$1.gif" alt="" />');

		if($.trim(str)!=='') {

			let to_uid = "0"
			let status = 3
			if (isPrivateChat()) {
				// 私聊
				to_uid = getQueryVariable("uid")
				status = 5
			}


			let tempMsgId = sends_message($('.room').attr('data-username'), $('.room').attr('data-avatar_id'), str); // sends_message(昵称,头像id,聊天内容);

			let send_data = JSON.stringify({
				"status": status,
				"data": {
					"uid": $('.room').attr('data-uid').toString(),
					"username": $('.room').attr('data-username'),
					"avatar_id": $('.room').attr('data-avatar_id'),
					"room_id": $('.room').attr('data-room_id'),
					"content": str,
					"image_url" : "",
					"to_uid" : to_uid,
				}
			})

			// 保存临时ID到全局变量，等待服务器返回消息ID
			if (tempMsgId) {
				window.lastTempMsgId = tempMsgId;
			}

			ws.send(send_data);

			// 滚动条滚到最下面
			toLow();

		}

		$("#emojionearea2")[0].emojioneArea.setText("")
		$("#emojionearea2")[0].emojioneArea.setFocus()
	});





























// -----下边的代码不用管---------------------------------------



	jQuery('.scrollbar-macosx').scrollbar();
	$('.topnavlist li a.a-user-list').click(function(event) {
		$('.topnavlist .popover').not($(this).next('.popover')).removeClass('show');
		$(this).next('.popover').toggleClass('show');
		if($(this).next('.popover').attr('class')!='popover fade bottom in') {
			$('.clapboard').removeClass('hidden');
		}else{
			$('.clapboard').click();
		}
	});
	$('.clapboard').click(function(event) {
		$('.topnavlist .popover').removeClass('show');
		$(this).addClass('hidden');
		$('.user_portrait img').attr('portrait_id', $('.user_portrait img').attr('ptimg'));
		$('.user_portrait img').attr('src', '/static/images/user/' + $('.user_portrait img').attr('ptimg') + '.png');
		$('.select_portrait img').removeClass('t');
		$('.select_portrait img').eq($('.user_portrait img').attr('ptimg')-1).addClass('t');
		$('.rooms .user_name input').val('');
	});
	$('.select_portrait img').hover(function() {
		var portrait_id = $(this).attr('portrait_id');
		$('.user_portrait img').attr('src', '/static/images/user/' + portrait_id + '.png');
	}, function() {
		var t_id = $('.user_portrait img').attr('portrait_id');
		$('.user_portrait img').attr('src', '/static/images/user/' + t_id + '.png');
	});
	$('.select_portrait img').click(function(event) {
		var portrait_id = $(this).attr('portrait_id');
		$('.user_portrait img').attr('portrait_id', portrait_id);
		$('.select_portrait img').removeClass('t');
		$(this).addClass('t');
	});
	$('.face_btn,.faces').hover(function() {
		$('.faces').addClass('show');
	}, function() {
		$('.faces').removeClass('show');
	});
	$('.faces img').click(function(event) {
		if($(this).attr('alt')!='') {
			$('.text input').val($('.text input').val() + '[em_' + $(this).attr('alt') + ']');
		}
		$('.faces').removeClass('show');
		$('.text input').focus();
	});
	$('.imgFileico').click(function(event) {
		$('.imgFileBtn').click();
	});
	// 存储当前发送消息的临时ID（用于关联服务器返回的消息ID）
	let currentMsgTempId = 0;

	function sends_message (userName, userPortrait, message) {
		if(message!='') {

			let myDate = new Date();
			let time = myDate.toLocaleDateString() + myDate.toLocaleTimeString()
			currentMsgTempId++;
			let tempId = 'temp_' + currentMsgTempId;
			let msgHtml = '<li class="right" id="' + tempId + '" data-msg-id=""><img src="/static/images/user/' + userPortrait + '.png" alt=""><b>' + userName + '</b><i>'+ time +'</i><div class="aaa">' + message  +'</div></li>';
			$('.main .chat_info').html($('.main .chat_info').html() + msgHtml);
			return tempId;
		}
		return null;
	}
	$('.text input').keypress(function(e) {
		if (e.which == 13){
			$('#subxx').click();
		}
	});


	function replaceImg() {
		$(".load-img").each(function () {
			let realImgUrl = $(this).attr("data-src");
			if (realImgUrl !== "") {
				$(this).attr("src", $(this).attr("data-src"))
			}
		});
	}
	setTimeout(replaceImg, 1500);

	// 右键菜单 - 撤回消息功能
	$(document).on('contextmenu', '.chat_info li.right', function(e) {
		e.preventDefault();
		let $this = $(this);
		let msgId = $this.attr('data-msg-id');

		// 检查是否有消息ID
		if (!msgId) {
			return false;
		}

		// 检查消息是否在2分钟内
		let msgData = sentMessages[msgId];
		if (!msgData || !canRecallMessage(msgData.time)) {
			layer.msg('消息发送超过2分钟，无法撤回');
			return false;
		}

		// 检查消息是否已被撤回
		if ($this.find('.aaa').hasClass('recalled')) {
			return false;
		}

		// 创建右键菜单
		let menuHtml = '<div class="recall-context-menu" style="position: fixed; z-index: 9999; background: #fff; border: 1px solid #ccc; border-radius: 4px; padding: 5px 0; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">' +
			'<div class="recall-menu-item" style="padding: 8px 20px; cursor: pointer; font-size: 14px; color: #333;" data-msg-id="' + msgId + '">撤回消息</div>' +
			'</div>';

		// 移除已存在的菜单
		$('.recall-context-menu').remove();

		// 添加新菜单
		$('body').append(menuHtml);

		// 设置菜单位置
		let menuWidth = 100;
		let menuHeight = 40;
		let posX = e.clientX;
		let posY = e.clientY;

		// 边界检查
		if (posX + menuWidth > $(window).width()) {
			posX = $(window).width() - menuWidth - 10;
		}
		if (posY + menuHeight > $(window).height()) {
			posY = $(window).height() - menuHeight - 10;
		}

		$('.recall-context-menu').css({
			left: posX,
			top: posY
		});

		return false;
	});

	// 点击撤回菜单项
	$(document).on('click', '.recall-menu-item', function() {
		let msgId = $(this).attr('data-msg-id');
		if (msgId) {
			recallMessage(parseInt(msgId));
		}
		$('.recall-context-menu').remove();
	});

	// 点击其他地方关闭菜单
	$(document).on('click', function() {
		$('.recall-context-menu').remove();
	});

});

function getQueryVariable(variable)
{
	var query = window.location.search.substring(1);
	var vars = query.split("&");
	for (var i=0;i<vars.length;i++) {
		var pair = vars[i].split("=");
		if(pair[0] == variable){return pair[1];}
	}
	return(false);
}

function isPrivateChat()
{
	return window.location.href.search('private-chat') > 0
}

function toLow() {
	$('.scrollbar-macosx.scroll-content.scroll-scrolly_visible').animate({
		scrollTop: $('.scrollbar-macosx.scroll-content.scroll-scrolly_visible').prop('scrollHeight')
	}, 500);
}


