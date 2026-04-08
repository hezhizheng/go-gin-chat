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
                    let contentHtml = '';
                    let revokeBtnHtml = '';
                    
                    if (value.is_revoked) {
                        contentHtml = '<div style="color:#999;font-style:italic;">[消息已撤回]</div>';
                    } else {
                        contentHtml = '<div class="aaa">' + value.content + '</div>';
                        if (value.user_id == $("#body-room").attr("data-uid")) {
                            revokeBtnHtml = '<button class="revoke-btn" style="font-size:12px;cursor:pointer;">撤回</button>';
                        }
                    }
                    
                    if ( value.user_id == $("#body-room").attr("data-uid") )
                    {
                        $('#chat-list-li-top').after(
                            '<li class="right" data-msg-id="' + value.id + '" data-created-at="' + value.created_at + '"><img src="/static/images/user/' +
                            value.avatar_id +
                            '.png" alt=""><b>' +
                            value.username +
                            '</b><i>' +
                            value.created_at +
                            '</i>' + revokeBtnHtml + contentHtml + '</li>');
                    }else{
                        $('#chat-list-li-top').after(
                            '<li class="left" data-msg-id="' + value.id + '"><img src="/static/images/user/' +
                            value.avatar_id +
                            '.png" alt=""><b>' +
                            value.username +
                            '</b><i>' +
                            value.created_at +
                            '</i>' + contentHtml + '</li>');
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