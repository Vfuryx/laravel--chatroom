<template>
  <div class="loign">
    <van-cell-group>
      <van-field
        v-model="username"
        required
        clearable
        label="用户名"
        right-icon="question-o"
        placeholder="请输入用户名"
      />

      <van-field
        v-model="password"
        required
        type="password"
        abel="密码"
        placeholder="请输入密码"
        label="密码"
      />
    </van-cell-group>
    <van-button type="default" size="large" class="login-btn" @click="login">提交</van-button>

  </div>

</template>

<script>
import { mapGetters } from 'vuex'
import { Dialog } from 'vant';

export default {
  name: 'login',
  data() {
    return {
      username: "",
      password: "",
    }
  },
  methods: {
    login() {

      if (this.username == '' || this.password == '') {
        Dialog.alert({
          title: '账号或密码错误',
          message: '请输入正确的账号和密码'
        }).then(() => {
          // on close
        });
        return false
      }
      this.$store.dispatch('InitWebSocket', {
        username: this.username,
        password: this.password
      })

      this.ws.onopen = this.websocketOnOpen
      this.ws.onmessage = this.websocketOnMessage
      this.ws.onclose = this.websocketClose
      this.ws.onerror = this.websocketError
      if (this.$store.getters.isLogin == true) {
        this.$router.push({ name: 'index', params: {} })
      }
    },
    websocketOnOpen() {
      let data = {
        Source: 4,
        From: this.user,
        To: '',
        Type: 0,
        SendAt: ''
      }
      // 获取用户列表
      this.$store.getters.ws.send(JSON.stringify(data))

      data = {
        Source: 9,
        From: this.user,
        To: '',
        Type: 0,
        SendAt: ''
      }
      // 获取群组列表
      this.$store.getters.ws.send(JSON.stringify(data))

      this.reset()
    },
    reset() { //心跳检测重置
      clearTimeout(this.$store.getters.timeoutObj)

      this.start()
    },
    start() {
      var that = this;
      let obj = setTimeout(function () {
        //这里发送一个心跳，后端收到后，返回一个心跳消息，
        //onmessage拿到返回的心跳就说明连接正常
        let data = {
          "type": "ping"
        };
        that.send(JSON.stringify(data));
      }, this.timeout)

      this.$store.dispatch('SetTimeout', obj)
    },
    websocketOnMessage(e) {
      let redata = JSON.parse(e.data);
      console.log(e.data);
      this.msgHandle(redata)
      // switch (redata.stat) {
      //     case 0: //err
      //         console.log('err');
      //         break;

      // }
      this.reset(); //心跳检测重置
    },
    websocketClose(e) { //关闭
      console.log("connection closed (" + e.code + ")");

      if (this.$store.state.websocket.login === false) {
        this.reconnect(this.$store.state.websocket.url); //重新连接
      }
    },
    send(data) {
      this.ws.send(data)
    },
    websocketError(e) {
      this.reconnect(this.$store.state.websocket.url);
    },
    reconnect() { //重新连接
      if (this.$store.state.websocket.lockReconnect) return;
      this.$store.state.websocket.lockReconnect = true;

      //没连接上会一直重连，设置延迟避免请求过多
      var that = this
      setTimeout(() => {
        this.initWebSocket(that.user, that.state.websocket.password);
        this.$store.state.websocket.lockReconnect = false;
      }, 2000);
    },
    msgHandle(msg) {
      switch (msg.Source) {
        case 1:  //私信
          this.giveMe(msg)
          break;
        case 2:  //群聊
          this.group(msg)
          break;
        case 3:  //ping pong
          break;
        case 4:  //用户列表
          this.getUserList(msg)
          break;
        case 5:  //添加用户
          this.addUser(msg)
          break;
        case 6:  //移除用户
          this.removeUser(msg)
          break;
        case 7:  //公共群聊
          this.createPublicGroup(msg)
          break;
        case 9:  //获取群组
          this.getGroupList(msg)
          break;

      }
    },
    getUserList(msg) {
      let userList = JSON.parse(msg.Content)
      this.$store.dispatch('setList', userList)
    },
    getGroupList(msg) {
      let groupList = JSON.parse(msg.Content)
      this.$store.dispatch('setGroupList', groupList)
    },
    createPublicGroup(msg) {
      console.log(msg)
      this.$store.dispatch('creatGroup', msg)
    },
    giveMe(msg) {
      let data = {
        from: msg.From,
        type: msg.Type,
        content: msg.Content
      }
      console.log(data)
      this.$store.dispatch('giveMe', data)
    },
    //群聊
    group(msg) {
      let data = {
        from: msg.From,
        to: msg.To,
        type: msg.Type,
        content: msg.Content
      }
      console.log(data)
      this.$store.dispatch('group', data)
    },
    addUser(msg) {
      let data = {
        name: msg.Content,
      }
      this.$store.dispatch('addUser', data)
      console.log(data)
      this.$store.dispatch('addRecord', msg.Content)
    },
    removeUser(msg) {
      this.$store.dispatch('removeUser', msg.Content)
    }
  },
  computed: {
    ...mapGetters([
      'ws',
      'user'
    ])
  },
  components: {
  }
}
</script>

<style>
.login-btn {
  display: block;
}
</style>

