import axios from 'axios';

const USER_SERVICE_URL = import.meta.env.VITE_USER_SERVICE_URL || 'http://localhost:12346';
const BASE_URL = `${USER_SERVICE_URL}/api/user-manager`;

const getAuthToken = () => {
    const token = localStorage.getItem('token');
    if (!token) {
        throw new Error('Authentication required. Please log in.');
    }
    return token;
};

const createAuthHeaders = () => ({
    'Authorization': `Bearer ${getAuthToken()}`,
    'Content-Type': 'application/json'
});

export const userService = {
    getUser: async (id) => {
        const response = await axios.get(`${BASE_URL}/users/${id}`, {
            headers: createAuthHeaders()
        });
        return response.data.user || response.data;
    },

    createUser: async (userData) => {
        const response = await axios.post(`${BASE_URL}/users`, userData, {
            headers: {
                'Content-Type': 'application/json'
            }
        });
        return response.data.user || response.data;
    },

    updateUser: async (id, updates) => {
        const response = await axios.patch(`${BASE_URL}/users/${id}`, updates, {
            headers: createAuthHeaders()
        });
        return response.data.user || response.data;
    },

    deleteUser: async (id) => {
        const response = await axios.delete(`${BASE_URL}/users/${id}`, {
            headers: createAuthHeaders()
        });
        return response.data.user || response.data;
    },

    createTicketForUser: async (userId, packetId, eventId) => {
        const data = {};
        if (packetId) data.packet_id = packetId;
        if (eventId) data.event_id = eventId;

        const response = await axios.post(`${BASE_URL}/clients/${userId}/tickets`, data, {
            headers: createAuthHeaders()
        });

        return response.data.ticket_code;
    },

    getCustomersByEventID: async (eventId) => {
        const response = await axios.get(`${BASE_URL}/events/${eventId}/customers`, {
            headers: createAuthHeaders()
        });
        return response.data.users || [];
    },

    getCustomersByPacketID: async (packetId) => {
        const response = await axios.get(`${BASE_URL}/packets/${packetId}/customers`, {
            headers: createAuthHeaders()
        });
        return response.data.users || [];
    },
};
