import React, { useState, useEffect } from 'react';
import { eventService } from '../services/eventService';
import { useAuth } from '../contexts/AuthContext';
import { parseErrorMessage } from '../utils/errorParser';

export const InclusionsPage = () => {
    const { userInfo } = useAuth();
    const [events, setEvents] = useState([]);
    const [packets, setPackets] = useState([]);
    const [selectedEventId, setSelectedEventId] = useState('');
    const [selectedPacketId, setSelectedPacketId] = useState('');
    const [eventPackets, setEventPackets] = useState({});
    const [message, setMessage] = useState('');
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        setLoading(true);
        try {
            const response = await eventService.filterEvents();
            const eventList = response.events || [];
            const myEvents = userInfo?.user_id
                ? eventList.filter(e => e.id_owner === parseInt(userInfo.user_id))
                : eventList;
            setEvents(myEvents);

            
            setPackets([]);

            setMessage('ℹ️ Enter Event ID and Packet ID below to bind them together');
        } catch (err) {
            console.error('Load data error:', err.response?.status, err.response?.data?.error || err.message);
            const errorMsg = parseErrorMessage(err, 'Failed to load data');
            setMessage(`❌ ${errorMsg}`);
            setEvents([]);
        } finally {
            setLoading(false);
        }
    };

    const loadEventPackets = async (eventId) => {
        try {
            const packets = await eventService.getPacketsByEvent(eventId);
            setEventPackets(prev => ({ ...prev, [eventId]: packets }));
        } catch (err) {
            console.error('Load event packets error:', err.response?.status, err.response?.data?.error || err.message);
        }
    };

    const handleAddInclusion = async (e) => {
        e.preventDefault();
        setMessage('');

        if (!selectedEventId || !selectedPacketId) {
            setMessage('⚠️ Please select both an event and a packet');
            return;
        }

        try {
            await eventService.createInclusion(parseInt(selectedEventId), parseInt(selectedPacketId));
            setMessage(`✓ Successfully added Event ${selectedEventId} to Packet ${selectedPacketId}!`);
            setSelectedEventId('');
            setSelectedPacketId('');

            loadEventPackets(parseInt(selectedEventId));
        } catch (err) {
            console.error('Create inclusion error:', err.response?.status, err.response?.data?.error || err.message);
            const errorMsg = parseErrorMessage(err, 'Failed to bind event to packet');
            setMessage(`❌ ${errorMsg}`);
        }
    };

    const handleRemoveInclusion = async (eventId, packetId) => {
        if (!window.confirm(`Remove Event ${eventId} from Packet ${packetId}?`)) return;

        try {
            await eventService.deleteInclusion(eventId, packetId);
            setMessage(`✓ Removed Event ${eventId} from Packet ${packetId}`);
            loadEventPackets(eventId);
        } catch (err) {
            console.error('Delete inclusion error:', err.response?.status, err.response?.data?.error || err.message);
            const errorMsg = parseErrorMessage(err, 'Failed to remove binding');
            setMessage(`❌ ${errorMsg}`);
        }
    };

    return (
        <div className="page-container">
            <h1>Event-Packet Bindings</h1>
            <p className="info-text">
                Bind events to packets to create ticket bundles. Customers can purchase a packet and get access to all included events.
            </p>

            {message && (
                <div
                    className="message"
                    style={{
                        padding: '15px',
                        margin: '10px 0',
                        borderRadius: '5px',
                        backgroundColor: message.includes('❌') || message.includes('⚠️')
                            ? '#fee'
                            : message.includes('✓')
                                ? '#efe'
                                : '#fff3cd',
                        border: message.includes('❌') || message.includes('⚠️')
                            ? '1px solid #fcc'
                            : message.includes('✓')
                                ? '1px solid #cfc'
                                : '1px solid #ffd700',
                        color: message.includes('❌') || message.includes('⚠️')
                            ? '#c00'
                            : message.includes('✓')
                                ? '#0a0'
                                : '#856404',
                        whiteSpace: 'pre-line'
                    }}
                >
                    {message}
                </div>
            )}

            <div className="section">
                <h2>Add Event to Packet</h2>
                <form onSubmit={handleAddInclusion}>
                    <div className="form-group">
                        <label htmlFor="event">Select Event</label>
                        {loading ? (
                            <p>Loading events...</p>
                        ) : events.length === 0 ? (
                            <p>No events found. Create an event first!</p>
                        ) : (
                            <select
                                id="event"
                                value={selectedEventId}
                                onChange={(e) => setSelectedEventId(e.target.value)}
                                required
                            >
                                <option value="">-- Select an Event --</option>
                                {events.map(event => (
                                    <option key={event.id} value={event.id}>
                                        ID {event.id}: {event.name} ({event.location || 'No location'})
                                    </option>
                                ))}
                            </select>
                        )}
                    </div>

                    <div className="form-group">
                        <label htmlFor="packet">Packet ID</label>
                        <input
                            id="packet"
                            type="number"
                            placeholder="Enter Packet ID"
                            value={selectedPacketId}
                            onChange={(e) => setSelectedPacketId(e.target.value)}
                            required
                            min="1"
                        />
                        <small style={{ color: '#666', display: 'block', marginTop: '5px' }}>
                            Enter the ID of the packet you created in the "My Packages" page
                        </small>
                    </div>

                    <button type="submit" className="btn-primary">
                        Bind Event to Packet
                    </button>
                </form>
            </div>

            <div className="section">
                <h2>Current Bindings</h2>
                {events.length === 0 ? (
                    <p>No events to display</p>
                ) : (
                    <div className="event-grid">
                        {events.map(event => (
                            <div key={event.id} className="card">
                                <h3>{event.name}</h3>
                                <p><strong>Event ID:</strong> {event.id}</p>
                                <p><strong>Location:</strong> {event.location || 'Not specified'}</p>

                                <div style={{ marginTop: '15px' }}>
                                    <button
                                        onClick={() => loadEventPackets(event.id)}
                                        className="btn-secondary"
                                        style={{ fontSize: '0.9em' }}
                                    >
                                        Load Packets for this Event
                                    </button>
                                </div>

                                {eventPackets[event.id] && (
                                    <div style={{ marginTop: '10px', padding: '10px', background: '#f5f5f5', borderRadius: '5px' }}>
                                        <strong>Included in Packets:</strong>
                                        {eventPackets[event.id].length === 0 ? (
                                            <p style={{ margin: '5px 0', fontSize: '0.9em' }}>Not in any packets yet</p>
                                        ) : (
                                            <ul style={{ margin: '5px 0', paddingLeft: '20px' }}>
                                                {eventPackets[event.id].map((pkt, idx) => (
                                                    <li key={idx} style={{ fontSize: '0.9em', margin: '5px 0' }}>
                                                        Packet ID: {pkt.packet_id || 'Unknown'}
                                                        <button
                                                            onClick={() => handleRemoveInclusion(event.id, pkt.packet_id)}
                                                            className="btn-danger"
                                                            style={{ marginLeft: '10px', fontSize: '0.8em', padding: '2px 8px' }}
                                                        >
                                                            Remove
                                                        </button>
                                                    </li>
                                                ))}
                                            </ul>
                                        )}
                                    </div>
                                )}
                            </div>
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
};
