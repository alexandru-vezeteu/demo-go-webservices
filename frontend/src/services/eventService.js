import axios from 'axios';

const EVENT_SERVICE_URL = import.meta.env.VITE_EVENT_SERVICE_URL || 'http://localhost:12345';
const BASE_URL = `${EVENT_SERVICE_URL}/api/event-manager`;

const getAuthToken = () => {
    const token = localStorage.getItem('token');
    return token;
};

const createAuthHeaders = () => {
    const token = getAuthToken();
    return token ? {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
    } : {
        'Content-Type': 'application/json'
    };
};

export const eventService = {
    createEvent: async (event) => {
        const response = await axios.post(`${BASE_URL}/events`, event, {
            headers: createAuthHeaders()
        });
        return response.data.event || response.data;
    },

    getEventById: async (id) => {
        const response = await axios.get(`${BASE_URL}/events/${id}`, {
            headers: createAuthHeaders()
        });
        return response.data.event || response.data;
    },

    filterEvents: async (filters = {}) => {
        const response = await axios.get(`${BASE_URL}/events`, {
            params: filters,
            headers: createAuthHeaders()
        });
        return response.data;
    },

    updateEvent: async (id, updates) => {
        const response = await axios.patch(`${BASE_URL}/events/${id}`, updates, {
            headers: createAuthHeaders()
        });
        return response.data.event || response.data;
    },

    deleteEvent: async (id) => {
        const response = await axios.delete(`${BASE_URL}/events/${id}`, {
            headers: createAuthHeaders()
        });
        return response.data.event || response.data;
    },

    createPacket: async (packet) => {
        const response = await axios.post(`${BASE_URL}/event-packets`, packet, {
            headers: createAuthHeaders()
        });
        return response.data.event_packet || response.data;
    },

    getPacketById: async (id) => {
        const response = await axios.get(`${BASE_URL}/event-packets/${id}`, {
            headers: createAuthHeaders()
        });
        return response.data.event_packet || response.data;
    },

    updatePacket: async (id, updates) => {
        const response = await axios.patch(`${BASE_URL}/event-packets/${id}`, updates, {
            headers: createAuthHeaders()
        });
        return response.data.event_packet || response.data;
    },

    deletePacket: async (id) => {
        const response = await axios.delete(`${BASE_URL}/event-packets/${id}`, {
            headers: createAuthHeaders()
        });
        return response.data.event_packet || response.data;
    },

    createTicket: async (eventId, packetId) => {
        const data = {};
        if (eventId) data.event_id = eventId;
        if (packetId) data.packet_id = packetId;

        const response = await axios.post(`${BASE_URL}/tickets`, data, {
            headers: createAuthHeaders()
        });
        return response.data;
    },

    getTicketByCode: async (code) => {
        const response = await axios.get(`${BASE_URL}/tickets/${code}`, {
            headers: createAuthHeaders()
        });
        return response.data;
    },

    updateTicket: async (code, updates) => {
        const response = await axios.patch(`${BASE_URL}/tickets/${code}`, updates, {
            headers: createAuthHeaders()
        });
        return response.data;
    },

    deleteTicket: async (code) => {
        const response = await axios.delete(`${BASE_URL}/tickets/${code}`, {
            headers: createAuthHeaders()
        });
        return response.data;
    },

    createInclusion: async (eventId, packetId) => {
        const response = await axios.post(
            `${BASE_URL}/event-packet-inclusions/event/${eventId}/packet/${packetId}`,
            {},
            { headers: createAuthHeaders() }
        );
        return response.data;
    },

    getPacketsByEvent: async (eventId) => {
        const response = await axios.get(
            `${BASE_URL}/event-packet-inclusions/event/${eventId}`,
            { headers: createAuthHeaders() }
        );
        return response.data;
    },

    getEventsByPacket: async (packetId) => {
        const response = await axios.get(
            `${BASE_URL}/event-packet-inclusions/packet/${packetId}`,
            { headers: createAuthHeaders() }
        );
        return response.data;
    },

    deleteInclusion: async (eventId, packetId) => {
        const response = await axios.delete(
            `${BASE_URL}/event-packet-inclusions/event/${eventId}/packet/${packetId}`,
            { headers: createAuthHeaders() }
        );
        return response.data;
    },

    filterPackets: async (filters = {}) => {
        const response = await axios.get(`${BASE_URL}/event-packets`, {
            params: filters,
            headers: createAuthHeaders()
        });
        return response.data;
    },
};
