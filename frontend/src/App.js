import React, { useState, useEffect, useCallback } from 'react';
import { BrowserRouter as Router } from 'react-router-dom';
import OAuthLogin from './components/OAuthLogin';
import SongUpload from './components/SongUpload';
import CustomAudioPlayer from './components/CustomAudioPlayer';
import { authAPI, songsAPI } from './services/api';
import './App.css';

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [user, setUser] = useState(null);
  const [songs, setSongs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [currentSong, setCurrentSong] = useState(null);

  const loadUserSongs = useCallback(async () => {
    try {
      const response = await songsAPI.getUserSongs();
      setSongs(response.data);
    } catch (error) {
      console.error('Failed to load songs:', error);
    }
  }, []);

  const checkAuthStatus = useCallback(async () => {
    const token = localStorage.getItem('authToken');
    if (token) {
      try {
        const profileResponse = await authAPI.getProfile();
        setUser(profileResponse.data);
        setIsLoggedIn(true);
        loadUserSongs();
      } catch (error) {
        console.error('Auth check failed:', error);
        localStorage.removeItem('authToken');
        localStorage.removeItem('user');
      }
    }
    setLoading(false);
  }, [loadUserSongs]);

  useEffect(() => {
    checkAuthStatus();
  }, [checkAuthStatus]);

  const handleLoginSuccess = useCallback((userData) => {
    setUser(userData);
    setIsLoggedIn(true);
    loadUserSongs();
  }, [loadUserSongs]);

  const handleLogout = () => {
    localStorage.removeItem('authToken');
    localStorage.removeItem('user');
    setIsLoggedIn(false);
    setUser(null);
    setSongs([]);
  };

  const handlePlaySong = (song) => {
    setCurrentSong(song);
  };

  const handleUploadSuccess = (newSong) => {
    setSongs(prev => [newSong, ...prev]);
  };

  const handleDeleteSong = async (songId) => {
    try {
      await songsAPI.deleteSong(songId);
      setSongs(prev => prev.filter(song => song.id !== songId));
    } catch (error) {
      console.error('Failed to delete song:', error);
    }
  };

  // Handle OAuth callback
  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const token = urlParams.get('token');
    const userId = urlParams.get('user_id');
    const email = urlParams.get('email');
    const name = urlParams.get('name');

    if (token && userId) {
      // Store the token and user info
      localStorage.setItem('authToken', token);
      localStorage.setItem('user', JSON.stringify({
        id: parseInt(userId),
        email: email,
        name: name
      }));
      
      // Update state
      setUser({
        id: parseInt(userId),
        email: email,
        name: name
      });
      setIsLoggedIn(true);
      
      // Clean up URL
      window.history.replaceState({}, document.title, '/');
      
      // Load user songs
      loadUserSongs();
    }
  }, [loadUserSongs]);

  if (loading) {
    return (
      <div className="App">
        <div className="loading">
          <h2>Loading...</h2>
        </div>
      </div>
    );
  }

  if (!isLoggedIn) {
    return (
      <div className="App">
        <header className="App-header">
          <h1>🎵 Musica Universalis</h1>
          <p>Your Personal Music Cloud</p>
        </header>
        <main className="App-main">
          <OAuthLogin onLoginSuccess={handleLoginSuccess} />
        </main>
        <footer className="App-footer">
          <p>Built with ❤️ using Go, React, PostgreSQL, and MinIO</p>
        </footer>
      </div>
    );
  }

  return (
    <Router>
      <div className="App">
        <header className="App-header">
          <h1>🎵 Musica Universalis</h1>
          <p>Welcome, {user?.name || user?.email}!</p>
          <button onClick={handleLogout} className="logout-button">
            Logout
          </button>
        </header>

        <main className="App-main">
          <div className="dashboard">
            <div className="upload-section">
              <SongUpload onUploadSuccess={handleUploadSuccess} />
            </div>

            <div className="songs-section">
              <h3>Your Music Library ({songs.length} songs)</h3>
              {songs.length === 0 ? (
                <p className="no-songs">No songs uploaded yet. Start by uploading your first song!</p>
              ) : (
                <div className="songs-grid">
                  {songs.map(song => (
                    <div key={song.id} className="song-card">
                      <div className="song-info">
                        <h4>{song.title}</h4>
                        <p>{song.artist || 'Unknown Artist'}</p>
                        <p>{song.album || 'Unknown Album'}</p>
                        <p className="song-meta">
                          {song.duration ? `${Math.floor(song.duration / 60)}:${(song.duration % 60).toString().padStart(2, '0')}` : ''}
                          {song.file_size && ` • ${(song.file_size / (1024 * 1024)).toFixed(2)} MB`}
                        </p>
                      </div>
                      <div className="song-actions">
                        <button 
                          onClick={() => handlePlaySong(song)}
                          className="play-btn"
                        >
                          ▶️ Play
                        </button>
                        <button 
                          onClick={() => handleDeleteSong(song.id)}
                          className="delete-btn"
                        >
                          🗑️ Delete
                        </button>
                      </div>
                      {currentSong && currentSong.id === song.id && (
                        <CustomAudioPlayer 
                          songId={song.id} 
                          onEnded={() => setCurrentSong(null)} 
                        />
                      )}
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </main>

        <footer className="App-footer">
          <p>Built with ❤️ using Go, React, PostgreSQL, and MinIO</p>
        </footer>
      </div>
    </Router>
  );
}

export default App; 