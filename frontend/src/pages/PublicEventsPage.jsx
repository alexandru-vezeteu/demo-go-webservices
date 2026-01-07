import React, { useState, useEffect } from 'react';
import { eventService } from '../services/eventService';
import { useAuth } from '../contexts/AuthContext';
import { Link } from 'react-router-dom';

export const PublicEventsPage = () => {
  const [activeTab, setActiveTab] = useState('events');
  const [eventsData, setEventsData] = useState({ events: [], _links: {}, _metadata: {} });
  const [packetsData, setPacketsData] = useState({ event_packets: [], _links: {}, _metadata: {} });
  const [filters, setFilters] = useState({
    name: '',
    location: '',
    min_seats: '',
    max_seats: '',
  });
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [selectedPacket, setSelectedPacket] = useState(null);
  const [packetEvents, setPacketEvents] = useState([]);
  const [loadingEvents, setLoadingEvents] = useState(false);
  const { isAuthenticated, userInfo } = useAuth();
  const isClient = userInfo?.role === 'client';

  useEffect(() => {
    if (activeTab === 'events') {
      loadEvents();
    } else {
      loadPackets();
    }
  }, [activeTab]);

  const loadEvents = async (searchFilters = {}) => {
    setLoading(true);
    setMessage('');
    try {
      const data = await eventService.filterEvents({ ...searchFilters, per_page: 2 });
      setEventsData(data);
    } catch (err) {
      setMessage(`Error loading events: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  const loadPackets = async (searchFilters = {}) => {
    setLoading(true);
    setMessage('');
    try {
      const data = await eventService.filterPackets({ ...searchFilters, per_page: 2 });
      setPacketsData(data);
    } catch (err) {
      setMessage(`Error loading packets: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (e) => {
    e.preventDefault();
    const searchFilters = {};
    if (filters.name) searchFilters.name = filters.name;
    if (filters.location) searchFilters.location = filters.location;
    if (filters.min_seats) searchFilters.min_seats = parseInt(filters.min_seats);
    if (filters.max_seats) searchFilters.max_seats = parseInt(filters.max_seats);

    if (activeTab === 'events') {
      loadEvents(searchFilters);
    } else {
      loadPackets(searchFilters);
    }
  };

  const handleClearFilters = () => {
    setFilters({
      name: '',
      location: '',
      min_seats: '',
      max_seats: '',
    });
    if (activeTab === 'events') {
      loadEvents();
    } else {
      loadPackets();
    }
  };

  const handleTabChange = (tab) => {
    setActiveTab(tab);
    setFilters({ name: '', location: '', min_seats: '', max_seats: '' });
    setMessage('');
  };

  const navigateToPage = async (linkUrl) => {
    if (!linkUrl) return;

    setLoading(true);
    setMessage('');
    try {
      const url = new URL(linkUrl);
      const params = {};
      url.searchParams.forEach((value, key) => {
        params[key] = value;
      });

      params.per_page = 2;

      if (activeTab === 'events') {
        const data = await eventService.filterEvents(params);
        setEventsData(data);
      } else {
        const data = await eventService.filterPackets(params);
        setPacketsData(data);
      }
    } catch (err) {
      console.error('Navigation error:', err);
      setMessage(`Error navigating: ${err.message}`);
    } finally {
      setLoading(false);
    }
  };

  const handleViewPacketEvents = async (packet) => {
    setSelectedPacket(packet);
    setShowModal(true);
    setLoadingEvents(true);
    setPacketEvents([]);

    try {
      const response = await eventService.getEventsByPacket(packet.id);
      const eventsArray = Array.isArray(response) ? response : (response?.events || response?.event_inclusions || []);
      setPacketEvents(eventsArray);
    } catch (err) {
      console.error('Error loading packet events:', err);
      setPacketEvents([]);
    } finally {
      setLoadingEvents(false);
    }
  };

  const closeModal = () => {
    setShowModal(false);
    setSelectedPacket(null);
    setPacketEvents([]);
  };

  const currentData = activeTab === 'events' ? eventsData : packetsData;
  const items = activeTab === 'events' ? (eventsData.events || []) : (packetsData.event_packets || []);
  const links = currentData._links || {};
  const metadata = currentData._metadata || {};

  return (
    <div className="page-container">
      <h1>Browse Events & Packets</h1>

      {!isAuthenticated && (
        <div className="info-banner">
          <p>Viewing as guest. <Link to="/login">Login</Link> or <Link to="/register">Register</Link> to purchase tickets.</p>
        </div>
      )}

      {message && <div className="message">{message}</div>}

      <div className="tabs">
        <button
          className={`tab ${activeTab === 'events' ? 'active' : ''}`}
          onClick={() => handleTabChange('events')}
        >
          Events
        </button>
        <button
          className={`tab ${activeTab === 'packets' ? 'active' : ''}`}
          onClick={() => handleTabChange('packets')}
        >
          Event Packets
        </button>
      </div>

      <div className="section">
        <h2>Search {activeTab === 'events' ? 'Events' : 'Packets'}</h2>
        <form onSubmit={handleSearch} className="search-form">
          <input
            type="text"
            placeholder={`${activeTab === 'events' ? 'Event' : 'Packet'} name...`}
            value={filters.name}
            onChange={(e) => setFilters({ ...filters, name: e.target.value })}
          />
          <input
            type="text"
            placeholder="Location..."
            value={filters.location}
            onChange={(e) => setFilters({ ...filters, location: e.target.value })}
          />
          <input
            type="number"
            placeholder="Min seats"
            value={filters.min_seats}
            onChange={(e) => setFilters({ ...filters, min_seats: e.target.value })}
          />
          <input
            type="number"
            placeholder="Max seats"
            value={filters.max_seats}
            onChange={(e) => setFilters({ ...filters, max_seats: e.target.value })}
          />
          <button type="submit" className="btn-primary">Search</button>
          <button type="button" onClick={handleClearFilters} className="btn-secondary">
            Clear
          </button>
        </form>
      </div>

      {(links.first || links.prev || links.next || links.last) && (
        <div className="pagination" style={{ display: 'flex', gap: '10px', justifyContent: 'center', marginBottom: '1rem' }}>
          {links.first && (
            <button
              onClick={() => navigateToPage(links.first.href)}
              className="btn-secondary"
              disabled={loading}
            >
              « First
            </button>
          )}
          {links.prev && (
            <button
              onClick={() => navigateToPage(links.prev.href)}
              className="btn-secondary"
              disabled={loading}
            >
              ‹ Previous
            </button>
          )}
          {metadata.page && (
            <span style={{ alignSelf: 'center', padding: '0 10px' }}>
              Page {metadata.page}{metadata.total_pages ? ` of ${metadata.total_pages}` : ''}
            </span>
          )}
          {links.next && (
            <button
              onClick={() => navigateToPage(links.next.href)}
              className="btn-secondary"
              disabled={loading}
            >
              Next ›
            </button>
          )}
          {links.last && (
            <button
              onClick={() => navigateToPage(links.last.href)}
              className="btn-secondary"
              disabled={loading}
            >
              Last »
            </button>
          )}
        </div>
      )}

      <div className="section">
        <h2>{activeTab === 'events' ? 'Available Events' : 'Available Packets'} ({items.length})</h2>

        {loading ? (
          <p>Loading...</p>
        ) : items.length === 0 ? (
          <p>No {activeTab === 'events' ? 'events' : 'packets'} found. Try adjusting your filters.</p>
        ) : (
          <div className="event-grid">
            {items.map((item) => (
              <div key={item.id} className="card event-card">
                <h3>{item.name}</h3>
                <p><strong>Location:</strong> {item.location || 'Not specified'}</p>
                <p><strong>Description:</strong> {item.description || 'No description'}</p>
                <p><strong>{activeTab === 'events' ? 'Available Seats' : 'Allocated Seats'}:</strong> {item.seats || item.allocated_seats || 'Unlimited'}</p>
                <p className="event-id">{activeTab === 'events' ? 'Event' : 'Packet'} ID: {item.id}</p>

                {activeTab === 'packets' && (
                  <button
                    className="btn-secondary"
                    onClick={() => handleViewPacketEvents(item)}
                    style={{ marginBottom: '10px', width: '100%' }}
                  >
                    View Events Inside
                  </button>
                )}

                {isClient && (
                  <Link to="/tickets">
                    <button className="btn-primary">Purchase Ticket</button>
                  </Link>
                )}
              </div>
            ))}
          </div>
        )}
      </div>

      {(links.first || links.prev || links.next) && items.length > 0 && (
        <div className="pagination" style={{ display: 'flex', gap: '10px', justifyContent: 'center', marginTop: '2rem' }}>
          {links.first && (
            <button
              onClick={() => navigateToPage(links.first.href)}
              className="btn-secondary"
              disabled={loading}
            >
              « First
            </button>
          )}
          {links.prev && (
            <button
              onClick={() => navigateToPage(links.prev.href)}
              className="btn-secondary"
              disabled={loading}
            >
              ‹ Previous
            </button>
          )}
          {metadata.page && (
            <span style={{ alignSelf: 'center', padding: '0 10px' }}>
              Page {metadata.page}{metadata.total_pages ? ` of ${metadata.total_pages}` : ''}
            </span>
          )}
          {links.next && (
            <button
              onClick={() => navigateToPage(links.next.href)}
              className="btn-secondary"
              disabled={loading}
            >
              Next ›
            </button>
          )}
          {links.last && (
            <button
              onClick={() => navigateToPage(links.last.href)}
              className="btn-secondary"
              disabled={loading}
            >
              Last »
            </button>
          )}
        </div>
      )}

      {/* Modal for viewing packet events */}
      {showModal && (
        <div
          className="modal-overlay"
          onClick={closeModal}
          style={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            backgroundColor: 'rgba(0, 0, 0, 0.5)',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            zIndex: 1000
          }}
        >
          <div
            className="modal-content"
            onClick={(e) => e.stopPropagation()}
            style={{
              backgroundColor: 'white',
              padding: '30px',
              borderRadius: '10px',
              maxWidth: '600px',
              maxHeight: '80vh',
              overflow: 'auto',
              boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)'
            }}
          >
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
              <h2 style={{ margin: 0 }}>{selectedPacket?.name}</h2>
              <button
                onClick={closeModal}
                style={{
                  background: 'none',
                  border: 'none',
                  fontSize: '24px',
                  cursor: 'pointer',
                  padding: '0',
                  color: '#666'
                }}
              >
                ×
              </button>
            </div>

            {selectedPacket && (
              <div style={{ marginBottom: '20px', padding: '15px', backgroundColor: '#f8f9fa', borderRadius: '5px' }}>
                <p><strong>Location:</strong> {selectedPacket.location || 'Not specified'}</p>
                <p><strong>Description:</strong> {selectedPacket.description || 'No description'}</p>
                <p><strong>Allocated Seats:</strong> {selectedPacket.allocated_seats || 'Unlimited'}</p>
              </div>
            )}

            <h3>Included Events</h3>
            {loadingEvents ? (
              <p>Loading events...</p>
            ) : packetEvents.length === 0 ? (
              <p style={{ color: '#666', fontStyle: 'italic' }}>No events included in this packet yet.</p>
            ) : (
              <div style={{ display: 'flex', flexDirection: 'column', gap: '15px' }}>
                {packetEvents.map((event, idx) => (
                  <div
                    key={idx}
                    style={{
                      padding: '15px',
                      border: '1px solid #dee2e6',
                      borderRadius: '5px',
                      backgroundColor: '#fff'
                    }}
                  >
                    <h4 style={{ marginTop: 0, color: '#2c3e50' }}>
                      {event.name || `Event ${event.id}`}
                    </h4>
                    {event.location && (
                      <p style={{ margin: '5px 0', color: '#666' }}>
                        <strong>Location:</strong> {event.location}
                      </p>
                    )}
                    {event.description && (
                      <p style={{ margin: '5px 0', color: '#666' }}>
                        <strong>Description:</strong> {event.description}
                      </p>
                    )}
                    {event.seats && (
                      <p style={{ margin: '5px 0', color: '#666' }}>
                        <strong>Seats:</strong> {event.seats}
                      </p>
                    )}
                    <p style={{ margin: '5px 0', fontSize: '0.9em', color: '#999' }}>Event ID: {event.id}</p>
                  </div>
                ))}
              </div>
            )}

            <button
              onClick={closeModal}
              className="btn-secondary"
              style={{ marginTop: '20px', width: '100%' }}
            >
              Close
            </button>
          </div>
        </div>
      )}
    </div>
  );
};
