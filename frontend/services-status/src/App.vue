<template>
  <div id="app">
    <ServicesStatus />
  </div>
</template>

<script>
import ServicesStatus from './components/ServicesStatus.vue';

export default {
  name: 'App',
  components: {
    ServicesStatus,
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
