import { Outlet, Link, useLocation } from 'react-router-dom';
import { useAuth } from './lib/auth';
import './App.css';

function App() {
  const { isAuthenticated, logout } = useAuth();
  const location = useLocation();

  const isAuthPage = ['/login', '/register', '/confirm'].some(path =>
    location.pathname.startsWith(path)
  );

  return (
    <div className="app">
      {isAuthenticated && !isAuthPage && (
        <nav className="navbar">
          <div className="nav-container">
            <Link to="/feed" className="nav-brand">Social</Link>
            <div className="nav-links">
              <Link to="/feed" className="nav-link">Feed</Link>
              <Link to="/explore" className="nav-link">Explore</Link>
              <Link to="/create-post" className="nav-link">Create</Link>
              <button onClick={logout} className="nav-link logout-btn">Logout</button>
            </div>
          </div>
        </nav>
      )}
      <main>
        <Outlet />
      </main>
    </div>
  );
}

export default App;
