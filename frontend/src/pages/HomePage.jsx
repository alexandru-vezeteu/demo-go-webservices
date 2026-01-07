import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

export const HomePage = () => {
  const { isAuthenticated, userInfo } = useAuth();

  return (
    <div className="home-page">
      <div className="hero">
        <h1>Welcome to POS Event Manager</h1>
        <p>Your complete platform for managing and attending events</p>
      </div>

      {!isAuthenticated ? (
        <div className="cta-section">
          <h2>Get Started</h2>
          <div className="cta-cards">
            <div className="cta-card">
              <h3>Browse Events</h3>
              <p>Explore public events and ticket packages</p>
              <Link to="/events-public">
                <button className="btn-primary">View Events</button>
              </Link>
            </div>

            <div className="cta-card">
              <h3>Create Events</h3>
              <p>Register as an event owner to manage your events</p>
              <Link to="/register">
                <button className="btn-secondary">Register as Owner</button>
              </Link>
            </div>

            <div className="cta-card">
              <h3>Buy Tickets</h3>
              <p>Register as a client to purchase event tickets</p>
              <Link to="/register">
                <button className="btn-secondary">Register as Client</button>
              </Link>
            </div>
          </div>
        </div>
      ) : (
        <div className="welcome-back">
          <h2>Welcome back, {userInfo?.email}!</h2>
          <Link to="/dashboard">
            <button className="btn-primary">Go to Dashboard</button>
          </Link>
        </div>
      )}

      <div className="features">
        <h2>Features</h2>
        <div className="feature-grid">
          <div className="feature">
            <h3>Event Management</h3>
            <p>Create and manage events with ease</p>
          </div>
          <div className="feature">
            <h3>Ticket Packages</h3>
            <p>Bundle events into attractive packages</p>
          </div>
          <div className="feature">
            <h3>Customer Insights</h3>
            <p>View who's attending your events</p>
          </div>
          <div className="feature">
            <h3>🔒 Secure Tickets</h3>
            <p>Unique codes for each ticket purchase</p>
          </div>
        </div>
      </div>
    </div>
  );
};
