@extends('layout')

@section('title','聊天室')

@section('css')
    <style>
        * {
            padding: 0;
            margin: 0;
            list-style-type: none;
            font-size: 62.5%;
        }

        * focus {
            outline: none;
        }

        a {
            text-decoration: none;
        }

        html,
        body {
            width: 100%;
            height: 100%;
            overflow: hidden;
        }

        .clearfix:before,
        .clearfix:after {
            content: "";
            display: block;
            clear: both;
        }

        .clearfix {

            zoom: 1;

        }

        #login,
        #chat {
            height: 100%;
            width: 100%;
            position: absolute;
            z-index: 0;
            background: rgb(246, 247, 249);
            overflow: hidden;
            word-break: break-all;
        }

        #login .nick_input,
        #login .login_input {
            position: absolute;
            top: 50%;
            left: 50%;

            padding: 0.5rem 0;
            width: 80%;

            font-size: 2rem;
            line-height: 1.8;

            font-family: "Microsoft YaHei UI";
            text-align: center;

            border: none;
            outline: none;
            border-radius: 1rem;

            transition: box-shadow ease-out;
            transform: translate(-50%, -50%);
        }

        #login .nick_input:focus {
            transition: box-shadow ease-out;
            box-shadow: 0 0 10px #3a9692;
        }

        #login .nick_input {
            margin-top: -7rem;
            color: #317571;
            background: #fff;
        }

        #login .login_input {
            color: #ffffff;
            background: #6b7396;
        }

        #chat .chat_l,
        #chat .chat_r {
            width: 100%;
        }

        #chat .chat_l {
            position: absolute;
            height: 14rem;
            background: #ffffff;
        }

        #chat .chat_r {
            box-sizing: border-box;
            padding-top: 14rem;
            height: 100%;
        }

        #chat .chat_l .user_num {
            height: 15%
        }

        #chat .chat_l .user_lsit {
            /* width: 100%; */
            height: 55%;
            overflow-x: auto;
            /* 内容会被裁剪，会以滚动条显示 */
            overflow-y: hidden;
            /* 超出内容不可见 */
            white-space: nowrap;
            /* 不换行 */
        }

        #chat .chat_l .user_info {
            display: inline-block;
            box-sizing: border-box;
            /* height: 8rem; */
            padding: 0.5rem 0.8rem 0.5rem 1.5rem;
        }

        #chat .chat_l .user_info .avatar {
            /* margin: 1rem 0 0 1.5rem; */
            width: 4.5rem;
            height: 4.5rem;

            border-radius: 50%;
        }

        #chat .chat_l .user_info .chat_l_nick {
            display: block;
            margin: 0.2rem 0;

            font-size: 1rem;
            text-align: center;

            color: lightslategrey;
        }

        .btn_out {

            display: block;
            margin: 0 1rem;
            padding: 0.5rem 0;

            font-size: 1.5rem;
            color: white;
            text-align: center;

            background: rgb(250, 102, 102);

            cursor: pointer;
            outline: none;
            border: none;
            border-radius: 10px;
        }

        #chat .chat_r .top {
            height: 70%;
            background: #93a3bf;

            background: transparent;
            overflow-y: auto;
        }

        #chat .chat_r .top .line,
        #chat .chat_r .top .line_r {
            margin: 0.5rem 0;
        }

        #chat .chat_r .top .line .avatar {
            float: left;
            margin: 10px 0 0 15px;
            width: 4.5rem;
            height: 4.5rem;

            border-radius: 50%;

            background: transparent;
        }

        #chat .chat_r .top .line aside {
            padding-left: 7rem;
            padding-top: 0.8rem;
            padding-right: 2.5rem;

        }

        #chat .chat_r .top .line .chat_l_nick {
            font-size: 1.4rem;
        }

        #chat .chat_r .top .line .message {
            float: left;
            margin-top: 0.4rem;
            padding: 1rem;

            font-size: 1.4rem;
            color: #abc3d3;
            background: #e6edf4;

            border-radius: 1.5rem;
        }

        #chat .chat_r .top .line_r .avatar {
            float: right;
            margin: 1rem 1.5rem 0 0;
            width: 4.5rem;
            height: 4.5rem;

            border-radius: 50%;

            background: transparent;
        }

        #chat .chat_r .top .line_r aside {
            padding-right: 7.0rem;
            padding-top: 0.8rem;
            padding-left: 2.5rem;

        }

        #chat .chat_r .top .line_r .chat_l_nick {
            text-align: right;
            font-size: 1.4rem;
        }

        #chat .chat_r .top .line_r .message {
            float: right;
            margin-top: 0.4rem;
            padding: 1rem;

            font-size: 1.4rem;
            text-align: right;
            color: #abc3d3;
            background: #d5f4d5;

            border-radius: 1.5rem;
        }


        #chat .chat_r .bottom {
            position: relative;
            box-sizing: border-box;
            display: block;
            height: 30%;
            width: 100%;
            padding: 1rem;

            font-size: 2.4rem;
            color: white;

            border: none;
            border-radius: 1.5rem 1.5rem 0 0;

            outline: none;

            background: #9eaebf;
        }

        /* .btn_send {
            display: block;

            padding: 0.5rem 0;
            height: 6%;

            font-size: 1.5rem;
            color: white;
            text-align: center;

            background: rgb(106, 164, 252);

            cursor: pointer;
            outline: none;
            border: none;

        } */
    </style>
@endsection

@section('content')
    <div id="app">
        <div id="login" v-if='login'>
            <input class="nick_input" type="text" placeholder='昵称' v-model="nick">
            <button class="login_input" v-on:click="toLogin">登录</button>
        </div>
        <div id="chat" class="clearfix" v-else>
            <div class="chat_l">
                <p class="user_num">在线人数:@{{ num }}</p>
                <div class="user_lsit clearfix">

                    <p class="user_info clearfix">
                        <img class="avatar" src="http://qiniu.furyx.cn/photo.jpg">
                        <span class="chat_l_nick">@{{nick}}</span>
                        <input type="hidden" :value="members.fd">
                    </p>


                    <section class="user_info clearfix" v-for="member in members">
                        <img class="avatar" src="http://qiniu.furyx.cn/photo.jpg">
                        <p class="chat_l_nick">@{{ member.nick}}</p>
                        <input type="hidden" :value="member.fd">
                    </section>
                </div>

                <p class="btn_out" v-on:click="toOut">退出</p>
            </div>
            <div class="chat_r">
                <section class="top" id="data-list-content">
                    <article class=" line clearfix ">
                        <img class="avatar" src="http://qiniu.furyx.cn/photo.jpg">
                        <aside>
                            <p class="chat_l_nick">admin</p>
                            <p class="message">富强、民主、文明、和谐、自由、平等、公正、法治、爱国、敬业、诚信、友善</p>
                        </aside>
                    </article>
                    <article :class="[char.fd == fd ? 'line_r' : 'line', 'clearfix']" v-for=" char in chars">
                        <img class="avatar" src="http://qiniu.furyx.cn/photo.jpg">
                        <aside>
                            <p class="chat_l_nick">@{{char.nick}}</p>
                            <p class="message">@{{char.msg}}</p>
                        </aside>
                    </article>
                </section>
                <textarea class="bottom" @keyup.enter="sendMsg" v-model='message'></textarea>
                <!-- <p class="btn_send" v-on:click="sendMsg">发送</p> -->
            </div>
        </div>
    </div>



    <script>
        //    var i = 0;
        //    window.setInterval(function(){ //每隔5秒钟发送一次心跳，避免websocket连接因超时而自动断开
        //        console.log(i++);
        //    },1000);

        var app = new Vue({
            el: '#app',
            data: {
                //login: sessionStorage.getItem('login') == 'false' ? false : true,
                login: true,
                websock: null,
                fd: 0,
                chars: [],
                //nick: sessionStorage.getItem('nick') || '',
                nick: '',
                num: 1, //在线人数
                members: [],
                message: 'Hello',
                timeout: 58000, //58秒
                timeoutObj: null,
                lockReconnect: false, //避免重复连接
                wsurl: "ws://47.98.181.245/ws", //ws地址
            },
            methods: {
                toLogin() {
                    if (this.nick) {
                        this.login = false;
                        sessionStorage.setItem('login', false);
                        this.initWebSocket(this.wsurl)
                    }
                },
                toOut() {
                    this.login = true;
                    this.fd = 0;
                    this.member = [];
                    this.members = [];
                    this.message = '';
                    // sessionStorage.setItem('login', true);
                    // sessionStorage.setItem('nick', '');
                    this.websock.close();
                },
                addNick(nick) {
                    // sessionStorage.setItem('nick', this.nick);
                    $data = JSON.stringify({
                        'type': 'login',
                        'nick': this.nick
                    });
                    this.threadPoxi($data);
                },
                threadPoxi($data) { // 实际调用的方法
                    //参数

                    const agentData = this.message;
                    //若是ws开启状态
                    if (this.websock.readyState === this.websock.OPEN) {
                        console.log($data);
                        this.websocketsend($data);
                    }
                    // 若是 正在开启状态，则等待300毫秒
                    else if (this.websock.readyState === this.websock.CONNECTING) {
                        let that = this; //保存当前对象this
                        console.log(113);
                        setTimeout(function () {
                            that.websocketsend(that.data);
                        }, 300);

                    }
                    // 若未开启 ，则等待500毫秒
                    else {
                        console.log(114);

                        this.initWebSocket(this.wsurl);
                        let that = this; //保存当前对象this
                        setTimeout(function () {
                            that.websocketsend(that.data);
                        }, 500);
                    }
                },
                initWebSocket(url) { //初始化weosocket
                    this.websock = new WebSocket(url);
                    this.websock.onopen = this.websocketonopen;
                    this.websock.onmessage = this.websocketonmessage;
                    this.websock.onclose = this.websocketclose;
                    this.websock.onerror = this.websocketerror;
                },
                websocketonopen() {
                    this.addNick(this.nick);
                    this.reset(); //心跳检测重置
                },
                websocketonmessage(e) { //数据接收

                    let redata = JSON.parse(e.data);
                    console.log(redata);
                    switch (redata.stat) {
                        case 0: //err
                            console.log('err');
                            break;
                        case 1:
                            this.getMessage(redata);
                            break;
                        case 2: //open
                            this.getStart(redata);
                            break;
                        case 3: //ping
                            this.getStart(redata);
                            break;
                        case 4: //addmember
                            this.addMember(redata);
                            break;
                        case 5: //addmember
                            this.delMember(redata);
                            break;
                    }
                    this.reset(); //心跳检测重置
                },
                websocketsend(agentData) { //数据发送
                    this.websock.send(agentData);
                },
                websocketclose(e) { //关闭
                    console.log("connection closed (" + e.code + ")");
                    if (this.login === false) {
                        this.reconnect(this.wsurl); //重新连接
                    }
                },
                websocketerror() { //重新连接
                    this.reconnect(this.wsurl);
                },
                sendMsg() { //发送消息

                    let data = JSON.stringify({
                        'type': 'content',
                        'msg': this.message
                    });
                    this.threadPoxi(data);
                    this.message = '';
                },
                getMessage(data) {
                    console.log(data)
                    this.chars.push({
                        'fd': data.fd,
                        'nick': data.nick,
                        'time': data.time,
                        'msg': data.msg
                    });
                    // let add = data.nick + '-时间:' + data.time + ';' + '内容:' + data.msg + ";\n"
                    // this.char += add;
                },
                getStart(data) {
                    console.log(data);
                    this.fd = data.fd;
                    this.num = data.num;
                    for (member in data.members) {
                        this.members.push({
                            fd: member,
                            nick: data.members[member]
                        });
                    }

                    console.log(this.members);

                },
                addMember(data) {
                    console.log(this.members);
                    this.num = data.num;
                    this.members.push({
                        fd: data.fd,
                        nick: data.nick
                    });
                    console.log(this.members);
                },
                delMember(data) {
                    console.log(data);
                    this.num -= 1;
                    for (member in this.members) {
                        console.log(member);
                        if (this.members[member].fd == data.fd) {
                            this.members.splice(member, 1);
                            break;
                        }
                    }
                    console.log(this.members);
                },
                reset: function () { //心跳检测重置
                    clearTimeout(this.timeoutObj);
                    this.start();
                },
                start: function () {
                    var that = this;
                    this.timeoutObj = setTimeout(function () {
                        //这里发送一个心跳，后端收到后，返回一个心跳消息，
                        //onmessage拿到返回的心跳就说明连接正常
                        $data = {
                            "type": "ping"
                        };
                        that.websock.send(JSON.stringify($data));
                    }, this.timeout)
                },
                reconnect: function (url) { //重新连接
                    if (this.lockReconnect) return;
                    this.lockReconnect = true;
                    //没连接上会一直重连，设置延迟避免请求过多
                    setTimeout(function () {
                        this.initWebSocket(url);
                        this.lockReconnect = false;
                    }, 2000);
                },
                scrollToBottom: function () {
                    this.$nextTick(() => {
                        var div = document.getElementById('data-list-content')
                        div.scrollTop = div.scrollHeight
                })
                }

            },
            created() {
                // if (sessionStorage.getItem('login') == 'false') {
                //     this.initWebSocket(this.wsurl)
                // }
            },
            watch: {
                'chars': 'scrollToBottom' //监视char变化更改导航条的变化
            }
        })
    </script>
@endsection