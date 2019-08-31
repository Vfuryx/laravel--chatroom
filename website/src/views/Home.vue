<template>
  <div class="home">
    <van-button>122</van-button>
    <van-button>565</van-button>
  </div>
</template>

<script>

import { Button } from "vant2"

export default {
  name: 'home',
  components: {
    [Button.name]: Button
  }
}
</script>
