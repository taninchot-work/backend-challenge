# Backend Golang Coding Test

This project implements a RESTful API in Golang to manage a list of users, using MongoDB and JWT for authentication.

## Project Setup and Run Instructions

### Prerequisites

- Go (version 1.24.0 or higher)
- Docker and Docker Compose (if using Docker setup)
- MongoDB (if running locally without Docker)

### Local Setup (Without Docker)

1. **Clone the repository:**
   ```bash
   git clone git@github.com/taninchot-work/backend-challenge.git
   ```

   ```bash
   cd backend-challenge
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Configure environment variables:**
   Create a `config.yaml` file by copying `config.example.yaml` and update the necessary configurations, especially the
   database connection details and JWT secret.
   ```bash
   cp config.example.yaml config.yaml
   ```

4. **Ensure MongoDB is running and accessible.**

5. **Run the application:**
   ```bash
   go run cmd/main.go
   ```
   The API server will start, typically on port that you configure.

### Docker Setup

1. **Clone the repository:**
   ```bash
   git clone git clone github.com/taninchot-work/backend-challenge
   cd backend-challenge
   ```
2. **Build and run with Docker Compose:**
   ```bash
   docker-compose up --build
   ```

   This will start the API service and a MongoDB instance. on port that you configure

### API Endpoints

- **GET /api/v1/users/get/list**: List all users.
- **GET /api/v1/users/get/me**: Fetch user by ID (requires JWT).
- **POST /api/v1/users/register**: Register a new user.
- **POST /api/v1/users/login**: Authenticate user and return a JWT.
- **POST /api/v1/users/update**: Update a user's name or email (requires JWT).
- **POST /api/v1/users/delete**: Delete a user (requires JWT).

### Api Documentation

I provide a postman collection for testing the API endpoints. You can import it into Postman to test the API easily.

file: `postman_collection.json`

### Protected Endpoints

In the protected endpoints, you need to include the JWT in the Authorization header as follows:

```http
Authorization
Bearer <your_jwt_token>
```