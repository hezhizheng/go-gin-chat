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
                    var isWithdrawn = value.is_withdrawn == 1;
                    var contentHtml = isWithdrawn 
                        ? '<span style="color: #999; font-style: italic;">消息已撤回</span>' 
                        : value.content;
                    var withdrawnClass = isWithdrawn ? 'withdrawn' : '';
                    var msgIdAttr = 'data-msg-id="' + value.id + '"';
                    var uidAttr = 'data-uid="' + value.user_id + '"';
                    
                    if ( value.user_id == $("#body-room").attr("data-uid") )
                    {
                        $('#chat-list-li-top').after(
                            '<li class="right" ' + msgIdAttr + ' ' + uidAttr + '><img src="/static/images/user/' +
                            value.avatar_id +
                            '.png" alt=""><b>' +
                            value.username +
                            '</b><i>' +
                            value.created_at +
                            '</i><div class="aaa ' + withdrawnClass + '">' +
                            contentHtml +
                            '</div></li>');
                    }else{
                        $('#chat-list-li-top').after(
                            '<li class="left" ' + msgIdAttr + ' ' + uidAttr + '><img src="/static/images/user/' +
                            value.avatar_id +
                            '.png" alt=""><b>' +
                            value.username +
                            '</b><i>' +
                            value.created_at +
                            '</i><div class="aaa ' + withdrawnClass + '">' +
                            contentHtml +
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