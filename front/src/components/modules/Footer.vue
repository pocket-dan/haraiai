<script setup lang="ts">
</script>

<template>
  <footer>
    <router-link :to="{'name': 'about'}">
      サービスについて
    </router-link>
    <router-link :to="{'name': 'help'}">
      よくある質問
    </router-link>
    <router-link :to="{'name': 'inquiry'}">
      フィードバック
    </router-link>
    <br>

    <router-link :to="{'name': 'tos'}">
      利用規約
    </router-link>
    <router-link :to="{'name': 'privacy'}">
      プライバシーポリシー
    </router-link>

    <br>

    <p>
      Developed by <a
        href="https://github.com/raahii"
        target="_blank"
      >raahii</a>.
    </p>
    <p class="credit">
      Special thanks to <br>
      <a
        href="https://www.flaticon.com/free-icons/donation"
        title="donation icons"
      >Donation icons created by Freepik - Flaticon</a>
    </p>
  </footer>
</template>

<style lang="scss" scoped>
footer {
  margin: 50px 0 30px;
  display: flex;
  align-items: center;
  flex-direction: column;

  .credit {
    font-size: 0.8rem;
    color: gray;
  }
}
</style>
