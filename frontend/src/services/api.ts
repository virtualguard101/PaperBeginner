import axios from 'axios'
import { useAuthStore } from '@/stores/authStore'

const api = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add auth token
api.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      useAuthStore.getState().logout()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// Auth API
export const authApi = {
  register: (data: { email: string; password: string; name: string }) =>
    api.post('/auth/register', data),
  
  login: (data: { email: string; password: string }) =>
    api.post('/auth/login', data),
  
  refresh: () =>
    api.post('/auth/refresh'),
}

// User API
export const userApi = {
  getMe: () => api.get('/users/me'),
  updateMe: (data: { name?: string; avatar?: string; preferences?: Record<string, unknown> }) =>
    api.put('/users/me', data),
  changePassword: (data: { old_password: string; new_password: string }) =>
    api.put('/users/me/password', data),
}

// Paper API
export const paperApi = {
  upload: (formData: FormData) =>
    api.post('/papers/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    }),
  list: (params?: { page?: number; per_page?: number }) =>
    api.get('/papers', { params }),
  get: (id: string) =>
    api.get(`/papers/${id}`),
  delete: (id: string) =>
    api.delete(`/papers/${id}`),
  analyze: (id: string, types?: string[]) =>
    api.post(`/papers/${id}/analyze`, { types }),
  getAnalyses: (id: string) =>
    api.get(`/papers/${id}/analyses`),
}

// Trending API
export const trendingApi = {
  getItems: (params?: { source?: string; category_id?: number; limit?: number }) =>
    api.get('/trending', { params }),
  getReports: (params?: { category_id?: number; limit?: number }) =>
    api.get('/trending/reports', { params }),
}

// Learning API
export const learningApi = {
  generate: (data: { category_id: number; difficulty: string; prerequisites?: string[]; focus_areas?: string[] }) =>
    api.post('/learning/generate', data),
  getPaths: () =>
    api.get('/learning/paths'),
  getPath: (id: string) =>
    api.get(`/learning/paths/${id}`),
}

// Review API
export const reviewApi = {
  generate: (data: { title: string; paper_ids: string[]; category_id?: number; style?: string }) =>
    api.post('/reviews/generate', data),
  list: (params?: { page?: number; per_page?: number }) =>
    api.get('/reviews', { params }),
  get: (id: string) =>
    api.get(`/reviews/${id}`),
  score: (id: string, content?: string) =>
    api.post(`/reviews/${id}/score`, { content }),
  delete: (id: string) =>
    api.delete(`/reviews/${id}`),
}

// CCF API
export const ccfApi = {
  getCategories: () =>
    api.get('/ccf/categories'),
}

export default api

