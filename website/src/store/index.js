import Vue from 'vue'
import Vuex from 'vuex'
import websocket from './modules/websocket'
import chat from './modules/chat'

Vue.use(Vuex)

const store = new Vuex.Store({
  modules: {
    websocket,
    chat,
  }
});

export default store