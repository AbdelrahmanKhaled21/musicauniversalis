import React, { useState } from 'react';
import './App.css';

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [user, setUser] = useState(null);

  const handleLogin = () => {
    // TODO: Implement OAuth login
    console.log('OAuth login will be implemented here');
    // For now, simulate login
    setIsLoggedIn(true);
    setUser({
      name: 'Test User',
      email: 'test@example.com'
    });
  };

  const handleLogout = () => {
    setIsLoggedIn(false);
    setUser(null);
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>🎵 Musica Universalis</h1>
        <p>Your Personal Music Cloud</p>
      </header>

      <main className="App-main">
        {!isLoggedIn ? (
          <div className="login-section">
            <h2>Welcome to Musica Universalis</h2>
            <p>Store, organize, and stream your music library</p>
            <button onClick={handleLogin} className="login-button">
              Login with GitHub
            </button>
            <p className="note">
              Note: OAuth integration coming soon. This is a placeholder.
            </p>
          </div>
        ) : (
          <div className="dashboard">
            <div className="user-info">
              <h3>Welcome, {user?.name}!</h3>
              <p>Email: {user?.email}</p>
              <button onClick={handleLogout} className="logout-button">
                Logout
              </button>
            </div>
            
            <div className="features">
              <h3>Features Coming Soon:</h3>
              <ul>
                <li>🎵 Upload and store music files</li>
                <li>📱 Stream music from anywhere</li>
                <li>📋 Create and manage playlists</li>
                <li>🔍 Search and organize your library</li>
              </ul>
            </div>
          </div>
        )}
      </main>

      <footer className="App-footer">
        <p>Built with Go, React, PostgreSQL, and MinIO</p>
      </footer>
    </div>
  );
}

export default App; 