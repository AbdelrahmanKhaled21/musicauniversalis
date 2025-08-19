# Musica Universalis 🎵

A full-stack music library platform built with Go, React, PostgreSQL, and MinIO in a dockerized environment.

## 🚀 Features

- **OAuth Authentication** with GitHub/Google
- **Music Upload & Storage** with chunked uploads
- **Audio Streaming** with range request support
- **Playlist Management** with batch operations
- **Background Processing** using Go workers and goroutines
- **Responsive Design** for mobile and desktop
- **Production Ready** with Docker and health checks

## 🏗️ Architecture

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐
│   React     │────│   Go API     │────│ PostgreSQL  │
│  Frontend   │    │   Backend    │    │  Database   │
└─────────────┘    └──────────────┘    └─────────────┘
                           │
                    ┌──────────────┐
                    │    MinIO     │
                    │ File Storage │
                    └──────────────┘
                           │
                    ┌──────────────┐
                    │ Go Workers   │
                    │ (Goroutines) │
                    │              │
                    │ • Metadata   │
                    │   Extraction │
                    │ • File       │
                    │   Processing │
                    │ • Background │
                    │   Tasks      │
                    └──────────────┘
```

## 🛠️ Tech Stack

- **Frontend**: React, Create React App
- **Backend**: Go with Gin framework, GORM ORM, golang-migrate
- **Database**: PostgreSQL
- **File Storage**: MinIO (S3-compatible)
- **Authentication**: OAuth 2.0, JWT tokens
- **Caching**: Go-native in-memory caching
- **Background Processing**: Go goroutines and worker pools
- **Containerization**: Docker Compose
- **Testing**: Go testing package

## 📋 Prerequisites

- Docker & Docker Compose
- Go 1.19+ (for local development)
- Node.js 16+ (for React development)
- Git

## 🚀 Quick Start

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd musica-universalis
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your OAuth credentials
   ```

3. **Start the services**
   ```bash
   docker-compose up -d
   ```

4. **Run database migrations**
   ```bash
   cd backend
   go run cmd/migrate/main.go
   ```

5. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - MinIO Console: http://localhost:9001

## 🔧 Development

### Backend Development
```bash
cd backend
go mod tidy
go run main.go
```

### Frontend Development
```bash
cd frontend
npm install
npm start
```

### Running Tests
```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test
```

## 📁 Project Structure

```
musica-universalis/
├── docker-compose.yml
├── .env
├── README.md
├── backend/
│   ├── Dockerfile
│   ├── go.mod
│   ├── main.go
│   ├── migrations/
│   ├── models/
│   ├── handlers/
│   ├── middleware/
│   ├── config/
│   ├── workers/
│   ├── utils/
│   └── errors/
├── frontend/
│   ├── Dockerfile
│   ├── package.json
│   ├── src/
│   └── public/
└── scripts/
```

## 🔐 OAuth Setup

### GitHub OAuth
1. Create a new OAuth app in GitHub
2. Set homepage URL to `http://localhost:3000`
3. Set callback URL to `http://localhost:8080/auth/callback`
4. Add client ID and secret to `.env`

### Google OAuth (Alternative)
1. Create OAuth 2.0 credentials in Google Cloud Console
2. Add authorized redirect URIs
3. Add client ID and secret to `.env`

## 🧪 Testing

- **Unit Tests**: Run after each feature implementation
- **Integration Tests**: Test API endpoints end-to-end
- **Coverage Goal**: >80% backend coverage

## 📊 API Documentation

API endpoints are available at `/api/v1/` with the following structure:

- **Authentication**: `/auth/*`
- **Music Library**: `/api/v1/songs/*`
- **Playlists**: `/api/v1/playlists/*`
- **Users**: `/api/v1/users/*`
- **Health**: `/health`

## 🚀 Deployment

The application is containerized and ready for deployment:

```bash
# Production build
docker-compose -f docker-compose.prod.yml up -d

# Environment-specific configs
cp .env.production .env
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📝 License

This project is for portfolio demonstration purposes.

## 🎯 Portfolio Highlights

This project demonstrates:
- Full-stack development with modern technologies
- Go concurrency and worker pool patterns
- Microservices architecture and containerization
- Authentication and security best practices
- File handling and media streaming
- Database design and migrations
- API design and REST principles
- Modern frontend with React
- Testing practices and code quality
- DevOps skills with Docker

---

Built with ❤️ using Go, React, and modern web technologies. 