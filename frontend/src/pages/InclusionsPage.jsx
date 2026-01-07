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


            const packetResponse = await eventService.filterPackets();
            const packetList = packetResponse.event_packets || [];
            const myPackets = userInfo?.user_id
                ? packetList.filter(p => p.id_owner === parseInt(userInfo.user_id))
                : packetList;
            setPackets(myPackets);

            setMessage('Enter Event ID and Packet ID below to bind them together');
        } catch (err) {
            console.error('Load data error:', err.response?.status, err.response?.data?.error || err.message);
            const errorMsg = parseErrorMessage(err, 'Failed to load data');
            setMessage(`Error: ${errorMsg}`);
            setEvents([]);
        } finally {
            setLoading(false);
        }
    };

    const loadEventPackets = async (eventId) => {
        try {
            const response = await eventService.getPacketsByEvent(eventId);
            const packetsArray = Array.isArray(response) ? response : (response.event_packets || response.packets || []);
            setEventPackets(prev => ({ ...prev, [eventId]: packetsArray }));
        } catch (err) {
            console.error('Load event packets error:', err.response?.status, err.response?.data?.error || err.message);
            setEventPackets(prev => ({ ...prev, [eventId]: [] }));
        }
    };

    const handleAddInclusion = async (e) => {
        e.preventDefault();
        setMessage('');

        if (!selectedEventId || !selectedPacketId) {
            setMessage('Please select both an event and a packet');
            return;
        }

        try {
            await eventService.createInclusion(parseInt(selectedEventId), parseInt(selectedPacketId));
            setMessage(`Successfully added Event ${selectedEventId} to Packet ${selectedPacketId}!`);
            setSelectedEventId('');
            setSelectedPacketId('');

            loadEventPackets(parseInt(selectedEventId));
        } catch (err) {
            console.error('Create inclusion error:', err.response?.status, err.response?.data?.error || err.message);
            const errorMsg = parseErrorMessage(err, 'Failed to bind event to packet');
            setMessage(`Error: ${errorMsg}`);
        }
    };

    const handleRemoveInclusion = async (eventId, packetId) => {
        if (!window.confirm(`Remove Event ${eventId} from Packet ${packetId}?`)) return;

        try {
            await eventService.deleteInclusion(eventId, packetId);
            setMessage(`Removed Event ${eventId} from Packet ${packetId}`);
            loadEventPackets(eventId);
        } catch (err) {
            console.error('Delete inclusion error:', err.response?.status, err.response?.data?.error || err.message);
            const errorMsg = parseErrorMessage(err, 'Failed to remove binding');
            setMessage(`Error: ${errorMsg}`);
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
                        backgroundColor: message.startsWith('Error:')
                            ? '#fee'
                            : message.includes('Successfully') || message.includes('Removed')
                                ? '#efe'
                                : '#fff3cd',
                        border: message.startsWith('Error:')
                            ? '1px solid #fcc'
                            : message.includes('Successfully') || message.includes('Removed')
                                ? '1px solid #cfc'
                                : '1px solid #ffd700',
                        color: message.startsWith('Error:')
                            ? '#c00'
                            : message.includes('Successfully') || message.includes('Removed')
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
                        <label htmlFor="packet">Select Packet</label>
                        {loading ? (
                            <p>Loading packets...</p>
                        ) : packets.length === 0 ? (
                            <p>No packets found. Create a packet first!</p>
                        ) : (
                            <select
                                id="packet"
                                value={selectedPacketId}
                                onChange={(e) => setSelectedPacketId(e.target.value)}
                                required
                            >
                                <option value="">-- Select a Packet --</option>
                                {packets.map(packet => (
                                    <option key={packet.id} value={packet.id}>
                                        ID {packet.id}: {packet.name} ({packet.location || 'No location'})
                                    </option>
                                ))}
                            </select>
                        )}
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
                                                        <strong>{pkt.name || `Packet ${pkt.id}`}</strong>
                                                        {pkt.location && ` - ${pkt.location}`}
                                                        <span style={{ color: '#666', fontSize: '0.85em', marginLeft: '8px' }}>
                                                            (ID: {pkt.id})
                                                        </span>
                                                        <button
                                                            onClick={() => handleRemoveInclusion(event.id, pkt.id)}
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
