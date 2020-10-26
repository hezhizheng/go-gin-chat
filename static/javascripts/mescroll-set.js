var mescroll = new MeScroll("mescroll", { //第一个参数"mescroll"对应上面布局结构div的id (1.3.5版本支持传入dom对象)
    //如果您的下拉刷新是重置列表数据,那么down完全可以不用配置,具体用法参考第一个基础案例
    //解析: down.callback默认调用mescroll.resetUpScroll(),而resetUpScroll会将page.num=1,再触发up.callback
    down: {
        callback: function (mescroll ){
            let offset = $("#hidden-chat-list-li-top").attr("data-offset")
            let room_id = 1
            let uid = getURLParam('uid')
            console.log("uid",uid)
            $.ajax({
                url: '/pagination?room_id='+room_id+'&offset='+offset+'&uid='+uid,
                success: function(data) {
                    //联网成功的回调,隐藏下拉刷新的状态;
                    //无参. 注意结束下拉刷新是无参的
                    mescroll.endSuccess();

                    //设置数据
                    var item = data.data.list
                    if ( item == null )
                    {
                        layer.msg('没有更多了！')
                        // 销毁 mescroll
                        mescroll.lockDownScroll( true );
                        // 还原 offset
                        $("#hidden-chat-list-li-top").attr("data-offset",offset)

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
                    //联网失败的回调,隐藏下拉刷新的状态
                    mescroll.endErr();
                }
            });

        }, //下拉刷新的回调,别写成downCallback(),多了括号就自动执行方法了
        auto : false
    },
});

var getURLParam = function(name) {
    return decodeURIComponent((new RegExp('[?|&]' + name + '=' + '([^&;]+?)(&|#|;|$)', "ig").exec(location.search) || [, ""])[1].replace(/\+/g, '%20')) || '';
};