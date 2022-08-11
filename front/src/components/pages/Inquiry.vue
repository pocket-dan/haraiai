<script setup lang="ts">
import { reactive, ref } from 'vue';
import api from '@/api/api';

const error = reactive({
  isError: false,
  message: '',
});
const submitted = ref(false);
const inquiryText = ref('');

async function sendInquiry() {
  const text = inquiryText.value;
  if (text.length === 0) {
    error.isError = true;
    error.message = '内容を入力してください';
    return;
  }

  await api
    .sendInquiry(text)
    .then(() => {
      error.isError = false;
      submitted.value = true;
    })
    .catch((err: Error) => {
      error.isError = true;
      error.message =
        '送信に失敗しました。もう一度送信するか、時間を空けてお試しください。';
      console.error(err);
    });
}
</script>

<template>
  <div class="inquiry-page">
    <main>
      <h1>フィードバック</h1>
      <form v-if="!submitted">
        <label for="inquery">内容:</label>
        <span
          v-if="error.isError"
          class="error"
        >{{ error.message }}</span>
        <textarea
          id="inquery"
          v-model="inquiryText"
          :class="{ error: error.isError }"
          placeholder="要望・改善・感想など何でもお待ちしています"
        />
        <button
          type="button"
          @click="sendInquiry"
        >
          送る
        </button>
      </form>
      <p
        v-else
        id="done-message"
      >
        ありがとうございました！
      </p>
    </main>
  </div>
</template>

<style lang="scss" scoped>
.inquiry-page {
  main {
    max-width: 800px;
    margin: 0 auto;
    padding: 0 20px;

    h1 {
      margin-bottom: 40px;
    }

    label {
      display: block;
      margin-bottom: 10px;
    }

    textarea {
      font-size: 16px;
      width: 100%;
      display: block;
      height: 200px;
      padding: 12px 20px;
      box-sizing: border-box;
      border: 2px solid #ccc;
      border-radius: 4px;
      resize: none;
    }

    span.error {
      font-size: 0.9rem;
      color: red;
    }

    textarea.error {
      border-color: #ff8888;
    }

    button {
      margin: 20px auto 0px;
      appearance: none;
      background-color: #fafbfc;
      border: 1px solid rgba(27, 31, 35, 0.15);
      border-radius: 6px;
      box-shadow: rgba(27, 31, 35, 0.04) 0 1px 0,
        rgba(255, 255, 255, 0.25) 0 1px 0 inset;
      box-sizing: border-box;
      color: #24292e;
      cursor: pointer;
      display: block;
      line-height: 20px;
      list-style: none;
      padding: 6px 16px;
      transition: background-color 0.2s cubic-bezier(0.3, 0, 0.5, 1);
    }

    button:hover {
      background-color: #f3f4f6;
      text-decoration: none;
      transition-duration: 0.1s;
    }

    button:disabled {
      background-color: #fafbfc;
      border-color: rgba(27, 31, 35, 0.15);
      color: #959da5;
      cursor: default;
    }

    button:active {
      background-color: #edeff2;
      box-shadow: rgba(225, 228, 232, 0.2) 0 1px 0 inset;
      transition: none 0s;
    }

    button:focus {
      outline: 1px transparent;
    }

    button:before {
      display: none;
    }

    button:-webkit-details-marker {
      display: none;
    }

    #done-message {
      text-align: center;
    }
  }
}
</style>
