<template>
  <div class="index">
    <van-nav-bar :title="'在线好友('+ (this.chatUserList.length - 1) +')'">
      <div slot="right" style="height:46px;line-height:46px">
        <i style="color:rgb(30, 228, 32);">•</i> 在线
      </div>
    </van-nav-bar>
    <van-cell-group>
      <van-cell v-for="user of this.chatUserList" value="12:00" v-if="user.name != name"  :to="{name: 'message', params: { to_name: user.name }}">
        <template slot="icon">
          <van-icon
            size="60px"
            name="http://juntao.oss-cn-shenzhen.aliyuncs.com/images/wechat/avatar/xX9hSG0utRz0hmC4hA5b4nWL46Jgn8YGhQ7psWyU.png"
            class="custom-icon"
          />
        </template>
        <template slot="title">
          <div style="padding-left:10px">
            <span class="custom-text">{{user.name}}</span>
            <p style="margin:5px 0;color:rgb(152, 152, 152);">{{ name }}</p>
          </div>
        </template>
      </van-cell>
    </van-cell-group>
    <van-button type="default" @click="debug">默认按钮</van-button>
    <van-tabbar route v-model="active" style="height:55px;">
      <van-tabbar-item replace to="/" info="10" icon="chat">消息</van-tabbar-item>
      <van-tabbar-item replace to="/about" icon="manager">个人</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'index',
  data() {
    return {
      active: 0,
      list: [],
      name: ''
    }
  },
  created(){
    this.name = this.$store.getters.user
  },
  methods: {
    debug() {
      console.log(this.$store.getters.getRecord)
      // for (let i of this.chatUserList) {
      //   console.log(i)
      // }
    }
  },
  computed: {
    ...mapGetters([
      'chatUserList'
    ])
  },
  components: {
  }
}
</script>

<style>
</style>

