import React, { useState, useEffect } from 'react';
import { eventService } from '../services/eventService';
import { useAuth } from '../contexts/AuthContext';
import { parseErrorMessage } from '../utils/errorParser';

export const PacketsPage = () => {
  const { userInfo } = useAuth();
  const [packets, setPackets] = useState([]);
  const [formData, setFormData] = useState({
    name: '',
    location: '',
    description: '',
    allocated_seats: '',
  });
  const [editingId, setEditingId] = useState(null);
  const [message, setMessage] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadPackets();
  }, []);

  const loadPackets = async () => {
    setLoading(true);
    setMessage('');
    try {
      const response = await eventService.filterPackets();
      const packetList = response.event_packets || [];
      const ownerID = parseInt(userInfo?.user_id);
      const ownerPackets = packetList.filter(p => p.id_owner === ownerID);
      setPackets(ownerPackets);
    } catch (err) {
      console.error('Load packets error:', err.response?.status, err.response?.data?.error || err.message);
      const errorMsg = parseErrorMessage(err, 'Failed to load packets');
      setMessage(`Error: ${errorMsg}`);
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = async (e) => {
    e.preventDefault();
    setMessage('');

    try {
      await eventService.createPacket({
        name: formData.name,
        location: formData.location || undefined,
        description: formData.description || undefined,
        allocated_seats: formData.allocated_seats ? parseInt(formData.allocated_seats) : undefined,
        id_owner: parseInt(userInfo.user_id),
      });
      setMessage('Packet created successfully!');
      setFormData({ name: '', location: '', description: '', allocated_seats: '' });
      loadPackets();
    } catch (err) {
      console.error('Create packet error:', err.response?.status, err.response?.data?.error || err.message);
      const errorMsg = parseErrorMessage(err, 'Failed to create packet');
      setMessage(`Error: ${errorMsg}`);
    }
  };

  const handleUpdate = async (e) => {
    e.preventDefault();
    setMessage('');

    try {
      await eventService.updatePacket(editingId, {
        name: formData.name,
        location: formData.location || undefined,
        description: formData.description || undefined,
        allocated_seats: formData.allocated_seats ? parseInt(formData.allocated_seats) : undefined,
      });
      setMessage('Packet updated successfully!');
      setEditingId(null);
      setFormData({ name: '', location: '', description: '', allocated_seats: '' });
      loadPackets();
    } catch (err) {
      console.error('Update packet error:', err.response?.status, err.response?.data?.error || err.message);
      const errorMsg = parseErrorMessage(err, 'Failed to update packet');
      setMessage(`Error: ${errorMsg}`);
    }
  };

  const handleEdit = (packet) => {
    setEditingId(packet.id);
    setFormData({
      name: packet.name,
      location: packet.location || '',
      description: packet.description || '',
      allocated_seats: packet.allocated_seats ? packet.allocated_seats.toString() : '',
    });
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleCancelEdit = () => {
    setEditingId(null);
    setFormData({ name: '', location: '', description: '', allocated_seats: '' });
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this packet?')) return;

    try {
      await eventService.deletePacket(id);
      setMessage('Packet deleted successfully!');
      loadPackets();
    } catch (err) {
      console.error('Delete packet error:', err.response?.status, err.response?.data?.error || err.message);
      const errorMsg = parseErrorMessage(err, 'Failed to delete packet');
      setMessage(`Error: ${errorMsg}`);
    }
  };

  return (
    <div className="page-container">
      <h1>My Event Packets</h1>

      {message && <div className="message">{message}</div>}

      <div className="section">
        <h2>{editingId ? 'Edit Packet' : 'Create New Packet'}</h2>
        <p className="info-text">
          Event packets allow you to bundle multiple events together as a package deal.
        </p>

        <form onSubmit={editingId ? handleUpdate : handleCreate}>
          <div className="form-group">
            <label htmlFor="name">Packet Name *</label>
            <input
              id="name"
              type="text"
              placeholder="Enter packet name"
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
              placeholder="Packet location (optional)"
              value={formData.location}
              onChange={(e) => setFormData({ ...formData, location: e.target.value })}
            />
          </div>

          <div className="form-group">
            <label htmlFor="description">Description</label>
            <textarea
              id="description"
              placeholder="Packet description (optional)"
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              rows={4}
            />
          </div>

          <div className="form-group">
            <label htmlFor="allocated_seats">Allocated Seats</label>
            <input
              id="allocated_seats"
              type="number"
              placeholder="Number of seats allocated (optional)"
              value={formData.allocated_seats}
              onChange={(e) => setFormData({ ...formData, allocated_seats: e.target.value })}
              min="1"
            />
          </div>

          <div style={{ display: 'flex', gap: '10px' }}>
            <button type="submit" className="btn-primary">
              {editingId ? 'Update Packet' : 'Create Packet'}
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
        <h2>My Packets ({packets.length})</h2>
        {loading ? (
          <p>Loading packets...</p>
        ) : packets.length === 0 ? (
          <p>You haven't created any packets yet. Create your first packet above!</p>
        ) : (
          <div className="event-grid">
            {packets.map((packet) => (
              <div key={packet.id} className="card">
                <h3>{packet.name}</h3>
                <p><strong>Location:</strong> {packet.location || 'Not specified'}</p>
                <p><strong>Description:</strong> {packet.description || 'No description'}</p>
                <p><strong>Seats:</strong> {packet.allocated_seats || 'Not specified'}</p>
                <p className="event-id">ID: {packet.id}</p>

                <div style={{ display: 'flex', gap: '10px', marginTop: '10px' }}>
                  <button onClick={() => handleEdit(packet)} className="btn-secondary">
                    Edit
                  </button>
                  <button onClick={() => handleDelete(packet.id)} className="btn-danger">
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
