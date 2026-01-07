import axios from 'axios';

const IDM_GATEWAY_URL = import.meta.env.VITE_IDM_GATEWAY_URL || 'http://localhost:8000';
const BASE_URL = `${IDM_GATEWAY_URL}/api/idm/auth`;

export const authService = {
    register: async (email, password, role) => {
        const response = await axios.post(`${BASE_URL}/register`, {
            email,
            password,
            role
        });
        return response.data;
    },

    login: async (email, password) => {
        const response = await axios.post(`${BASE_URL}/login`, {
            email,
            password
        });
        return response.data;
    },

    verifyToken: async (token) => {
        const response = await axios.post(`${BASE_URL}/verify`, {
            token
        });
        return response.data;
    },

    revokeToken: async (token) => {
        const response = await axios.post(`${BASE_URL}/revoke`, {
            token
        });
        return response.data;
    },
};

export const setAuthToken = (token) => {

};
