import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('authToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Handle auth errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('authToken');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth API calls
export const authAPI = {
  // Start OAuth flow
  initiateOAuth: () => api.get('/api/v1/auth/login'),
  
  // Get user profile
  getProfile: () => api.get('/api/v1/user/profile'),
};

// Songs API calls
export const songsAPI = {
  // Get user's songs
  getUserSongs: () => api.get('/api/v1/songs/'),
  
  // Upload a song
  uploadSong: (formData) => api.post('/api/v1/songs/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  }),
  
  // Get song details
  getSong: (id) => api.get(`/api/v1/songs/${id}`),
  
  // Stream song
  streamSong: (id) => `${API_BASE_URL}/api/v1/songs/${id}/stream`,
  
  // Delete song
  deleteSong: (id) => api.delete(`/api/v1/songs/${id}`),
};

export default api; 