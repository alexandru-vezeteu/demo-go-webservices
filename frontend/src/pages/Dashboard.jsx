import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

export const Dashboard = () => {
  const { userInfo } = useAuth();

  const role = userInfo?.role?.toLowerCase();
  const isOwner = role === 'owner' || role === 'owner-event';
  const isClient = role === 'client';

  return (
    <div className="dashboard">
      <h1>Welcome, {userInfo?.email}!</h1>
      <p>Role: {isOwner ? 'Event Owner' : isClient ? 'Client' : userInfo?.role || 'Unknown'}</p>

      {isOwner && (
        <div className="dashboard-grid">
          <Link to="/events" className="dashboard-card">
            <h2>My Events</h2>
            <p>Create and manage your events</p>
          </Link>

          <Link to="/packets" className="dashboard-card">
            <h2>My Packages</h2>
            <p>Bundle events into ticket packages</p>
          </Link>

          <Link to="/customers" className="dashboard-card">
            <h2>Customers</h2>
            <p>View clients who purchased tickets</p>
          </Link>

          <Link to="/events-public" className="dashboard-card">
            <h2>Browse All Events</h2>
            <p>View public events from all owners</p>
          </Link>
        </div>
      )}

      {isClient && (
        <div className="dashboard-grid">
          <Link to="/tickets" className="dashboard-card">
            <h2>My Tickets</h2>
            <p>View and purchase event tickets</p>
          </Link>

          <Link to="/profile" className="dashboard-card">
            <h2>My Profile</h2>
            <p>Manage your personal information</p>
          </Link>

          <Link to="/events-public" className="dashboard-card">
            <h2>Browse Events</h2>
            <p>Discover events and packages</p>
          </Link>
        </div>
      )}

      <div style={{ marginTop: '40px', padding: '20px', background: '#f8f9fa', borderRadius: '8px' }}>
        <h3>Quick Stats</h3>
        <p>User ID: {userInfo?.user_id}</p>
        <p>Account Type: {isOwner ? 'Event Owner' : 'Client'}</p>
      </div>
    </div>
  );
};
