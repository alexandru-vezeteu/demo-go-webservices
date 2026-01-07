import React, { useState, useEffect } from 'react';
import { eventService } from '../services/eventService';
import { useAuth } from '../contexts/AuthContext';

export const EventsPage = () => {
  const { userInfo } = useAuth();
  const [events, setEvents] = useState([]);
  const [formData, setFormData] = useState({
    name: '',
    location: '',
    description: '',
    seats: '',
  });
  const [editingId, setEditingId] = useState(null);
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadEvents();
  }, []);

  const loadEvents = async () => {
    setLoading(true);
    try {
      const response = await eventService.filterEvents();
      const eventList = response.events || [];
      const myEvents = userInfo?.user_id
        ? eventList.filter(e => e.id_owner === parseInt(userInfo.user_id))
        : eventList;
      setEvents(myEvents);
    } catch (err) {
      setMessage(`Error loading events: ${err.message}`);
      setEvents([]);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = async (e) => {
    e.preventDefault();
    setMessage('');

    try {
      await eventService.createEvent({
        name: formData.name,
        location: formData.location || undefined,
        description: formData.description || undefined,
        seats: formData.seats ? parseInt(formData.seats) : undefined,
        id_owner: parseInt(userInfo.user_id),
      });
      setMessage('Event created successfully!');
      setFormData({ name: '', location: '', description: '', seats: '' });
      loadEvents();
    } catch (err) {
      setMessage(`Error: ${err.response?.data?.error || err.message}`);
    }
  };

  const handleUpdate = async (e) => {
    e.preventDefault();
    setMessage('');

    try {
      await eventService.updateEvent(editingId, {
        name: formData.name,
        location: formData.location || undefined,
        description: formData.description || undefined,
        seats: formData.seats ? parseInt(formData.seats) : undefined,
      });
      setMessage('Event updated successfully!');
      setEditingId(null);
      setFormData({ name: '', location: '', description: '', seats: '' });
      loadEvents();
    } catch (err) {
      setMessage(`Error: ${err.response?.data?.error || err.message}`);
    }
  };

  const handleEdit = (event) => {
    setEditingId(event.id);
    setFormData({
      name: event.name,
      location: event.location || '',
      description: event.description || '',
      seats: event.seats ? event.seats.toString() : '',
    });
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleCancelEdit = () => {
    setEditingId(null);
    setFormData({ name: '', location: '', description: '', seats: '' });
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this event?')) return;

    try {
      await eventService.deleteEvent(id);
      setMessage('Event deleted successfully!');
      loadEvents();
    } catch (err) {
      setMessage(`Error: ${err.response?.data?.error || err.message}`);
    }
  };

  return (
    <div className="page-container">
      <h1>My Events</h1>

      {message && <div className="message">{message}</div>}

      <div className="section">
        <h2>{editingId ? 'Edit Event' : 'Create New Event'}</h2>
        <form onSubmit={editingId ? handleUpdate : handleCreate}>
          <div className="form-group">
            <label htmlFor="name">Event Name *</label>
            <input
              id="name"
              type="text"
              placeholder="Enter event name"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="location">Location</label>
            <input
              id="location"
              type="text"
              placeholder="Event location (optional)"
              value={formData.location}
              onChange={(e) => setFormData({ ...formData, location: e.target.value })}
            />
          </div>

          <div className="form-group">
            <label htmlFor="description">Description</label>
            <textarea
              id="description"
              placeholder="Event description (optional)"
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              rows={4}
            />
          </div>

          <div className="form-group">
            <label htmlFor="seats">Total Seats</label>
            <input
              id="seats"
              type="number"
              placeholder="Number of seats (optional)"
              value={formData.seats}
              onChange={(e) => setFormData({ ...formData, seats: e.target.value })}
              min="1"
            />
          </div>

          <div style={{ display: 'flex', gap: '10px' }}>
            <button type="submit" className="btn-primary">
              {editingId ? 'Update Event' : 'Create Event'}
            </button>
            {editingId && (
              <button type="button" onClick={handleCancelEdit} className="btn-secondary">
                Cancel
              </button>
            )}
          </div>
        </form>
      </div>

      <div className="section">
        <h2>My Events ({events.length})</h2>
        {loading ? (
          <p>Loading events...</p>
        ) : events.length === 0 ? (
          <p>You haven't created any events yet. Create your first event above!</p>
        ) : (
          <div className="event-grid">
            {events.map((event) => (
              <div key={event.id} className="card">
                <h3>{event.name}</h3>
                <p><strong>Location:</strong> {event.location || 'Not specified'}</p>
                <p><strong>Description:</strong> {event.description || 'No description'}</p>
                <p><strong>Seats:</strong> {event.seats || 'Unlimited'}</p>
                <p className="event-id">ID: {event.id}</p>

                <div style={{ display: 'flex', gap: '10px', marginTop: '10px' }}>
                  <button onClick={() => handleEdit(event)} className="btn-secondary">
                    Edit
                  </button>
                  <button onClick={() => handleDelete(event.id)} className="btn-danger">
                    Delete
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};
