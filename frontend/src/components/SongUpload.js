import React, { useState, useRef } from 'react';
import { songsAPI } from '../services/api';
import './SongUpload.css';

const SongUpload = ({ onUploadSuccess }) => {
  const [isDragging, setIsDragging] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const fileInputRef = useRef(null);

  const [metadata, setMetadata] = useState({
    title: '',
    artist: '',
    album: '',
    genre: '',
    year: '',
    trackNumber: ''
  });

  const handleDragOver = (e) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e) => {
    e.preventDefault();
    setIsDragging(false);
    
    const files = e.dataTransfer.files;
    if (files.length > 0) {
      handleFileSelect(files[0]);
    }
  };

  const handleFileSelect = (file) => {
    // Validate file type
    const allowedTypes = ['.mp3', '.flac', '.wav', '.m4a', '.ogg'];
    const fileExt = file.name.toLowerCase().substring(file.name.lastIndexOf('.'));
    
    if (!allowedTypes.includes(fileExt)) {
      setError('Invalid file type. Please select an audio file (MP3, FLAC, WAV, M4A, OGG)');
      return;
    }

    // Validate file size (100MB limit)
    if (file.size > 100 * 1024 * 1024) {
      setError('File too large. Maximum size is 100MB');
      return;
    }

    // Auto-fill title if empty
    if (!metadata.title) {
      setMetadata(prev => ({
        ...prev,
        title: file.name.replace(fileExt, '')
      }));
    }

    // Store file for upload
    setSelectedFile(file);
    setError('');
  };

  const [selectedFile, setSelectedFile] = useState(null);

  const handleUpload = async () => {
    if (!selectedFile) {
      setError('Please select a file to upload');
      return;
    }

    setUploading(true);
    setError('');
    setSuccess('');
    setUploadProgress(0);

    try {
      const formData = new FormData();
      formData.append('song', selectedFile);
      
      // Add metadata
      Object.keys(metadata).forEach(key => {
        if (metadata[key]) {
          formData.append(key, metadata[key]);
        }
      });

      const response = await songsAPI.uploadSong(formData);
      
      setSuccess('Song uploaded successfully!');
      setSelectedFile(null);
      setMetadata({
        title: '',
        artist: '',
        album: '',
        genre: '',
        year: '',
        trackNumber: ''
      });
      
      if (onUploadSuccess) {
        onUploadSuccess(response.data);
      }
    } catch (err) {
      setError(err.response?.data?.error || 'Upload failed. Please try again.');
      console.error('Upload error:', err);
    } finally {
      setUploading(false);
      setUploadProgress(0);
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setMetadata(prev => ({
      ...prev,
      [name]: value
    }));
  };

  return (
    <div className="song-upload">
      <h3>Upload New Song</h3>
      
      <div 
        className={`drop-zone ${isDragging ? 'dragging' : ''}`}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        onClick={() => fileInputRef.current?.click()}
      >
        <input
          ref={fileInputRef}
          type="file"
          accept=".mp3,.flac,.wav,.m4a,.ogg"
          onChange={(e) => e.target.files[0] && handleFileSelect(e.target.files[0])}
          style={{ display: 'none' }}
        />
        
        {selectedFile ? (
          <div className="file-selected">
            <p>Selected: {selectedFile.name}</p>
            <p>Size: {(selectedFile.size / (1024 * 1024)).toFixed(2)} MB</p>
          </div>
        ) : (
          <div className="drop-zone-content">
            <p>Drag & drop audio files here</p>
            <p>or click to browse</p>
            <p className="file-types">Supported: MP3, FLAC, WAV, M4A, OGG</p>
          </div>
        )}
      </div>

      {selectedFile && (
        <div className="metadata-form">
          <h4>Song Metadata</h4>
          
          <div className="form-row">
            <div className="form-group">
              <label>Title *</label>
              <input
                type="text"
                name="title"
                value={metadata.title}
                onChange={handleInputChange}
                placeholder="Song title"
                required
              />
            </div>
            
            <div className="form-group">
              <label>Artist</label>
              <input
                type="text"
                name="artist"
                value={metadata.artist}
                onChange={handleInputChange}
                placeholder="Artist name"
              />
            </div>
          </div>

          <div className="form-row">
            <div className="form-group">
              <label>Album</label>
              <input
                type="text"
                name="album"
                value={metadata.album}
                onChange={handleInputChange}
                placeholder="Album name"
              />
            </div>
            
            <div className="form-group">
              <label>Genre</label>
              <input
                type="text"
                name="genre"
                value={metadata.genre}
                onChange={handleInputChange}
                placeholder="Genre"
              />
            </div>
          </div>

          <div className="form-row">
            <div className="form-group">
              <label>Year</label>
              <input
                type="number"
                name="year"
                value={metadata.year}
                onChange={handleInputChange}
                placeholder="2024"
                min="1900"
                max="2030"
              />
            </div>
            
            <div className="form-group">
              <label>Track #</label>
              <input
                type="number"
                name="trackNumber"
                value={metadata.trackNumber}
                onChange={handleInputChange}
                placeholder="1"
                min="1"
              />
            </div>
          </div>

          <button
            onClick={handleUpload}
            disabled={uploading || !metadata.title}
            className="upload-btn"
          >
            {uploading ? 'Uploading...' : 'Upload Song'}
          </button>
        </div>
      )}

      {error && <p className="error-message">{error}</p>}
      {success && <p className="success-message">{success}</p>}
      
      {uploading && (
        <div className="upload-progress">
          <div className="progress-bar">
            <div 
              className="progress-fill" 
              style={{ width: `${uploadProgress}%` }}
            ></div>
          </div>
          <p>{uploadProgress}%</p>
        </div>
      )}
    </div>
  );
};

export default SongUpload; 