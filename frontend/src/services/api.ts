import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 60000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor — injeta user_id do localStorage
api.interceptors.request.use((config) => {
  const userId = localStorage.getItem('user_id')
  if (userId) {
    config.headers['X-User-ID'] = userId
  }
  return config
})

// Response interceptor — limpa user_id se 401
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('user_id')
    }
    return Promise.reject(error)
  },
)

export default api
