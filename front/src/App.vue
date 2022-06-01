<script setup lang="ts">
// import { onMounted } from 'vue'
import Footer from "@/components/modules/Footer.vue"

// import liff from "@line/liff";

// onMounted(async () => {
//   liff
//     .init({
//       liffId: "1657131873-ZNGno007",
//     })
//     .then(() => {
//       console.debug("LIFF init succeeded.")
//       if (!liff.isLoggedIn()) {
//         liff.login()
//       }
//     })
//     .catch((e: Error) => {
//       console.error("LIFF init failed.", e)
//     });
// })
</script>

<template>
  <router-view></router-view>
  <Footer />
</template>

<style>
html, body {
  margin: 0;
  background-color: #F8FAFB;
}

#app {
  font-family: 'Zen Maru Gothic', Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;

  width: 100%;
}
</style>
