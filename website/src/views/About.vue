<template>
  <div class="about">
    <van-tabbar route v-model="active" style="height:55px;">
      <van-tabbar-item replace to="/" info="10" icon="chat">消息</van-tabbar-item>
      <van-tabbar-item replace to="/about" icon="manager">个人</van-tabbar-item>
    </van-tabbar>
    <van-cell-group>
      <van-cell>
        <template slot="icon">
          <van-icon
            size="60px"
            name="http://juntao.oss-cn-shenzhen.aliyuncs.com/images/wechat/avatar/xX9hSG0utRz0hmC4hA5b4nWL46Jgn8YGhQ7psWyU.png"
            class="custom-icon"
          />
        </template>
        <template slot="title">
          <div style="padding-left:10px;padding-top:10px;">
            <span class="custom-text">张三</span>
            <p style="margin:5px 0;color:rgb(102, 102, 102);">121221111212</p>
          </div>
        </template>
      </van-cell>
    </van-cell-group>
    <van-cell-group style="margin:5px 0 0;">
      <van-cell title="????" icon="location-o" is-link/>
    </van-cell-group>
  </div>
</template>

<script>
export default {
  name: 'about',
  data() {
    return {
      active: 1
    }
  },
  components: {
  },
  methods: {
    onClickLeft() {
      // Toast('返回');
    },
  }
}
</script>

<style>
.custom-icon {
  margin: 10px;
}
</style>
