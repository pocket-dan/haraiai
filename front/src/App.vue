<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useRoute, RouteRecordName, RouteMeta, RouteLocationNormalizedLoaded } from 'vue-router'

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
const route = useRoute()
const titleSuffix = " - haraiai"
const siteUrl = "https://haraiai.netlify.app"

onMounted(async () => {
  setTags(route)
})

watch(route, () => {
  setTags(route)
})

const setTags = async (route: RouteLocationNormalizedLoaded) => {
  if (!route.name || !route.meta) return
  setTitleTags(route.name, route.meta.title)
  setDescTags(route.meta.desc)
  setOtherTags(route.path)
}

const setTitleTags = async (pageName: RouteRecordName, title: string) => {
  if (!title) return

  if (pageName !== "about") {
    title += titleSuffix
  }

  document.title = title
  document.querySelector("meta[name='title']")?.setAttribute("content", title)
  document.querySelector("meta[property='og:title']")?.setAttribute("content", title)
}

const setDescTags = async (desc: string) => {
  if (!desc) return

  document.querySelector("meta[name='description']")?.setAttribute("content", desc)
  document.querySelector("meta[property='og:description']")?.setAttribute("content", desc)
}

const setOtherTags = async (path: string) => {
  const canonicalUrl = siteUrl + path
  document.querySelector("link[rel='canonical']")?.setAttribute("href", canonicalUrl)
}
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
