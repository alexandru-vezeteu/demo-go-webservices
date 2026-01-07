export const parseErrorMessage = (error, defaultMessage = 'An error occurred') => {
    if (!error.response) {
        if (error.message === 'Network Error') {
            return 'Unable to connect to the server. Please check your internet connection.';
        }
        return error.message || defaultMessage;
    }

    const status = error.response.status;
    const data = error.response.data;

    switch (status) {
        case 400:
            return parseBadRequest(data) || 'Invalid request. Please check your input.';
        case 401:
            return 'You need to log in to perform this action.';
        case 403:
            return 'You don\'t have permission to perform this action.';
        case 404:
            return parseNotFound(data) || 'The requested resource was not found.';
        case 409:
            return 'This item already exists or conflicts with existing data.';
        case 422:
            return parseValidationError(data) || 'Invalid data provided. Please check your input.';
        case 500:
            return 'Server error. Please try again later.';
        case 503:
            return 'Service temporarily unavailable. Please try again later.';
        default:
            break;
    }

    if (typeof data === 'string') {
        return cleanErrorMessage(data);
    }

    if (data?.error) {
        return cleanErrorMessage(data.error);
    }

    if (data?.message) {
        return cleanErrorMessage(data.message);
    }

    return defaultMessage;
};


const parseBadRequest = (data) => {
    if (typeof data === 'string' && data.includes('validation error')) {
        return parseValidationError(data);
    }
    return null;
};

const parseNotFound = (data) => {
    const message = typeof data === 'string' ? data : (data?.error || data?.message || '');

    if (message.includes('User') || message.includes('user')) {
        return 'User not found.';
    }
    if (message.includes('Event') || message.includes('event')) {
        return 'Event not found.';
    }
    if (message.includes('Packet') || message.includes('packet')) {
        return 'Event packet not found.';
    }
    if (message.includes('Ticket') || message.includes('ticket')) {
        return 'Ticket not found.';
    }
    if (message.match(/ID \d+ not found/)) {
        return 'The requested item was not found. It may have been deleted.';
    }

    return null;
};


const parseValidationError = (data) => {
    const message = typeof data === 'string' ? data : (data?.error || data?.message || '');

    let cleanMessage = message;
    const jsonMatch = message.match(/\{.*\}/);
    if (jsonMatch) {
        try {
            const parsed = JSON.parse(jsonMatch[0]);
            if (parsed.error) {
                cleanMessage = parsed.error;
            }
        } catch (e) {
        }
    }

    if (cleanMessage.includes('ID') && cleanMessage.includes('not found')) {
        if (message.includes('ticket')) {
            return 'Unable to create ticket: Event or packet not found.';
        }
        return 'The selected item no longer exists. Please refresh and try again.';
    }

    if (cleanMessage.includes('email') && cleanMessage.includes('required')) {
        return 'Email address is required.';
    }

    if (cleanMessage.includes('email') && cleanMessage.includes('invalid')) {
        return 'Please provide a valid email address.';
    }

    if (cleanMessage.includes('password') && cleanMessage.includes('required')) {
        return 'Password is required.';
    }

    if (cleanMessage.includes('already exists')) {
        return 'This item already exists.';
    }

    if (cleanMessage.includes('cannot be empty')) {
        return cleanMessage;
    }

    if (cleanMessage.includes('validation')) {
        return 'Please check your input and try again.';
    }

    return cleanErrorMessage(cleanMessage);
};


const cleanErrorMessage = (message) => {
    if (!message || typeof message !== 'string') {
        return '';
    }

    let cleaned = message.replace(/validation error on field '[^']+': ?/gi, '');

    cleaned = cleaned.replace(/failed to /gi, 'Unable to ');

    cleaned = cleaned.replace(/\{[^}]*"error"[^}]*\}/g, '');

    cleaned = cleaned.charAt(0).toUpperCase() + cleaned.slice(1);

    if (cleaned && !cleaned.match(/[.!?]$/)) {
        cleaned += '.';
    }

    return cleaned.trim();
};


export const formatErrorList = (errors) => {
    if (!errors || errors.length === 0) {
        return '';
    }

    if (errors.length === 1) {
        return errors[0];
    }

    return '• ' + errors.join('\n• ');
};
