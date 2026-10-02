import { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { api, type User } from '../lib/api';
import './Explore.css';

export function Explore() {
  const [users, setUsers] = useState<User[]>([]);
  const [filteredUsers, setFilteredUsers] = useState<User[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadUsers();
  }, []);

  useEffect(() => {
    if (searchQuery.trim() === '') {
      setFilteredUsers(users);
    } else {
      const query = searchQuery.toLowerCase();
      setFilteredUsers(
        users.filter(
          (user) =>
            user.username.toLowerCase().includes(query) ||
            user.email.toLowerCase().includes(query)
        )
      );
    }
  }, [searchQuery, users]);

  const loadUsers = async () => {
    try {
      const allUsers = await api.getAllUsers();
      // Get current user info to filter out own profile
      const currentUserStr = localStorage.getItem('current_user');
      const currentUserId = currentUserStr ? JSON.parse(currentUserStr).id : null;

      // Filter out current user
      const filteredByAuth = currentUserId
        ? allUsers.filter(user => user.id !== currentUserId)
        : allUsers;

      setUsers(filteredByAuth);
      setFilteredUsers(filteredByAuth);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load users');
    } finally {
      setLoading(false);
    }
  };

  const handleFollow = async (userId: number) => {
    try {
      await api.followUser(userId);
      // Reload users to get updated follow status
      await loadUsers();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to follow user');
    }
  };

  const handleUnfollow = async (userId: number) => {
    try {
      await api.unfollowUser(userId);
      // Reload users to get updated follow status
      await loadUsers();
    } catch (err) {
      alert(err instanceof Error ? err.message : 'Failed to unfollow user');
    }
  };

  if (loading) {
    return (
      <div className="explore-container">
        <div className="explore-loading">
          <div className="spinner"></div>
          <p>Loading users...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="explore-container">
        <div className="explore-error">
          <p>{error}</p>
          <button onClick={loadUsers} className="retry-button">
            Try Again
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="explore-container">
      <div className="explore-header">
        <h1>Explore Users</h1>
        <p>Discover and connect with other members of the community</p>
      </div>

      <div className="explore-search">
        <input
          type="text"
          className="search-input"
          placeholder="Search by username or email..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>

      {filteredUsers.length === 0 ? (
        <div className="explore-empty">
          <svg width="64" height="64" viewBox="0 0 24 24" fill="none">
            <path
              d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
            <circle
              cx="9"
              cy="7"
              r="4"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
            <path
              d="M23 21v-2a4 4 0 00-3-3.87M16 3.13a4 4 0 010 7.75"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            />
          </svg>
          <h2>No users found</h2>
          <p>Try a different search query</p>
        </div>
      ) : (
        <div className="users-grid">
          {filteredUsers.map((user) => (
            <div key={user.id} className="user-card">
              <div className="user-avatar-large">
                {user.username[0].toUpperCase()}
              </div>
              <h3>{user.username}</h3>
              <p className="user-card-email">{user.email}</p>
              <div className="user-card-actions">
                <Link to={`/profile/${user.id}`} className="view-profile-button">
                  View Profile
                </Link>
                {user.is_following ? (
                  <button
                    onClick={() => handleUnfollow(user.id)}
                    className="unfollow-button"
                  >
                    Unfollow
                  </button>
                ) : (
                  <button
                    onClick={() => handleFollow(user.id)}
                    className="follow-button"
                  >
                    Follow
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
