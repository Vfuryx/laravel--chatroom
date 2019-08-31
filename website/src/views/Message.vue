<template>
  <div class="Message">
    <van-nav-bar
      left-text="返回"
      :title="this.$route.params.to_name"
      left-arrow
      @click-left="onClickLeft"
    ></van-nav-bar>
    <div style="padding-bottom: 300px;">
      <div v-for="record in records" style="margin-top: 15px;">
        <div
          :style="record.is_me ? 'padding:0 15px 0 0;float:right' : 'padding:0 0 0 15px;float:left'"
        >
          <van-icon
            size="40px"
            name="http://juntao.oss-cn-shenzhen.aliyuncs.com/images/wechat/avatar/xX9hSG0utRz0hmC4hA5b4nWL46Jgn8YGhQ7psWyU.png"
          />
        </div>
        <div class="msg-content">
          <div :class="record.is_me ? 'right-msg' : 'left-msg' " style="word-wrap:break-word">
            {{ record.type == 1 ? record.content : '' }}
            <img
              v-if="record.type == 2"
              :src="record.content"
            >
          </div>
        </div>
      </div>
    </div>

    <div class="send">
      <van-row>
        <van-col span="24" class="controller">
          <div class="one">
            <van-icon size="20px" name="volume" color="#737373" class="my-icon"/>
          </div>
          <div class="two">
            <van-field v-model="value" placeholder="请输入" @keyup.enter="send"/>
          </div>
          <div class="three">
            <van-icon
              size="20px"
              name="smile-o"
              color="#737373"
              class="my-icon"
              @click="showBottom('emoji')"
            />
            <van-icon
              size="20px"
              name="add"
              color="#737373"
              class="my-icon"
              @click="showBottom('pic')"
            />
          </div>
        </van-col>
      </van-row>
      <van-row v-if="active == 'emoji'">
        <van-col class="smile-content" span="24" style="background-color:#f8f8f8;">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/1.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/198.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/201.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/202.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/429.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/430.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/431.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/432.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/434.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/438.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/460.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/701.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/705.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/703.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/706.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/719.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/720.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/725.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/728.png" class="smile-icon" alt @click="sendEmoji">
          <img src="https://raw.githubusercontent.com/tmm1/emoji-extractor/master/images/160x160/741.png" class="smile-icon" alt @click="sendEmoji">
        </van-col>
      </van-row>
      <van-row v-if="active == 'pic'">
        <van-col span="24" style="background-color:#f8f8f8;">
          <van-uploader :after-read="onRead" style="margin: 20px 15px;">
            <van-icon size="45px" name="photo" color="#737373" class="send-pic"/>
            <div style="text-align: center;font-size:14px">上传图片</div>
          </van-uploader>
        </van-col>
      </van-row>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'Message',
  data() {
    return {
      active: "",
      emoji: "http://juntao.oss-cn-shenzhen.aliyuncs.com/images/video/courses/packages/f2HQCKV4GMcG4Ajfu0z4eXHEJQ5vOdW6S3Mcc4oo.jpeg",
      value: "",
      records: {}
    }
  },

  components: {
  },
  created() {
    this.records = this.$store.getters.getRecordBy(this.$route.params.to_name)
  },
  methods: {
    onClickLeft() {
      this.$router.push({ name: 'index', params: {} })
    },
    showBottom(active) {
      if (this.active == active) {
        this.active = ''
        return
      }
      this.active = active
    },
    send() {
      let data = {
        Source: 1,
        From: this.user,
        To: this.$route.params.to_name,
        Type: 1,
        Content: this.value,
        SendAt: "2015"
      }
      let record = {
        from: this.$route.params.to_name,
        type: 1,
        is_me: true,
        content: this.value
      }
      this.$store.dispatch('giveMe', record)
      this.ws.send(JSON.stringify(data))
      this.value = ''
    },
    onRead(file) {
      console.log(file)
    },
    sendEmoji(event) {

      console.log(event.currentTarget.src)

      let data = {
        Source: 1,
        From: this.user,
        To: this.$route.params.to_name,
        Type: 2,
        Content: event.currentTarget.src,
        SendAt: "2015"
      }
      let record = {
        from: this.$route.params.to_name,
        type: 2,
        is_me: true,
        content: event.currentTarget.src
      }

      this.$store.dispatch('giveMe', record)

      this.ws.send(JSON.stringify(data))
    }
  },
  computed: {
    ...mapGetters([
      'ws',
      'user'
    ])
  },
}
</script>

<style>
.send {
  width: 100%;
  left: 0;
  bottom: 0;
  position: fixed;
  background-color: #fff;
}

.send .controller {
  background-color: #fff;
  display: flex;
}

.send .controller .one {
  flex: 0 0 50px;
}

.send .controller .two {
  flex: 1;
}

.send .controller .three {
  flex: 0 0 100px;
}

.my-icon {
  margin: 8px;
  padding: 6px;
  border-radius: 50%;
  border: 1px solid #737373;
}

@media screen and (max-width: 1920px) {
  .smile-icon {
    box-sizing: border-box;
    width: 5%;

    padding: 23px;
    border-radius: 50%;
    border: 1px solid #f8f8f8;
  }
}

@media screen and (max-width: 1500px) {
  .smile-icon {
    box-sizing: border-box;
    width: 6.66%;

    padding: 19px;
    border-radius: 50%;
    border: 1px solid #f8f8f8;
  }
}

@media screen and (max-width: 1400px) {
  .smile-icon {
    box-sizing: border-box;
    width: 8.33%;

    padding: 16px;
    border-radius: 50%;
    border: 1px solid #f8f8f8;
  }
}

@media screen and (max-width: 1180px) {
  .smile-icon {
    box-sizing: border-box;
    width: 10%;

    padding: 16px;
    border-radius: 50%;
    border: 1px solid #f8f8f8;
  }
}

@media screen and (max-width: 1000px) {
  .smile-icon {
    box-sizing: border-box;
    width: 11.11%;

    padding: 13px;
    border-radius: 50%;
    border: 1px solid #f8f8f8;
  }
}

@media screen and (max-width: 750px) {
  .smile-icon {
    box-sizing: border-box;
    width: 12.5%;

    padding: 13px;
    border-radius: 50%;
    border: 1px solid #f8f8f8;
  }
}

@media screen and (max-width: 500px) {
  .smile-icon {
    box-sizing: border-box;
    width: 16.66%;

    padding: 10px;
    border-radius: 50%;
    border: 1px solid #f8f8f8;
  }
}

@media screen and (max-width: 300px) {
  .smile-icon {
    box-sizing: border-box;
    width: 20%;

    padding: 7px;
    border-radius: 50%;
    border: 1px solid #f8f8f8;
  }
}

.send-pic {
  
  padding: 8px;
  margin: 15px 15px 5px;

  border: 1px solid #fff;
  border-radius: 10%;

  background: #fff;
}

.msg-content {
  max-width: 100%;
  box-sizing: border-box;

  padding: 0 70px;
  color: rgb(0, 0, 0);
}

.msg-content::after {
  content: "";
  display: block;
  clear: both;
}

.left-msg {
  position: relative;
  float: left;
  padding: 10px;
  max-width: 100%;
  box-sizing: border-box;

  color: rgb(0, 0, 0);
  background: rgb(255, 255, 255);
  border-radius: 5px;
}

.left-msg:after {
  content: "";

  display: block;

  position: absolute;
  top: 10px;
  left: 0;

  margin-left: -19px;

  border: 10px solid transparent;
  border-right-color: rgb(255, 255, 255);
}

.right-msg {
  position: relative;
  float: right;
  padding: 10px;
  max-width: 100%;
  box-sizing: border-box;

  text-align: left;
  color: rgb(0, 0, 0);
  background: rgb(88, 228, 151);

  border-radius: 5px;
}

.right-msg:after {
  content: "";

  display: block;

  position: absolute;
  top: 10px;
  right: 0;

  margin-right: -19px;

  border: 10px solid transparent;
  border-left-color: rgb(88, 228, 151);
}
</style>
