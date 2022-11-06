<script setup lang="ts">
import BeMyFriendButton from '@/components/modules/BeMyFriendButton.vue';
</script>

<template>
  <div class="about-page">
    <section class="gray">
      <img
        class="logo"
        src="@/assets/img/logo.png"
        width="100"
        alt="Logo"
      >
      <h1>haraiai - 払い合い</h1>
      <p class="desc">
        LINE上ですぐに使える <br />
        カップルのための折半アプリ
      </p>
      <BeMyFriendButton />
    </section>

    <section
      id="feature1"
      class="white feature"
    >
      <h2>LINE上で手軽に使える</h2>
      <p class="desc">
        はじめ方は簡単。グループを作って<br />相手 と haraiai を追加するだけ。<br>
        アプリのインストールや会員登録なしで すぐに始められます。
      </p>

      <img
        src="@/assets/img/f1.png"
        alt="Feature1"
      >
    </section>

    <section
      id="feature2"
      class="gray feature"
    >
      <h2>支払い金額を自動で折半</h2>
      <p class="desc">
        あとはグループに支払いを記録するだけ。
        haraiai が自動で折半し、支払い総額も教えてくれます。
      </p>

      <img
        class="feature"
        src="@/assets/img/f2.png"
        alt="Feature2"
      >
    </section>

    <section
      id="feature3"
      class="white feature"
    >
      <h2>退屈な精算はもうしない</h2>
      <p class="desc">
        面倒な定期精算をなくして払い合いませんか？
        そんな気持ちが名前に込められています。
        支払いが少ない側が次は払うことで 総支払い額が釣り合うようにしましょう！
      </p>

      <img
        src="@/assets/img/f3.png"
        alt="Feature3"
      >
    </section>

    <section
      id="add-friend"
      class="gray"
    >
      <BeMyFriendButton />
    </section>
  </div>
</template>

<style lang="scss" scoped>
.about-page {
  margin: 0 auto;

  text-align: center;
  padding-top: 20px;

  .feature {
    height: auto;

    img {
      display: block;
      width: 90%;
      margin: 0 auto;
      max-width: 350px;
    }
  }

  #add-friend {
    padding: 40px 0 0;
    height: 40px;
  }

  section {
    height: 300px;

    padding: 30px 30px;

    &.gray {
      background: #f8fafb;
    }
    &.white {
      background: white;
    }

    p.desc {
      margin-bottom: 30px;
      word-break: keep-all;
    }
  }
}
</style>
