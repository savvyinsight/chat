# Requirements

## Overview
Chat is a self-hosted, real-time messaging platform built with Go and React. The project provides REST APIs and WebSocket support for scalable, private in-app messaging features, suitable for teams and organizations requiring self-hosted communication infrastructure.

## Functional Requirements

### Authentication & User Management
- **User Registration**: Allow new users to create accounts via `POST /user/register` with credentials (email/username and password)
- **User Login**: Authenticate users via `POST /user/login` and return JWT tokens
- **JWT Token-Based Auth**: All protected endpoints require valid JWT tokens in request headers
- **User Profile**: Retrieve authenticated user information via `GET /user/me`
- **Avatar Management**: Support user avatar uploads via `POST /user/avatar` and static file serving

### Real-Time Messaging
- **WebSocket Connection**: Establish persistent WebSocket connections at `/ws` endpoint with JWT authentication
- **Message Routing**: Support 1-on-1 direct messages and group chat rooms
- **Message Persistence**: Store messages in MySQL database via `GET /messages` history API
- **Heartbeat/Ping**: Implement connection reliability checks to detect and handle stale connections
- **Client Reconnection Handling**: Gracefully handle client disconnections and reconnections

### Message Features
- **Send and Receive Messages**: Transmit messages in real-time between connected clients. Text,emojis,pictures,files,voice messages,video messages, etc.
- **Message History**: Retrieve historical messages with pagination support
- **Message Receipts**: Track message delivery and read status (planned)
- **File Attachments**: Support message attachments (planned)

### Scalability
- **Redis Pub/Sub**: Enable cross-instance message routing for horizontal scaling
- **Multi-Instance Support**: Design hub to work with distributed deployments
- **Group Rooms**: Route group messages to multiple recipients

## Technical Requirements

### Backend (Go)
- **Language**: Go 1.21 or higher
- **Framework**: Gin web framework for routing and middleware
- **Database**: MySQL for persistent user and message storage
- **ORM**: GORM for database abstraction
- **Caching/Message Queue**: Redis for pub/sub and horizontal scaling
- **Architecture**:
  - WebSocket hub for managing client connections
  - JWT middleware for authentication and authorization
  - CORS middleware for cross-origin requests
  - Request logging and error handling middleware
  - YAML-based configuration for database, Redis, and server settings
- **API Documentation**: Swagger/OpenAPI support

### Frontend (React)
- **Language**: TypeScript 5.0+
- **Framework**: React 18.2+
- **Build Tool**: Vite 5.0+
- **Styling**: Tailwind CSS 3.4+
- **Package Manager**: npm or yarn
- **Components**:
  - Login component for user authentication
  - Register component for account creation
  - Chat component for message display and sending
  - User profile component (Me)
- **Features**:
  - WebSocket client connection management
  - Real-time message updates
  - User authentication state management
  - Responsive UI design

### Infrastructure & Deployment
- **Database**: MySQL 5.7+ (or compatible)
- **Cache**: Redis 6.0+ (optional for scaling)
- **Server**: Linux/Docker compatible deployment
- **Configuration Management**: YAML-based config files for environment-specific settings
- **API Port**: Default 8080 (configurable)

## Non-Functional Requirements

### Performance
- Support real-time message delivery with minimal latency (<100 ms)
- Handle multiple concurrent WebSocket connections
- Optimize database queries for message history retrieval

### Reliability
- Implement connection heartbeat/ping mechanism
- Gracefully handle disconnections and reconnections
- Ensure message persistence in database
- Support horizontal scaling with Redis pub/sub

### Security
- Implement JWT token-based authentication
- Validate and sanitize all user inputs
- Enforce CORS policies
- Support secure WebSocket connections (WSS for production)
- Protect sensitive configuration (database credentials, secrets)

### Maintainability
- Modular code structure separating concerns (API, WebSocket, persistence)
- Clear configuration management
- API documentation with Swagger/OpenAPI
- Logging for debugging and monitoring

## API Endpoints

### User Management
- `POST /user/register` – User registration
- `POST /user/login` – User authentication
- `GET /user/me` – Get current user profile
- `POST /user/avatar` – Upload user avatar
- `GET /avatar/:id` – Retrieve user avatar

### Messaging
- `GET /messages` – Retrieve message history
- `WS /ws` – WebSocket connection for real-time messaging

## Browser & Client Support
- Modern browsers with WebSocket support (Chrome, Firefox, Safari, Edge)
- HTTPS/WSS for production deployments
- Mobile-responsive UI

## Development Requirements
- Version control with Git
- Unit and integration testing support
- Development server with hot reload (Vite)
- Database migration support (GORM)
- Swagger documentation generation

## Constraints & Assumptions
- Single or multi-instance deployment with optional Redis
- MySQL as primary data store
- JWT tokens stored client-side (localStorage/sessionStorage)
- YAML configuration files for server setup
- Go modules for dependency management
- npm/yarn for JavaScript dependencies
