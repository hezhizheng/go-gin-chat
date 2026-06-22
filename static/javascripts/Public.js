
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

function _time(time = +new Date()) {
	var date = new Date(time + 8 * 3600 * 1000); // 增加8小时
	return date.toJSON().substr(0, 19).replace('T', ' ');
	//return date.toJSON().substr(0, 19).replace('T', ' ').replace(/-/g, '/');
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
				// 收到别人的消息
				var msgContent = received_msg.data.is_recalled == 1 
					? '<span style="color:#999;font-style:italic;">该消息已撤回</span>'
					: received_msg.data.content;
				chat_info.html(chat_info.html() +
					'<li class="left" data-msg-id="' + received_msg.data.msg_id + '"><img src="/static/images/user/' +
					received_msg.data.avatar_id +
					'.png" alt=""><b>' +
					received_msg.data.username +
					'</b><i>' +
					time +
					'</i><div class="aaa">' +
					msgContent +
					'</div></li>');
			} else if (received_msg.data.uid == userInfo.uid && !isPrivateChat()) {
				// 收到自己发送的消息回执，更新消息ID
				if (received_msg.data.msg_id && window.lastSentTempId) {
					var tempLi = $('#' + window.lastSentTempId);
					if (tempLi.length > 0) {
						tempLi.attr('data-msg-id', received_msg.data.msg_id);
						tempLi.attr('id', '');
						window.lastSentTempId = null;
					}
				}
			}
			break;
			case 6:
				// 撤回消息通知
				var msgLi = $('li[data-msg-id="' + received_msg.data.msg_id + '"]');
				if (msgLi.length > 0) {
					msgLi.find('.aaa').html('<span style="color:#999;font-style:italic;">该消息已撤回</span>');
				}
				// 显示系统提示
				chat_info.html(chat_info.html() +
					'<li class="systeminfo"> <span>' +
					received_msg.data.username +
					' 撤回了一条消息' +
					'</span></li>');
				break;
			case -2:
				// 撤回失败
				layer.msg('撤回失败：' + received_msg.data.content);
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

			sends_message($('.room').attr('data-username'), $('.room').attr('data-avatar_id'), str); // sends_message(昵称,头像id,聊天内容);

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


			sends_message($('.room').attr('data-username'), $('.room').attr('data-avatar_id'), str); // sends_message(昵称,头像id,聊天内容);

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

			ws.send(send_data);

			// 滚动条滚到最下面
			toLow();

			// 保存当前发送的消息内容，用于接收服务器返回的消息ID后更新DOM
			window.lastSentMessage = str;

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
	function sends_message (userName, userPortrait, message) {
		if(message!='') {

			let myDate = new Date();
			let time = myDate.toLocaleDateString() + myDate.toLocaleTimeString()
			// 创建唯一ID用于关联消息
			var tempMsgId = 'temp_' + Date.now();
			$('.main .chat_info').html($('.main .chat_info').html() + '<li class="right" id="' + tempMsgId + '" data-msg-id=""><img src="/static/images/user/' + userPortrait + '.png" alt=""><b>' + userName + '</b><i>'+ time +'</i><div class="aaa">' + message  +'</div></li>');
			
			// 保存临时ID，用于接收服务器返回的真实消息ID
			window.lastSentTempId = tempMsgId;
		}
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

// 右键菜单功能
$(document).ready(function(){
	// 创建右键菜单
	var contextMenuHtml = '<div id="context-menu" style="display:none;position:absolute;background:#fff;border:1px solid #ccc;border-radius:4px;box-shadow:0 2px 8px rgba(0,0,0,0.15);z-index:9999;padding:5px 0;">' +
		'<a href="#" id="recall-menu-item" style="display:block;padding:8px 16px;color:#333;text-decoration:none;font-size:14px;">撤回消息</a>' +
		'</div>';
	$('body').append(contextMenuHtml);
	
	// 右键点击自己的消息显示菜单
	$(document).on('contextmenu', '.chat_info li.right', function(e) {
		e.preventDefault();
		var msgId = $(this).attr('data-msg-id');
		if (!msgId) {
			layer.msg('消息正在发送中，请稍后再试');
			return;
		}
		
		// 检查是否在2分钟内
		var msgTime = $(this).find('i').text();
		var msgDate = new Date(msgTime);
		var now = new Date();
		var diffMinutes = (now - msgDate) / (1000 * 60);
		
		if (diffMinutes > 2) {
			layer.msg('消息超过2分钟，无法撤回');
			return;
		}
		
		$('#context-menu').css({
			top: e.pageY + 'px',
			left: e.pageX + 'px'
		}).show();
		
		// 保存当前消息ID
		$('#recall-menu-item').data('msg-id', msgId);
	});
	
	// 点击其他地方隐藏菜单
	$(document).click(function() {
		$('#context-menu').hide();
	});
	
	// 点击撤回
	$(document).on('click', '#recall-menu-item', function(e) {
		e.preventDefault();
		var msgId = $(this).data('msg-id');
		if (msgId) {
			recallMessage(msgId);
		}
		$('#context-menu').hide();
	});
});

// 撤回消息函数
function recallMessage(msgId) {
	let send_data = JSON.stringify({
		"status": 6, // 撤回消息类型
		"data": {
			"uid": $('.room').attr('data-uid').toString(),
			"username": $('.room').attr('data-username'),
			"avatar_id": $('.room').attr('data-avatar_id'),
			"room_id": $('.room').attr('data-room_id'),
			"msg_id": parseInt(msgId),
		}
	});
	ws.send(send_data);
}


