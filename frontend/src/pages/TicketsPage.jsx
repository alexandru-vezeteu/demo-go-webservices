import React, { useState, useEffect } from 'react';
import { eventService } from '../services/eventService';
import { userService } from '../services/userService';
import { useAuth } from '../contexts/AuthContext';

export const TicketsPage = () => {
  const { userInfo } = useAuth();
  const [tickets, setTickets] = useState([]);
  const [enrichedTickets, setEnrichedTickets] = useState([]);
  const [events, setEvents] = useState([]);
  const [packets, setPackets] = useState([]);
  const [formData, setFormData] = useState({
    event_id: '',
    packet_id: '',
  });
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const [loadingTickets, setLoadingTickets] = useState(false);

  useEffect(() => {
    loadEvents();
    loadPackets();
    loadUserTickets();
  }, []);

  const loadEvents = async () => {
    try {
      const response = await eventService.filterEvents();
      setEvents(response.events || []);
    } catch (err) {
      console.error('Error loading events:', err);
      setEvents([]);
    }
  };

  const loadPackets = async () => {
    try {
      const response = await eventService.filterPackets();
      setPackets(response.event_packets || []);
    } catch (err) {
      console.error('Error loading packets:', err);
      setPackets([]);
    }
  };

  const loadUserTickets = async () => {
    if (!userInfo?.user_id) return;

    setLoadingTickets(true);
    try {
      const user = await userService.getUser(parseInt(userInfo.user_id));
      const ticketList = user.ticket_list || [];
      setTickets(ticketList);

      const enriched = await Promise.all(
        ticketList.map(async (ticket) => {
          const enrichedTicket = { ...ticket };

          const hasEventId = ticket.event_id !== undefined && ticket.event_id !== null;
          const hasPacketId = ticket.packet_id !== undefined && ticket.packet_id !== null;

          if (hasEventId) {
            try {
              const event = await eventService.getEventById(ticket.event_id);
              enrichedTicket.eventDetails = event;
            } catch (err) {
              console.error(`Error loading event ${ticket.event_id}:`, err);
            }
          }

          if (hasPacketId) {
            try {
              const packet = await eventService.getPacketById(ticket.packet_id);
              enrichedTicket.packetDetails = packet;

              if (packet) {
                try {
                  const packetEvents = await eventService.getEventsByPacket(ticket.packet_id);
                  enrichedTicket.packetEvents = packetEvents;
                } catch (err) {
                  console.error(`Error loading events for packet ${ticket.packet_id}:`, err);
                }
              }
            } catch (err) {
              console.error(`Error loading packet ${ticket.packet_id}:`, err);
            }
          }

          return enrichedTicket;
        })
      );

      setEnrichedTickets(enriched);
    } catch (err) {
      console.error('Error in loadUserTickets:', err);
      setMessage(`Error loading tickets: ${err.message}`);
    } finally {
      setLoadingTickets(false);
    }
  };

  const handlePurchase = async (e) => {
    e.preventDefault();
    if (!userInfo?.user_id) {
      setMessage('Error: Not logged in');
      return;
    }

    if (!formData.event_id && !formData.packet_id) {
      setMessage('Please select an event or packet');
      return;
    }

    setMessage('');
    setLoading(true);

    try {
      const ticketCode = await userService.createTicketForUser(
        parseInt(userInfo.user_id),
        formData.packet_id ? parseInt(formData.packet_id) : undefined,
        formData.event_id ? parseInt(formData.event_id) : undefined
      );

      setMessage(`Ticket purchased successfully! Code: ${ticketCode}`);
      setFormData({ event_id: '', packet_id: '' });
      loadUserTickets();
    } catch (err) {
      setMessage(`Error: ${err.response?.data?.error || err.message}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="page-container">
      <h1>My Tickets</h1>

      {message && <div className="message">{message}</div>}

      <div className="section">
        <h2>Purchase New Ticket</h2>
        <form onSubmit={handlePurchase}>
          <div className="form-group">
            <label htmlFor="event">Select Event</label>
            <select
              id="event"
              value={formData.event_id}
              onChange={(e) => setFormData({ ...formData, event_id: e.target.value, packet_id: '' })}
              disabled={loading || formData.packet_id}
            >
              <option value="">-- Select an Event --</option>
              {events.map((event) => (
                <option key={event.id} value={event.id}>
                  {event.name} - {event.location || 'No location'}
                </option>
              ))}
            </select>
          </div>

          <div className="form-group">
            <label htmlFor="packet">Or Select Packet</label>
            <select
              id="packet"
              value={formData.packet_id}
              onChange={(e) => setFormData({ ...formData, packet_id: e.target.value, event_id: '' })}
              disabled={loading || formData.event_id}
            >
              <option value="">-- Select a Packet --</option>
              {packets.map((packet) => (
                <option key={packet.id} value={packet.id}>
                  {packet.name} - {packet.location || 'No location'}
                </option>
              ))}
            </select>
          </div>

          <button type="submit" className="btn-primary" disabled={loading}>
            {loading ? 'Purchasing...' : 'Purchase Ticket'}
          </button>
        </form>
      </div>

      <div className="section">
        <h2>My Purchased Tickets ({tickets.length})</h2>
        {loadingTickets ? (
          <p>Loading ticket details...</p>
        ) : tickets.length === 0 ? (
          <p>You haven't purchased any tickets yet. Purchase your first ticket above!</p>
        ) : (
          <div className="ticket-grid">
            {enrichedTickets.map((ticket, idx) => (
              <div key={idx} className="card ticket-card">
                <h3>Ticket #{idx + 1}</h3>
                <p><strong>Code:</strong> <code>{ticket.code}</code></p>

                {ticket.event_id && (
                  <div style={{ marginTop: '15px', padding: '15px', background: '#f0f8ff', borderRadius: '8px', border: '1px solid #b3d9ff' }}>
                    <h4 style={{ marginTop: 0, color: '#0066cc' }}>Event Ticket</h4>

                    {ticket.eventDetails ? (
                      <>
                        <p><strong>Event Name:</strong> {ticket.eventDetails.name}</p>
                        {ticket.eventDetails.location && (
                          <p><strong>Location:</strong> {ticket.eventDetails.location}</p>
                        )}
                        {ticket.eventDetails.description && (
                          <p><strong>Description:</strong> {ticket.eventDetails.description}</p>
                        )}
                        {ticket.eventDetails.seats && (
                          <p><strong>Total Seats:</strong> {ticket.eventDetails.seats}</p>
                        )}
                        <p className="event-id" style={{ fontSize: '0.85em', color: '#666' }}>Event ID: {ticket.event_id}</p>
                      </>
                    ) : (
                      <p style={{ color: '#666', fontSize: '0.9em' }}>Loading event details...</p>
                    )}
                  </div>
                )}

                {ticket.packet_id && (
                  <div style={{ marginTop: '15px', padding: '15px', background: '#fff5e6', borderRadius: '8px', border: '1px solid #ffd699' }}>
                    <h4 style={{ marginTop: 0, color: '#cc8800' }}>Packet Ticket</h4>

                    {ticket.packetDetails ? (
                      <>
                        <p><strong>Packet Name:</strong> {ticket.packetDetails.name}</p>
                        {ticket.packetDetails.location && (
                          <p><strong>Location:</strong> {ticket.packetDetails.location}</p>
                        )}
                        {ticket.packetDetails.description && (
                          <p><strong>Description:</strong> {ticket.packetDetails.description}</p>
                        )}
                        {ticket.packetDetails.allocated_seats && (
                          <p><strong>Allocated Seats:</strong> {ticket.packetDetails.allocated_seats}</p>
                        )}
                        <p className="event-id" style={{ fontSize: '0.85em', color: '#666' }}>Packet ID: {ticket.packet_id}</p>

                        {ticket.packetEvents && ticket.packetEvents.length > 0 && (
                          <div style={{ marginTop: '10px', paddingTop: '10px', borderTop: '1px solid #ffd699' }}>
                            <p style={{ fontWeight: 'bold', marginBottom: '8px' }}>Included Events:</p>
                            <ul style={{ margin: '0', paddingLeft: '20px', fontSize: '0.9em' }}>
                              {ticket.packetEvents.map((evt, i) => (
                                <li key={i} style={{ marginBottom: '4px' }}>
                                  {evt.event_name || `Event ID: ${evt.event_id}`}
                                  {evt.event_location && ` - ${evt.event_location}`}
                                </li>
                              ))}
                            </ul>
                          </div>
                        )}
                      </>
                    ) : (
                      <p style={{ color: '#666', fontSize: '0.9em' }}>Loading packet details...</p>
                    )}
                  </div>
                )}

                {!ticket.event_id && !ticket.packet_id && (
                  <p style={{ color: '#999', fontStyle: 'italic', marginTop: '10px' }}>
                    No event or packet associated with this ticket.
                  </p>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};
