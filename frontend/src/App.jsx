import React from 'react';
import { BrowserRouter, Routes, Route, Navigate, Link } from 'react-router-dom';
import { AuthProvider, useAuth } from './contexts/AuthContext';
import { HomePage } from './pages/HomePage';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';
import { PublicEventsPage } from './pages/PublicEventsPage';
import { Dashboard } from './pages/Dashboard';
import { UsersPage } from './pages/UsersPage';
import { EventsPage } from './pages/EventsPage';
import { PacketsPage } from './pages/PacketsPage';
import { TicketsPage } from './pages/TicketsPage';
import { ProfilePage } from './pages/ProfilePage';
import { InclusionsPage } from './pages/InclusionsPage';
import './App.css';

const ProtectedRoute = ({ children }) => {
  const { isAuthenticated } = useAuth();
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" />;
};

const Layout = ({ children }) => {
  const { isAuthenticated, logout, userInfo } = useAuth();
  const role = userInfo?.role?.toLowerCase();
  const isOwner = role === 'owner-event';
  const isClient = role === 'client';

  return (
    <div className="app">
      <nav className="navbar">
        <div className="nav-brand">
          <Link to="/" style={{ color: 'white', textDecoration: 'none' }}>
            POS Event Manager
          </Link>
        </div>

        <div className="nav-links">
          {!isAuthenticated && (
            <>
              <Link to="/events-public">Browse Events</Link>
            </>
          )}

          {isAuthenticated && (
            <>
              <Link to="/dashboard">Dashboard</Link>
              <Link to="/events-public">Browse Events</Link>

              {isOwner && (
                <>
                  <Link to="/events">My Events</Link>
                  <Link to="/packets">My Packets</Link>
                  <Link to="/inclusions">Bind Events</Link>
                  <Link to="/customers">Customers</Link>
                </>
              )}

              {isClient && (
                <>
                  <Link to="/tickets">My Tickets</Link>
                  <Link to="/profile">Profile</Link>
                </>
              )}
            </>
          )}
        </div>

        <div className="nav-user">
          {!isAuthenticated ? (
            <>
              <Link to="/login">
                <button className="btn-secondary">Login</button>
              </Link>
              <Link to="/register">
                <button className="btn-primary">Register</button>
              </Link>
            </>
          ) : (
            <>
              <span>{userInfo?.email} ({isOwner ? 'Owner' : 'Client'})</span>
              <button onClick={logout} className="btn-logout">Logout</button>
            </>
          )}
        </div>
      </nav>

      <main className="main-content">{children}</main>
    </div>
  );
};

function AppContent() {
  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/register" element={<RegisterPage />} />
          <Route path="/events-public" element={<PublicEventsPage />} />

          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            }
          />

          <Route
            path="/customers"
            element={
              <ProtectedRoute>
                <UsersPage />
              </ProtectedRoute>
            }
          />

          <Route
            path="/events"
            element={
              <ProtectedRoute>
                <EventsPage />
              </ProtectedRoute>
            }
          />

          <Route
            path="/packets"
            element={
              <ProtectedRoute>
                <PacketsPage />
              </ProtectedRoute>
            }
          />

          <Route
            path="/tickets"
            element={
              <ProtectedRoute>
                <TicketsPage />
              </ProtectedRoute>
            }
          />

          <Route
            path="/profile"
            element={
              <ProtectedRoute>
                <ProfilePage />
              </ProtectedRoute>
            }
          />

          <Route
            path="/inclusions"
            element={
              <ProtectedRoute>
                <InclusionsPage />
              </ProtectedRoute>
            }
          />
        </Routes>
      </Layout>
    </BrowserRouter>
  );
}

function App() {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  );
}

export default App;
