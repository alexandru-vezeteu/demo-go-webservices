import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { userService } from '../services/userService';
import { parseErrorMessage } from '../utils/errorParser';

export const ProfilePage = () => {
    const { userInfo } = useAuth();
    const [userData, setUserData] = useState(null);
    const [editing, setEditing] = useState(false);
    const [formData, setFormData] = useState({
        first_name: '',
        last_name: '',
        social_media_links: '',
        first_name_private: false,
        last_name_private: false,
    });
    const [message, setMessage] = useState('');
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        loadUserData();
    }, []);

    const loadUserData = async () => {
        if (!userInfo?.user_id) return;

        setLoading(true);
        try {
            const user = await userService.getUser(parseInt(userInfo.user_id));
            setUserData(user);
            setFormData({
                first_name: user.first_name || '',
                last_name: user.last_name || '',
                social_media_links: user.social_media_links || '',
                first_name_private: user.first_name_private || false,
                last_name_private: user.last_name_private || false,
            });
            setMessage('');
        } catch (err) {
            console.error('Profile load error:', {
                status: err.response?.status,
                statusText: err.response?.statusText,
                data: err.response?.data,
                message: err.message
            });
            const errorMsg = parseErrorMessage(err, 'Failed to load profile');
            setMessage(`Error: ${errorMsg}`);
        } finally {
            setLoading(false);
        }
    };

    const handleUpdate = async (e) => {
        e.preventDefault();
        setMessage('');

        const errors = [];
        if (!formData.first_name || formData.first_name.trim() === '') {
            errors.push('First name cannot be empty');
        }
        if (!formData.last_name || formData.last_name.trim() === '') {
            errors.push('Last name cannot be empty');
        }

        if (errors.length > 0) {
            setMessage(`Please fix the following errors:\n${errors.join('\n')}`);
            return;
        }

        try {
            const updateData = {
                first_name: formData.first_name,
                last_name: formData.last_name,
                first_name_private: formData.first_name_private,
                last_name_private: formData.last_name_private,
            };

            if (formData.social_media_links && formData.social_media_links.trim() !== '') {
                updateData.social_media_links = formData.social_media_links;
            }

            await userService.updateUser(parseInt(userInfo.user_id), updateData);
            setMessage('Profile updated successfully!');
            setEditing(false);
            loadUserData();
        } catch (err) {
            console.error('Profile update error:', err.response?.status, err.response?.data?.error || err.message);
            const errorMsg = parseErrorMessage(err, 'Failed to update profile');
            setMessage(`Error: ${errorMsg}`);
        }
    };

    if (loading) {
        return (
            <div className="page-container">
                <p>Loading profile...</p>
            </div>
        );
    }

    return (
        <div className="page-container">
            <h1>My Profile</h1>

            {message && (
                <div
                    className="message"
                    style={{
                        padding: '15px',
                        margin: '10px 0',
                        borderRadius: '5px',
                        backgroundColor: message.startsWith('Error:')
                            ? '#fee'
                            : message.includes('successfully')
                                ? '#efe'
                                : '#fff3cd',
                        border: message.startsWith('Error:')
                            ? '1px solid #fcc'
                            : message.includes('successfully')
                                ? '1px solid #cfc'
                                : '1px solid #ffd700',
                        color: message.startsWith('Error:')
                            ? '#c00'
                            : message.includes('successfully')
                                ? '#0a0'
                                : '#856404',
                        whiteSpace: 'pre-line'
                    }}
                >
                    {message}
                </div>
            )}

            {!editing ? (
                <div className="profile-view">
                    <div className="user-details">
                        <p><strong>Email:</strong> {userData?.email}</p>
                        <p><strong>First Name:</strong> {userData?.first_name || 'Not set'}</p>
                        <p><strong>Last Name:</strong> {userData?.last_name || 'Not set'}</p>
                        <p><strong>Social Media:</strong> {userData?.social_media_links || 'Not set'}</p>
                        <p><strong>User ID:</strong> {userData?.id}</p>
                    </div>

                    <div style={{ marginTop: '20px' }}>
                        <button onClick={() => setEditing(true)} className="btn-primary">
                            Edit Profile
                        </button>
                    </div>

                    {userData?.ticket_list && userData.ticket_list.length > 0 && (
                        <div className="section">
                            <h2>My Tickets</h2>
                            <div className="ticket-grid">
                                {userData.ticket_list.map((ticket, idx) => (
                                    <div key={idx} className="card">
                                        <p><strong>Code:</strong> {ticket.code}</p>
                                        {ticket.event_id && <p><strong>Event ID:</strong> {ticket.event_id}</p>}
                                        {ticket.packet_id && <p><strong>Packet ID:</strong> {ticket.packet_id}</p>}
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}
                </div>
            ) : (
                <div className="profile-edit">
                    <form onSubmit={handleUpdate}>
                        <div className="form-group">
                            <label htmlFor="first_name">First Name</label>
                            <input
                                id="first_name"
                                type="text"
                                value={formData.first_name}
                                onChange={(e) => setFormData({ ...formData, first_name: e.target.value })}
                                required
                            />
                        </div>

                        <div className="form-group">
                            <label htmlFor="last_name">Last Name</label>
                            <input
                                id="last_name"
                                type="text"
                                value={formData.last_name}
                                onChange={(e) => setFormData({ ...formData, last_name: e.target.value })}
                                required
                            />
                        </div>

                        <div className="form-group">
                            <label htmlFor="social_media">Social Media Links</label>
                            <input
                                id="social_media"
                                type="text"
                                value={formData.social_media_links}
                                onChange={(e) => setFormData({ ...formData, social_media_links: e.target.value })}
                                placeholder="e.g., twitter.com/username"
                            />
                        </div>

                        <div className="section">
                            <h3>Privacy Settings</h3>
                            <p className="info-text">Control what event owners can see when you purchase tickets</p>

                            <div className="form-group" style={{ marginBottom: '10px' }}>
                                <label style={{ display: 'flex', alignItems: 'center', cursor: 'pointer' }}>
                                    <input
                                        type="checkbox"
                                        checked={formData.first_name_private}
                                        onChange={(e) => setFormData({ ...formData, first_name_private: e.target.checked })}
                                        style={{ marginRight: '10px', cursor: 'pointer' }}
                                    />
                                    Make first name private (owners will see "[Private]")
                                </label>
                            </div>

                            <div className="form-group">
                                <label style={{ display: 'flex', alignItems: 'center', cursor: 'pointer' }}>
                                    <input
                                        type="checkbox"
                                        checked={formData.last_name_private}
                                        onChange={(e) => setFormData({ ...formData, last_name_private: e.target.checked })}
                                        style={{ marginRight: '10px', cursor: 'pointer' }}
                                    />
                                    Make last name private (owners will see "[Private]")
                                </label>
                            </div>
                        </div>

                        <div style={{ display: 'flex', gap: '10px' }}>
                            <button type="submit" className="btn-primary">Save Changes</button>
                            <button type="button" onClick={() => setEditing(false)} className="btn-secondary">
                                Cancel
                            </button>
                        </div>
                    </form>
                </div>
            )}
        </div>
    );
};
