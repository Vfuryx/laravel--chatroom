<template>
  <div class="index">
    <van-nav-bar :title="'在线好友('+ (this.chatUserList.length - 1) +')'">
      <div slot="left">
        <van-icon size="20px" name="add-o" color="#737373" @click="add = !add" />
        <div class="index-add-wrap" v-if="add" @click="add = !add"></div>
        <div class="index-add-mask" v-if="add">
          <div class="add-btn" @click="showPublic = true; add = false">
            <van-icon size="20px" name="chat" color="#fff" />发起公开群聊
          </div>
          <div class="add-btn" @click="showPrivate = true; add = false">
            <van-icon size="20px" name="chat" color="#fff" />发起私密群聊
          </div>
        </div>
        <van-dialog
          v-model="showPublic"
          title="发起公开群聊"
          show-cancel-button
          @confirm="createPublicGroup"
        >
          <van-cell-group>
            <van-field
              value
              v-model="createGroupName"
              label="群组名"
              left-icon="contact"
              placeholder="输入群组名称"
            />
          </van-cell-group>
        </van-dialog>
        <van-dialog
          v-model="showPrivate"
          title="发起私密群聊"
          show-cancel-button
          @confirm="createPrivateGroup"
        >
          <van-cell-group>
            <van-field
              value
              v-model="createGroupName"
              label="群组名"
              left-icon="contact"
              placeholder="输入群组名称"
            />
            <van-field value label="用户列表" left-icon="contact">
              <div slot="input">
                <van-checkbox-group v-model="addResult">
                  <van-cell-group v-for="user of this.chatUserList">
                    <van-cell v-if="user.name != name" clickable :key="user.name">
                      <div slot="title">
                        <van-icon
                          size="28px"
                          name="http://juntao.oss-cn-shenzhen.aliyuncs.com/images/wechat/avatar/xX9hSG0utRz0hmC4hA5b4nWL46Jgn8YGhQ7psWyU.png"
                          class="custom-icon"
                        />
                        <span class="padding-left-8px">{{ user.name }}</span>
                      </div>
                      <van-checkbox :name="user.name" ref="checkboxes" slot="right-icon" />
                    </van-cell>
                  </van-cell-group>
                </van-checkbox-group>
              </div>
            </van-field>
          </van-cell-group>
        </van-dialog>
      </div>
      <div slot="right" style="height:46px;line-height:46px">
        <i style="color:rgb(30, 228, 32);">•</i> 在线
      </div>
    </van-nav-bar>
    <van-cell-group>
      <van-cell v-for="group in this.groupList" value="12:00" :to="{name: 'message', params: { to_name: group, group: true }}">
        <template slot="icon">
          <van-icon
            size="60px"
            name="http://juntao.oss-cn-shenzhen.aliyuncs.com/images/wechat/avatar/xX9hSG0utRz0hmC4hA5b4nWL46Jgn8YGhQ7psWyU.png"
            class="custom-icon"
          />
        </template>
        <template slot="title">
          <div style="padding-left:10px">
            <span class="custom-text">{{ group }} (群聊)</span>
            <!-- <p style="margin:5px 0;color:rgb(152, 152, 152);">{{  }}</p> -->
          </div>
        </template>
      </van-cell>
    </van-cell-group>
    <van-cell-group>
      <van-cell
        v-for="user of this.chatUserList"
        value="12:00"
        v-if="user.name != name"
        :to="{name: 'message', params: { to_name: user.name, group: false }}"
      >
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
            <p style="margin:5px 0;color:rgb(152, 152, 152);">{{ user.name }}</p>
          </div>
        </template>
      </van-cell>
    </van-cell-group>
    <van-tabbar route v-model="active" style="height:55px;">
      <van-tabbar-item replace to="/" info="10" icon="chat">消息</van-tabbar-item>
      <van-tabbar-item replace to="/about" icon="manager">个人</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import axios from 'axios'
import { Dialog } from 'vant';

export default {
  name: 'index',
  data() {
    return {
      active: 0,
      list: [],
      name: '',
      add: false,
      addResult: [
      ],
      showPublic: false,
      showPrivate: false,
      createGroupName: '',
    }
  },
  created() {
    this.name = this.$store.getters.user
  },
  methods: {
    debug() {
      axios.get("http://web.zhbx/api/video/play/auth/fcec7d5010e34af892faa7108dc7c51e")
        .then((res) => {
          console.log(res)
        })
        .catch((err) => {
          console.log(err)
        })
      // console.log(this.$store.getters.getRecord)
      // for (let i of this.chatUserList) {
      //   console.log(i)
      // }
    },
    createPublicGroup() {
      if (this.createGroupName == '') {
        Dialog.alert({
          title: '请输入群组名称',
          message: '请输入正确的群组名称'
        }).then(() => {
          // on close
        });
        return false;
      }
      let data = {
        Source: 7,
        From: this.name,
        To: '',
        Type: 1,
        Content: this.createGroupName,
        SendAt: "2015"
      }
      this.ws.send(JSON.stringify(data))

    },
    createPrivateGroup() {
      console.log(this.createGroupName)
      console.log(this.addResult)
    }
  },
  computed: {
    ...mapGetters([
      'ws',
      'chatUserList',
      'groupList'
    ])
  },
  components: {
  }
}
</script>

<style>
.index-add-mask {
  position: absolute;
  margin-top: 5px;
  margin-left: -10px;
  padding: 2px 0;
  width: 128px;
  z-index: 3;

  background: #4c4c4c;

  border-radius: 5px;
}

.index-add-wrap {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1;

  background: transparent;
}

.index-add-mask::before {
  content: "";
  position: absolute;
  top: -18px;
  left: 10px;

  z-index: 999;

  border: 10px transparent solid;
  border-bottom-color: #4c4c4c;
}

.index-add-mask .add-btn {
  padding-left: 10px;
  margin: -5px 0;

  text-align: left;
  color: aliceblue;
}

.index-add-mask .add-btn i {
  position: relative;
  margin-top: -2px;
  margin-right: 5px;
}

.padding-left-8px {
  padding-left: 8px;
}
</style>

