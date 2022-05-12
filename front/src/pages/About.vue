<script setup lang="ts">
import BeMyFriendButton from "@/components/BeMyFriendButton.vue";
</script>

<template>
  <div class="about-page">
    <section class="gray">
      <img class="logo" src="@/assets/img/logo.png" width="100" alt="Logo" />
      <h1>haraiai - 払い合い</h1>
      <p class="desc">2人で使える<br/>「精算しない」シンプルな 割り勘サービス</p>
      <BeMyFriendButton />
    </section>

    <section class="white">
      <h2>LINE上ですぐに始められる</h2>
      <p class="desc">
      はじめ方は簡単。グループを作って相手と haraiai 公式アカウントを追加するだけ。<br />
      面倒なアプリのインストールなく、 すぐに開始できます。
      </p>
    </section>

    <section class="gray">
      <h2>支払った金額を自動で割り勘</h2>
      <p class="desc">
      お会計をまとめて払った人は、グループに タイトルと金額 の2行のメッセージを投稿。<br />
      あとは haraiai が5:5で割り勘してくれます。
      </p>

    </section>

    <section class="white">
      <h2>退屈な精算はもうしない</h2>
      <p class="desc">
      その精算、本当に必要ですか？ <br />
      支払いが少ない側は、 次回多めに払うことで バランスを保ちましょう。
      </p>
    </section>

    <footer>
      <router-link to="/about">サービスについて</router-link>
      <router-link to="/help">ヘルプ</router-link>
      <router-link to="/inquiry">問い合わせ</router-link>

      <br />

      <p> Developed by <a href="https://github.com/raahii">raahii</a>. </p>
      <p class="credit">
        Special thanks to <br />
        <a href="https://www.flaticon.com/free-icons/donation" title="donation icons">Donation icons created by Freepik - Flaticon</a>
      </p>
    </footer>
  </div>
</template>

<style lang="scss" scoped>
.about-page {
  text-align: center;
  padding-top: 20px;

  section {
    height: 300px;
    padding: 30px 30px;

    &.gray {
      background: #F8FAFB;
    }
    &.white {
      background: white;
    }

    p.desc {
      margin-bottom: 30px;
      word-break: keep-all;
    }
  }

  footer {
    padding: 30px 0;
    display: flex;
    align-items: center;
    flex-direction: column;

    .credit {
      font-size: 0.8rem;
      color: gray;
    }
  }
}
</style>
