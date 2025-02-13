# Dating App API

This repository contains a Golang-based API for a dating application. The API is built using:

- **Go Kit** for service endpoints and transport.
- **GORM** for ORM and database interactions.
- **Gorilla Mux** for HTTP routing and dispatching incoming requests to their respective handlers.
- **Viper** for configuration management.
- **JWT** for authentication.

The service provides endpoints for user registration, authentication, and swipe actions with validations such as daily swipe limits and prevention of duplicate swipes within a day.

## Table of Contents

- [Project Structure](#project-structure)
- [Features](#features)
- [Technology Stack](#technology-stack)
- [Installation & Running](#installation--running)
- [Testing](#testing)
- [Additional Enhancements](#additional-enhancements)
- [License](#license)

## Project Structure

The project follows a layered architecture, separating concerns as follows:

![image](https://github.com/user-attachments/assets/3a4cdb66-67eb-4916-8063-bc3f68c86e78)

## Features

- **User Registration & Authentication:**  
  - Register new users with secure password hashing.
  - Log in to receive a JWT for authenticated access.

- **Swipe Functionality:**  
  - Perform swipe actions ("like" or "pass") on other profiles.
  - Validations include preventing duplicate swipes on the same profile in a day and enforcing a daily swipe limit.

- **JWT Authentication:**  
  - Endpoints are protected using JWT middleware that validates tokens and attaches user context.

- **Structured Logging:**  
  - Uses Go Kit's logging module for consistent, leveled logging.

- **Configuration:**  
  - Application settings are managed via YAML files and loaded using Viper.

## Technology Stack

- **Go Kit:** Provides modular, transport-agnostic endpoints.
- **GORM:** Manages database interactions with PostgreSQL or SQLite.
- **Gorilla Mux:** Routes and dispatches HTTP requests to appropriate handlers.
- **Viper:** Loads configuration files and manages environment settings.
- **JWT:** Implements secure authentication.

## Installation & Running

### Prerequisites

- Go 1.24 or later
- PostgreSQL or SQLite (configured via the YAML file)
- Git

### Steps

1. **Clone the Repository**

   ```bash
   git clone https://github.com/maulvi67/dating-app-api.git
   cd dating-app-api

2. **Configure the Application**
    Create or update your configuration file (e.g., config-dev.yaml) with your environment settings. For example:
    ```bash
    url:
      basepath: "/api"
      baseprefix: "/dating"
    
    server:
      port: 8080
      env: "development"
      log:
        level: "debug"
        output: "stdout"
        file-path: "./logs/app.log"
    
    security:
      jwt:
        jwt-secret: "your_jwt_secret_here"
        jwt-expire-hours: 24
    
    database:
      driver: "postgres"
      host: "localhost"
      port: 5432
      username: "dbuser"
      password: "dbpass"
      dbname: "dating_app_db"
      schemaname: "public"
      max-idle-connections: 20
      max-open-connections: 100
      connection-max-lifetime: "1200s"
      connection-max-idle-time: "1s"
      logger:
        level: "info"
        slow-threshold: "200ms"
        ignore-not-found: true
    
    app:
      swipe-limit: 50
    ```
   
3. **Install Dependencies**
   Use Go modules to install dependencies:
   ```bash
   go mod tidy
   ```

5. **Run the Service**
   ```bash
   go run main.go
   ```
   The API will be accessible at http://localhost:8080/dating-svc/api/v1... based on your configuration.

## Testing
The project includes both unit and integration tests.
- Unit Tests: Test individual service functions by mocking dependencies.
- Integration Tests: Use httptest to simulate HTTP requests and verify responses
To run tests, execute:
```
go test ./tests
```

## Additional Enhancements
- Deployment:
  - A Dockerfile is provided to containerize the application.
- The project uses golangci-lint for linting.
  - Code formatting is enforced with gofmt and goimports.
- Graceful Shutdown:
  - The HTTP server supports graceful shutdown to handle termination signals cleanly.
- Environment Management:
  - Viper handles configuration for different environments (development, staging, production).

## License
This project is licensed under the MIT License.


