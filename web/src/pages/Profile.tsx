import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { api, type User } from '../lib/api';
import './Profile.css';

export function Profile() {
  const { userId } = useParams<{ userId: string }>();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useState(() => {
    if (!userId) return;

    const loadUser = async () => {
      try {
        const userData = await api.getUser(parseInt(userId));
        setUser(userData);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load user');
      } finally {
        setLoading(false);
      }
    };

    loadUser();
  });

  if (loading) {
    return (
      <div className="profile-container">
        <div className="profile-loading">
          <div className="spinner"></div>
          <p>Loading profile...</p>
        </div>
      </div>
    );
  }

  if (error || !user) {
    return (
      <div className="profile-container">
        <div className="profile-error">
          <p>{error || 'User not found'}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="profile-container">
      <div className="profile-card">
        <div className="profile-header">
          <div className="profile-avatar-large">
            {user.username[0].toUpperCase()}
          </div>
          <h1>{user.username}</h1>
          <p className="profile-email">{user.email}</p>
        </div>
      </div>
    </div>
  );
}
