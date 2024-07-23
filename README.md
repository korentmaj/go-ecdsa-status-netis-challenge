# Go ECDSA Status Management API

## Project Overview

This project is a Go-based API that provides a robust mechanism for managing and verifying ECDSA (Elliptic Curve Digital Signature Algorithm) signed statuses. The application includes functionalities for creating, storing, retrieving, and verifying statuses in a PostgreSQL database, utilizing JSON Web Signatures (JWS) for security. The API is designed to serve JSON REST endpoints and includes basic authentication for securing certain operations.

## Features

- **ECDSA Key Generation and Management**: Create and manage ECDSA-P256 keys in PEM format.
- **Message Signing and Verification**: Sign messages with ECDSA keys and verify the signatures.
- **Status Management**: Store and manipulate statuses, each represented by a single bit in a byte array.
- **REST API**: A fully functional REST API for managing statuses, including endpoints for creation, retrieval, updating, and deletion.
- **PostgreSQL Integration**: Use PostgreSQL for persistent storage of statuses.
- **Basic Authentication**: Secure API endpoints with basic authentication.
- **JWS**: Utilize JSON Web Signatures to ensure the integrity and authenticity of the statuses.

## Technologies Used

- **Go**: The primary programming language used for building the application.
- **ECDSA**: Elliptic Curve Digital Signature Algorithm for cryptographic signing and verification.
- **PostgreSQL**: A powerful, open-source relational database system used for storing statuses.
- **JWS (JSON Web Signature)**: A compact, URL-safe means of representing claims to be transferred between two parties.
- **Gorilla Mux**: A powerful HTTP router and URL matcher for building Go web servers.
- **Docker**: (Optional) Containerization for easy setup and deployment.

## Setup and Installation

### Prerequisites

- Go 1.16 or higher
- PostgreSQL
- Git

### Step-by-Step Installation

1. **Clone the repository**:

    ```sh
    git clone https://github.com/korentmaj/go-ecdsa-status-netis-challenge.git
    cd go-ecdsa-status-netis-challenge
    ```

2. **Setup the Database**:

    Ensure PostgreSQL is installed and running. Create the necessary database and user by running the following SQL script:

    ```sql
    -- Create the database
    CREATE DATABASE ecdsadb;

    -- Connect to the newly created database
    \c ecdsadb;

    -- Create the user
    CREATE USER ecdsa_user WITH PASSWORD 'majk';

    -- Grant all privileges on the database to the user
    GRANT ALL PRIVILEGES ON DATABASE ecdsadb TO ecdsa_user;

    -- Reconnect to the ecdsadb database as ecdsa_user
    \c ecdsadb ecdsa_user;

    -- Create the statuses table
    CREATE TABLE statuses (
        id SERIAL PRIMARY KEY,
        status_id VARCHAR(255) NOT NULL,
        status_list BYTEA NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );

    -- Create an index on status_id for faster lookup
    CREATE INDEX idx_status_id ON statuses (status_id);

    -- Insert a new row into the statuses table
    INSERT INTO statuses (status_id, status_list, created_at, updated_at)
    VALUES ('testStatusId', '\x01', NOW(), NOW());
    ```

3. **Install Dependencies**:

    Ensure all dependencies are up-to-date:

    ```sh
    go mod tidy
    ```

4. **Build the Project**:

    ```sh
    go build -o server.exe ./cmd/server
    ```

5. **Run the Server**:

    ```sh
    server.exe
    ```

## Usage

### API Endpoints

#### 1. **Get Status**

    ```sh
    GET /api/status/{statusId}?index={index}
    ```

    **Example**:
    ```sh
    curl -u ecdsa_user:majk "http://localhost:8000/api/status/testStatusId?index=1"
    ```

#### 2. **Create a New Status**

    ```sh
    POST /api/status/{statusId}
    ```

#### 3. **Set Status**

    ```sh
    PUT /api/status/{statusId}?index={index}
    ```

#### 4. **Delete Status**

    ```sh
    DELETE /api/status/{statusId}?index={index}
    ```

#### 5. **Get All Status IDs**

    ```sh
    GET /api/status
    ```

#### 6. **Create New Structure**

    ```sh
    POST /api/status
    ```

