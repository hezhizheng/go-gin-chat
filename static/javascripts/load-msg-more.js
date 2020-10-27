$(document).ready(function(){

    $("#chat-list-li-top").click(function (){
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
                    if ( value.user_id == $("#body-room").attr("data-uid") )
                    {
                        $('#chat-list-li-top').after(
                            '<li class="right"><img src="/static/images/user/' +
                            value.avatar_id +
                            '.png" alt=""><b>' +
                            value.username +
                            '</b><i>' +
                            value.created_at +
                            '</i><div class="aaa">' +
                            value.content+
                            '</div></li>');
                    }else{
                        $('#chat-list-li-top').after(
                            '<li class="left"><img src="/static/images/user/' +
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