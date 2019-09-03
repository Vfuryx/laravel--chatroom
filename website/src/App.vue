<template>
  <div id="app" style="bag">
    <router-view/>
    
  </div>
</template>

<script>
export default {
  name: 'app',
  data() {
    return {
     
    }
  }
}
</script>

<style>
body {
  background-color: whitesmoke;
}
#app {
  font-family: Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}
</style>
