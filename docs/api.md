# Chat Application - API Documentation

## Base URL

```
http://localhost:8080
```

## Authentication

JWT Bearer token in `Authorization` header:

```
Authorization: Bearer <JWT_TOKEN>
```

Token obtained via `/user/login` or `/user/register`. **Tokens expire in 1 hour.**

Alternatively, for WebSocket: `?token=<JWT_TOKEN>` query parameter.

---

## REST Endpoints

### Public Endpoints

#### 1. GET `/index`
Welcome endpoint (no auth required).

**Response:**
```json
{
  "message": "Welcome to the index page!"
}
```

---

#### 2. GET `/userList`
Retrieve all users in the system (no auth required).

**Response:**
```json
{
  "message": "User list retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Alice",
      "avatar_url": "/static/avatars/1_123456.jpg",
      "email": "alice@example.com",
      "phone": "+1234567890"
    },
    {
      "id": 2,
      "name": "Bob",
      "avatar_url": "/static/avatars/2_123457.jpg",
      "email": "bob@example.com",
      "phone": "+1234567891"
    }
  ]
}
```

---

#### 3. POST `/user/register`
Register a new user account.

**Request:**
```json
{
  "name": "Alice",
  "email": "alice@example.com",
  "phone": "+1234567890",
  "password": "securePassword123",
  "repassword": "securePassword123"
}
```

**Success Response (201/200):**
```json
{
  "message": "Register succeeded",
  "user_id": 1,
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Error Responses:**
- `400` - Password mismatch, missing field, email/phone already registered
- `500` - Internal error during registration

---

#### 4. POST `/user/login`
Authenticate user and receive JWT token.

**Request:**
```json
{
  "identifier": "alice@example.com",  // or phone number
  "password": "securePassword123"
}
```

**Success Response (200):**
```json
{
  "message": "Login succeeded",
  "user_id": 1,
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Error Responses:**
- `400` - Invalid request format
- `401` - Invalid credentials
- `500` - Token generation error

---

### Protected Endpoints (Requires JWT)

#### 5. GET `/user/me`
Get current authenticated user's profile.

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
```

**Success Response (200):**
```json
{
  "message": "ok",
  "data": {
    "id": 1,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T14:45:22Z",
    "name": "Alice",
    "avatar_url": "/static/avatars/1_1705330200.jpg",
    "email": "alice@example.com",
    "phone": "+1234567890",
    "is_logout": false
  }
}
```

**Error Responses:**
- `401` - Invalid/missing token
- `500` - User retrieval error

---

#### 6. GET `/messages`
Retrieve message history between two users.

**Query Parameters:**
- `with` (required): User ID to retrieve messages with
- `limit` (optional): Max messages to return (default: 100, max: 1000)

**Request Example:**
```
GET /messages?with=2&limit=50
Authorization: Bearer <JWT_TOKEN>
```

**Success Response (200):**
```json
{
  "message": "ok",
  "data": [
    {
      "id": 101,
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z",
      "from": 1,
      "to": 2,
      "type": "text",
      "body": "Hey Bob, how are you?",
      "delivered": true,
      "delivered_at": "2024-01-15T10:30:15Z"
    },
    {
      "id": 102,
      "created_at": "2024-01-15T10:31:00Z",
      "from": 2,
      "to": 1,
      "type": "text",
      "body": "Hi Alice! All good, thanks!",
      "delivered": true,
      "delivered_at": "2024-01-15T10:31:10Z"
    }
  ]
}
```

**Error Responses:**
- `400` - Missing or invalid `with` parameter
- `401` - Invalid token
- `500` - Database query error

---

#### 7. PUT `/user/{id}`
Full update of user profile (all fields provided).

**URL Parameters:**
- `id`: User ID to update

**Request:**
```json
{
  "name": "Alice Smith",
  "email": "alice.smith@example.com",
  "phone": "+1234567899",
  "password": "newPassword123"
}
```

**Success Response (200):**
```json
{
  "message": "Update User Succeeded!",
  "user_id": 1
}
```

**Error Responses:**
- `400` - Invalid request data or validation failed
- `404` - User not found
- `500` - Update error

**Notes:**
- Email/phone must be unique (except for current user)
- Password is hashed before storing
- Omitted fields are still updated to zero values

---

#### 8. PATCH `/user/{id}`
Partial update of user profile (only provided fields updated).

**URL Parameters:**
- `id`: User ID to update

**Request (example - update only name):**
```json
{
  "name": "Alice Johnson"
}
```

**Success Response (200):**
```json
{
  "message": "Update User Succeeded!",
  "user_id": 1
}
```

**Error Responses:**
- `400` - Invalid request or no fields to update
- `500` - Update error

**Notes:**
- Only provided fields are updated
- System fields (id, created_at) cannot be modified
- Email/phone uniqueness validated only if provided

---

#### 9. DELETE `/user/{id}`
Delete a user account.

**URL Parameters:**
- `id`: User ID to delete

**Success Response (200):**
```json
{
  "message": "Delete User Succeeded!",
  "user_id": 1
}
```

**Error Responses:**
- `400` - Invalid user ID format
- `404` - User not found
- `500` - Deletion error

**Notes:**
- Performs soft delete (sets `deleted_at` timestamp)
- User data retained in database

---

#### 10. POST `/user/avatar`
Upload avatar image for current user.

**Headers:**
```
Authorization: Bearer <JWT_TOKEN>
Content-Type: multipart/form-data
```

**Request (multipart):**
```
POST /user/avatar
Authorization: Bearer <JWT_TOKEN>

Form Data:
  - avatar: <FILE> (JPEG, PNG, WebP)
```

**Success Response (200):**
```json
{
  "message": "Avatar uploaded successfully",
  "avatar_url": "/static/avatars/1_1705330200.jpg"
}
```

**Error Responses:**
- `400` - Missing avatar file
- `401` - Invalid token
- `500` - File save or database update error

**Notes:**
- File saved to `/server/web/avatars/`
- Filename format: `{user_id}_{timestamp}.{ext}`
- URL automatically updated in user profile

---

## WebSocket Endpoint

### WebSocket `/ws`
Real-time bidirectional messaging via WebSocket.

**Connection URL:**
```
ws://localhost:8080/ws?token=<JWT_TOKEN>
// or
ws://localhost:8080/ws?user_id=<USER_ID>  (development only)
```

#### Authentication
- **JWT Token**: Preferred method (authorization header or query param)
- **User ID Query**: Fallback for development (no auth validation)

#### Message Format

All WebSocket messages are JSON objects:

```json
{
  "type": "message_type",
  "from": 1,
  "to": 2,
  "body": "message content",
  "id": 101,
  "room": "room_name"
}
```

#### Message Types

##### 1. Send Direct Message
**Client → Server:**
```json
{
  "type": "direct",
  "from": 1,
  "to": 2,
  "body": "Hello Bob!"
}
```

**Server processing:**
- Saves message to database
- Routes to recipient's connected clients
- Broadcasts delivery confirmation

**Server → Recipient Client:**
```json
{
  "type": "direct",
  "from": 1,
  "to": 2,
  "body": "Hello Bob!",
  "id": 101,
  "delivered": true,
  "delivered_at": "2024-01-15T10:30:15Z"
}
```

---

##### 2. Send Room Message
**Client → Server:**
```json
{
  "type": "room",
  "from": 1,
  "room": "engineering",
  "body": "Team standup at 10 AM"
}
```

**Server → All Room Members:**
```json
{
  "type": "room",
  "from": 1,
  "room": "engineering",
  "body": "Team standup at 10 AM",
  "id": 102,
  "delivered": true
}
```