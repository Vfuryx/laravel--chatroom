import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store/index.js'
import Vant from 'vant';
import 'vant/lib/index.css';

import { Checkbox, CheckboxGroup } from 'vant';

Vue.use(Checkbox).use(CheckboxGroup);
Vue.use(Vant)

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
