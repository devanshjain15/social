import { Link } from 'react-router-dom';
import './Home.css';

export function Home() {
  return (
    <div className="home-container">
      <div className="home-hero">
        <div className="home-content">
          <h1 className="home-title">
            Welcome to <span className="gradient-text">Social</span>
          </h1>
          <p className="home-subtitle">
            Connect, share, and engage with your community
          </p>
          <div className="home-buttons">
            <Link to="/register" className="btn btn-primary">
              Get Started
            </Link>
            <Link to="/login" className="btn btn-secondary">
              Login
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
