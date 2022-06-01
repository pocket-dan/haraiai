<script setup lang="ts">

type FAQ = {
  question: string
  answer: string
}

function questionIdOf(i: number): string {
  return "Q" + (i+1)
}

const FAQs: FAQ[] = [
  {
    question: "精算できますか？",
    answer: "haraiai をお使いのグループで「精算」と送ることで、精算することが可能です。",
  },
]

</script>

<template>
  <div class="help-page">
    <main>
      <h1>よくある質問</h1>

      <section id="toc">
        <h2>目次</h2>
        <ul>
          <a v-for="(e, i) in FAQs" :key="i" :href="'#' + questionIdOf(i)">
            <li >{{ e.question }}</li>
          </a>
        </ul>
      </section>

      <section id="faq">
        <h2>質問</h2>
        <p>下記に掲載されていない質問や要望がある方は <router-link to="/inquiry">こちら</router-link> よりお気軽にお問い合わせください</p>
        <div class="qa-list">
          <dl class="qa" v-for="(e, i) in FAQs" :key="i" :id="questionIdOf(i)">
            <dt><li>{{ e.question }}</li></dt>
            <dd>
              <p>{{ e.answer }}</p>
            </dd>
          </dl>
        </div>
      </section>
    </main>
  </div>
</template>

<style lang="scss" scoped>
.help-page {
  main {
    max-width: 800px;
    margin: 0 auto;
    padding: 0 40px;

    h2 {
      margin-top: 40px
    }
  }

  .qa-list {
    margin-top: 40px;

    .qa {
      margin-bottom: 30px;
    }
  }
}
</style>
