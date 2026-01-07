import React, { useState, useEffect } from 'react';
import { eventService } from '../services/eventService';
import { userService } from '../services/userService';
import { useAuth } from '../contexts/AuthContext';

export const UsersPage = () => {
  const { userInfo } = useAuth();
  const [activeTab, setActiveTab] = useState('events');
  const [myEvents, setMyEvents] = useState([]);
  const [myPackets, setMyPackets] = useState([]);
  const [eventCustomers, setEventCustomers] = useState({});
  const [packetCustomers, setPacketCustomers] = useState({});
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState('');

  useEffect(() => {
    loadOwnerData();
  }, []);

  const loadOwnerData = async () => {
    setLoading(true);
    setMessage('');
    try {
      const [eventResponse, packetResponse] = await Promise.all([
        eventService.filterEvents(),
        eventService.filterPackets(),
      ]);

      const eventList = eventResponse.events || [];
      const packetList = packetResponse.event_packets || [];

      const ownerID = parseInt(userInfo?.user_id);
      const ownerEvents = eventList.filter(e => e.id_owner === ownerID);
      const ownerPackets = packetList.filter(p => p.id_owner === ownerID);

      setMyEvents(ownerEvents);
      setMyPackets(ownerPackets);

      const eventCustomerMap = {};
      for (const event of ownerEvents) {
        try {
          const customers = await userService.getCustomersByEventID(event.id);
          eventCustomerMap[event.id] = customers;
        } catch (err) {
          console.error(`Error loading customers for event ${event.id}:`, err);
          eventCustomerMap[event.id] = [];
        }
      }
      setEventCustomers(eventCustomerMap);

      const packetCustomerMap = {};
      for (const packet of ownerPackets) {
        try {
          const customers = await userService.getCustomersByPacketID(packet.id);
          packetCustomerMap[packet.id] = customers;
        } catch (err) {
          console.error(`Error loading customers for packet ${packet.id}:`, err);
          packetCustomerMap[packet.id] = [];
        }
      }
      setPacketCustomers(packetCustomerMap);

      setMessage('Customer data loaded successfully');
    } catch (err) {
      setMessage(`Error loading data: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  const handleTabChange = (tab) => {
    setActiveTab(tab);
  };

  return (
    <div className="page-container">
      <h1>My Customers</h1>

      {message && <div className="message">{message}</div>}

      {myEvents.length === 0 && myPackets.length === 0 ? (
        <div className="info-banner">
          <p>You don't have any events or packets yet. Create events/packets to see customers who purchase tickets.</p>
        </div>
      ) : (
        <>
          <div className="tabs">
            <button
              className={`tab ${activeTab === 'events' ? 'active' : ''}`}
              onClick={() => handleTabChange('events')}
            >
              Event Customers
            </button>
            <button
              className={`tab ${activeTab === 'packets' ? 'active' : ''}`}
              onClick={() => handleTabChange('packets')}
            >
              Packet Customers
            </button>
          </div>

          {loading ? (
            <p>Loading...</p>
          ) : activeTab === 'events' ? (
            <div className="section">
              <h2>Customers by Event</h2>
              {myEvents.length === 0 ? (
                <p>You don't have any events yet.</p>
              ) : (
                myEvents.map((event) => (
                  <div key={event.id} className="section" style={{ borderBottom: '2px solid #ecf0f1', paddingBottom: '2rem', marginBottom: '2rem' }}>
                    <h3>{event.name}</h3>
                    <p><strong>Location:</strong> {event.location || 'Not specified'}</p>
                    <p><strong>Seats:</strong> {event.seats || 'Unlimited'}</p>

                    <h4 style={{ marginTop: '1rem', marginBottom: '0.5rem' }}>
                      Customers ({(eventCustomers[event.id] || []).length})
                    </h4>

                    {(eventCustomers[event.id] || []).length === 0 ? (
                      <p className="info-text">No tickets sold yet for this event.</p>
                    ) : (
                      <div className="user-grid">
                        {eventCustomers[event.id].map((customer) => (
                          <div key={customer.id} className="card">
                            <h3>{customer.first_name} {customer.last_name}</h3>
                            <p><strong>Email:</strong> {customer.email}</p>
                            {customer.social_media_links && (
                              <p><strong>Social:</strong> {customer.social_media_links}</p>
                            )}
                            <p className="event-id">Customer ID: {customer.id}</p>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                ))
              )}
            </div>
          ) : (
            <div className="section">
              <h2>Customers by Packet</h2>
              {myPackets.length === 0 ? (
                <p>You don't have any packets yet.</p>
              ) : (
                myPackets.map((packet) => (
                  <div key={packet.id} className="section" style={{ borderBottom: '2px solid #ecf0f1', paddingBottom: '2rem', marginBottom: '2rem' }}>
                    <h3>{packet.name}</h3>
                    <p><strong>Location:</strong> {packet.location || 'Not specified'}</p>
                    <p><strong>Allocated Seats:</strong> {packet.allocated_seats || 'Unlimited'}</p>

                    <h4 style={{ marginTop: '1rem', marginBottom: '0.5rem' }}>
                      Customers ({(packetCustomers[packet.id] || []).length})
                    </h4>

                    {(packetCustomers[packet.id] || []).length === 0 ? (
                      <p className="info-text">No tickets sold yet for this packet.</p>
                    ) : (
                      <div className="user-grid">
                        {packetCustomers[packet.id].map((customer) => (
                          <div key={customer.id} className="card">
                            <h3>{customer.first_name} {customer.last_name}</h3>
                            <p><strong>Email:</strong> {customer.email}</p>
                            {customer.social_media_links && (
                              <p><strong>Social:</strong> {customer.social_media_links}</p>
                            )}
                            <p className="event-id">Customer ID: {customer.id}</p>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                ))
              )}
            </div>
          )}
        </>
      )}
    </div>
  );
};
