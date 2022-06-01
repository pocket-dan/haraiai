let BACKEND_API_BASE_URL

if (import.meta.env.PROD) {
  BACKEND_API_BASE_URL = "https://asia-northeast1-haraiai.cloudfunctions.net"
} else {
  BACKEND_API_BASE_URL = "http://localhost:8080"
}

export {
  BACKEND_API_BASE_URL,
}
