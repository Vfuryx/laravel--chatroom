import Vue from 'vue'
import App from './App.vue'
import router from './router'
import vant from 'vant2';
import store from './store/index.js'
import 'vant2/lib/index.css';

Vue.use(vant);

Vue.config.productionTip = false

router.beforeEach((to, from, next) => {
    // 如果登录了
    if(store.getters.isLogin){

      if(to.name == 'login'){ next({path: '/'})}

      next()

    }else{

      if(to.name == 'login'){
        next();
      }else{
        next({path: '/login'})
      }

    }
});

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
