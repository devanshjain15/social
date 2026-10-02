import { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { api } from '../lib/api';
import './ConfirmUser.css';

export function ConfirmUser() {
  const { token } = useParams<{ token: string }>();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    if (!token) {
      setError('Invalid activation link');
      setLoading(false);
      return;
    }

    const activateAccount = async () => {
      try {
        await api.activateUser(token);
        setSuccess(true);
        setTimeout(() => {
          navigate('/login');
        }, 3000);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Activation failed');
      } finally {
        setLoading(false);
      }
    };

    activateAccount();
  }, [token, navigate]);

  if (loading) {
    return (
      <div className="confirm-container">
        <div className="confirm-card">
          <div className="spinner"></div>
          <p>Activating your account...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="confirm-container">
        <div className="confirm-card">
          <div className="error-icon">✕</div>
          <h1>Activation Failed</h1>
          <p className="confirm-error">{error}</p>
          <Link to="/register" className="confirm-button">
            Back to Register
          </Link>
        </div>
      </div>
    );
  }

  if (success) {
    return (
      <div className="confirm-container">
        <div className="confirm-card">
          <div className="success-icon">✓</div>
          <h1>Account Activated!</h1>
          <p className="confirm-success">Your account has been successfully activated.</p>
          <p className="confirm-redirect">Redirecting to login...</p>
        </div>
      </div>
    );
  }

  return null;
}
