<script setup lang="ts">
import BeMyFriendButton from "@/components/modules/BeMyFriendButton.vue"
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
      はじめ方は簡単。一緒に使う相手と haraiai をグループに追加するだけ。<br />
      面倒なアプリのインストールなく、 すぐにはじめられます。
      </p>
    </section>

    <section class="gray">
      <h2>支払った金額を自動で割り勘</h2>
      <p class="desc">
      お会計をまとめて払った人は、グループに タイトルと金額の2行のメッセージを投稿。<br />
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
}
</style>
