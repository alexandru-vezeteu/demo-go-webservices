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
      const data = await eventService.filterEvents(searchFilters);
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
      const data = await eventService.filterPackets(searchFilters);
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

      if (activeTab === 'events') {
        const data = await eventService.filterEvents(params);
        setEventsData(data);
      } else {
        const data = await eventService.filterPackets(params);
        setPacketsData(data);
      }
    } catch (err) {
      setMessage(`Error navigating: ${err.message}`);
    } finally {
      setLoading(false);
    }
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

      {(links.first || links.prev || links.next) && (
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
              Page {metadata.page} ({items.length} items)
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
              Page {metadata.page}
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
        </div>
      )}
    </div>
  );
};
