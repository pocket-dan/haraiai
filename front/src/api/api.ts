import axios, { AxiosResponse } from 'axios';
import { BACKEND_API_BASE_URL } from '@/config';

const instance = axios.create({
  baseURL: BACKEND_API_BASE_URL,
  timeout: 5000,
});

const sendInquiry = (text: string): Promise<AxiosResponse<void>> => {
  return instance.post('/NotifyInquiry', { text });
};

export default {
  sendInquiry,
};
