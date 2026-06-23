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
                    let isRecalled = value.is_recalled == 1;
                    let contentClass = isRecalled ? 'msg-content recalled' : 'msg-content';
                    let content = isRecalled ? '消息已撤回' : value.content;

                    if ( value.user_id == $("#body-room").attr("data-uid") )
                    {
                        let recallBtn = '';
                        // 只有2分钟内的消息才显示撤回按钮
                        if (!isRecalled && isWithinRecallTimeForHistory(value.created_at)) {
                            recallBtn = '<span class="recall-btn" onclick="recallMessage(' + value.id + ', \'' + value.created_at + '\')">撤回</span>';
                        }
                        $('#chat-list-li-top').after(
                            '<li class="right" data-msg-id="' + value.id + '" data-created-at="' + value.created_at + '"><img src="/static/images/user/' +
                            value.avatar_id +
                            '.png" alt=""><b>' +
                            value.username +
                            '</b><i>' +
                            value.created_at +
                            '</i><div class="' + contentClass + '">' +
                            content+
                            '</div>' + recallBtn + '</li>');
                    }else{
                        $('#chat-list-li-top').after(
                            '<li class="left" data-msg-id="' + value.id + '"><img src="/static/images/user/' +
                            value.avatar_id +
                            '.png" alt=""><b>' +
                            value.username +
                            '</b><i>' +
                            value.created_at +
                            '</i><div class="' + contentClass + '">' +
                            content+
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

// 判断历史消息是否在2分钟内（用于加载更多消息时）
function isWithinRecallTimeForHistory(createdAt) {
    let msgTime = new Date(createdAt).getTime();
    let now = new Date().getTime();
    return (now - msgTime) < 2 * 60 * 1000; // 2分钟
}