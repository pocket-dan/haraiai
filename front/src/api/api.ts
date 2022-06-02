import axios, {AxiosResponse} from 'axios';
import { BACKEND_API_BASE_URL } from '../config' // TODO: why '@/config' not work?

const instance = axios.create({
  baseURL: BACKEND_API_BASE_URL,
  timeout: 3000,
});

const sendInquiry = (text: string): Promise<AxiosResponse<void>> => {
  return instance.post("/NotifyInquiry", { text, })
}

export default {
  sendInquiry,
}
