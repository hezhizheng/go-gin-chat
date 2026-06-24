$(document).ready(function(){

    $(document).on('click', '#chat-list-li-top', function() {
    // $("#chat-list-li-top").click(function (){
        let offset = $("#hidden-chat-list-li-top").attr("data-offset")
        let room_id = $('.room').attr('data-room_id')
        let uid = getURLParam('uid')
        $.ajax({
            url: '/pagination?room_id='+room_id+'&offset='+offset+'&uid='+uid,
            success: function(data) {
                //设置数据
                var item = data.data.list
                if ( item == null )
                {
                    layer.msg('没有更多了！')
                    $("#hidden-chat-list-li-top").attr("data-offset",offset)
                    $("#chat-list-li-top").hide()
                    return false
                }


                $.each(item,function (index, value) {
                    if ( value.is_deleted ) {
                        // 已撤回的消息
                        $('#chat-list-li-top').after(
                            '<li class="systeminfo" data-msg-id="' + value.id + '"><span>【' + value.username + '】撤回了一条消息</span></li>');
                    } else if ( value.user_id == $("#body-room").attr("data-uid") ) {
                        // 自己发送的未撤回消息
                        $('#chat-list-li-top').after(
                            '<li class="right" data-msg-id="' + value.id + '"><img src="/static/images/user/' +
                            value.avatar_id +
                            '.png" alt=""><b>' +
                            value.username +
                            '</b><i>' +
                            value.created_at +
                            '</i><button class="btn撤回" style="font-size: 10px; margin-left: 10px; padding: 2px 5px;" data-msg-id="' + value.id + '">撤回</button><div class="aaa">' +
                            value.content+
                            '</div></li>');
                    } else {
                        // 他人发送的未撤回消息
                        $('#chat-list-li-top').after(
                            '<li class="left" data-msg-id="' + value.id + '"><img src="/static/images/user/' +
                            value.avatar_id +
                            '.png" alt=""><b>' +
                            value.username +
                            '</b><i>' +
                            value.created_at +
                            '</i><div class="aaa">' +
                            value.content+
                            '</div></li>');
                    }

                })
                $("#hidden-chat-list-li-top").attr("data-offset",parseInt(offset)+100)

            },
            error: function(data) {

            }
        });
    })
})


var getURLParam = function(name) {
    return decodeURIComponent((new RegExp('[?|&]' + name + '=' + '([^&;]+?)(&|#|;|$)', "ig").exec(location.search) || [, ""])[1].replace(/\+/g, '%20')) || '';
};