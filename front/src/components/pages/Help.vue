<script setup lang="ts">
type FAQ = {
  question: string;
  answer: string;
};

function questionIdOf(i: number): string {
  return 'Q' + (i + 1);
}

const FAQs: FAQ[] = [
  {
    question: '精算はできますか？',
    answer:
      'はい、haraiai をお使いのグループで「精算」と送ることで可能です。しかし精算はおすすめしていません。理由は、支払いが少ない側が次回払うようにしてバランスを保つと面倒な精算作業はしなくて済み、お互いが楽だと考えているからです。',
  },
  {
    question: '割り勘の割合を5:5以外に変更することは出来ますか？',
    answer: '現在、haraiai ではその機能には対応しておりません。',
  },
  {
    question: '3人以上での割り勘は出来ますか？',
    answer: '現在、haraiai ではその機能には対応しておりません。',
  },
];
</script>

<template>
  <div class="help-page">
    <h1>よくある質問</h1>

    <section id="toc">
      <h2>目次</h2>
      <ul>
        <a
          v-for="(e, i) in FAQs"
          :key="i"
          :href="'#' + questionIdOf(i)"
        >
          <li>{{ e.question }}</li>
        </a>
      </ul>
    </section>

    <section id="faq">
      <h2>質問と回答</h2>
      <div class="qa-list">
        <dl
          v-for="(e, i) in FAQs"
          :id="questionIdOf(i)"
          :key="i"
          class="qa"
        >
          <dt>
            <li>{{ e.question }}</li>
          </dt>
          <dd>
            <p>{{ e.answer }}</p>
          </dd>
        </dl>
      </div>

      <p class="note">
        下記に掲載されていない質問や要望がある方は
        <router-link to="/inquiry">
          こちら
        </router-link>
        よりお気軽にお問い合わせください
      </p>
    </section>
  </div>
</template>

<style lang="scss" scoped>
.help-page {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 40px;

  h2 {
    margin-top: 60px;
  }

  .note {
    margin-top: 60px;
  }

  .qa-list {
    margin-top: 40px;

    .qa {
      margin-bottom: 30px;
    }
  }
}
</style>
