<script setup lang="ts"></script>

<template>
  <a href="https://page.line.me/560mkygn">
    <button class="be-my-friend-btn">
      <img
        src="@/assets/img/line-round.png"
        width="30"
        class="line-logo-img"
        alt="LINE logo"
      >
      友だちに追加
    </button>
  </a>
</template>

<style lang="scss" scoped>
.be-my-friend-btn {
  width: 170px;
  height: 35px;
  cursor: pointer;

  border: none;
  border-radius: 17px;

  font-size: 1rem;
  font-weight: bold;
  line-height: 1.1rem;
  color: #ffffff;
  background: #06c755;

  position: relative;
  padding-left: 30px;

  .line-logo-img {
    position: absolute;
    top: 3px;
    left: 14px;
  }
}
</style>
