import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080/api/v1';

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle errors
api.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  login: (credentials) => api.post('/auth/login', credentials),
  register: (userData) => api.post('/auth/register', userData),
};

// Users API
export const usersAPI = {
  getProfile: () => api.get('/users/profile'),
  updateProfile: (data) => api.put('/users/profile', data),
  getAnalytics: () => api.get('/users/analytics'),
};

// Teams API
export const teamsAPI = {
  getTeams: () => api.get('/teams'),
  getTeam: (id) => api.get(`/teams/${id}`),
  createTeam: (data) => api.post('/teams', data),
  updateTeam: (id, data) => api.put(`/teams/${id}`, data),
  deleteTeam: (id) => api.delete(`/teams/${id}`),
  addMember: (id, data) => api.post(`/teams/${id}/members`, data),
  removeMember: (teamId, userId) => api.delete(`/teams/${teamId}/members/${userId}`),
  getAnalytics: (id) => api.get(`/teams/${id}/analytics`),
  getMemberActivity: (id) => api.get(`/teams/${id}/analytics/members`),
};

// Templates API
export const templatesAPI = {
  getTemplates: () => api.get('/templates'),
  getTemplate: (id) => api.get(`/templates/${id}`),
  getTemplateCategories: (id) => api.get(`/templates/${id}/categories`),
};

// Retrospectives API
export const retrospectivesAPI = {
  getRetrospectives: () => api.get('/retrospectives'),
  getRetrospective: (id) => api.get(`/retrospectives/${id}`),
  createRetrospective: (data) => api.post('/retrospectives', data),
  updateRetrospective: (id, data) => api.put(`/retrospectives/${id}`, data),
  deleteRetrospective: (id) => api.delete(`/retrospectives/${id}`),
  startRetrospective: (id) => api.post(`/retrospectives/${id}/start`),
  endRetrospective: (id) => api.post(`/retrospectives/${id}/end`),
  addItem: (id, data) => api.post(`/retrospectives/${id}/items`, data),
  voteItem: (itemId) => api.post(`/retrospectives/items/${itemId}/vote`),
  addActionItem: (id, data) => api.post(`/retrospectives/${id}/action-items`, data),
  updateActionItem: (actionItemId, data) => api.put(`/retrospectives/action-items/${actionItemId}`, data),
  deleteActionItem: (actionItemId) => api.delete(`/retrospectives/action-items/${actionItemId}`),
  joinRetrospective: (id) => api.post(`/retrospectives/${id}/join`),
  getParticipants: (id) => api.get(`/retrospectives/${id}/participants`),
  deleteItem: (itemId) => api.delete(`/retrospectives/items/${itemId}`),
  reopenRetrospective: (id) => api.post(`/retrospectives/${id}/reopen`),
  createGroup: (id, data) => api.post(`/retrospectives/${id}/groups`, data),
  voteGroup: (groupId) => api.post(`/retrospectives/groups/${groupId}/vote`),
  deleteGroup: (groupId) => api.delete(`/retrospectives/groups/${groupId}`),
  mergeItems: (id, data) => api.post(`/retrospectives/${id}/merge-items`, data),
  toggleBlur: (id, blurred) => api.put(`/retrospectives/${id}/blur`, { blurred }),
  exportRetrospective: (id) => api.get(`/retrospectives/${id}/export`, { responseType: 'blob' }),
};


// WebSocket connection
export const createWebSocketConnection = (retrospectiveId, token) => {
  const wsUrl = `${process.env.REACT_APP_WS_URL || 'ws://localhost:8080'}/api/v1/ws/retrospective?retrospective_id=${retrospectiveId}`;
  
  const socket = new WebSocket(wsUrl);
  
  // Add auth header (WebSocket doesn't support custom headers, so we'll use query params)
  // In production, you might want to use a different approach
  return socket;
};

export default api;
