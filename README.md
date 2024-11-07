# Learner Management API

This is a RESTful API for managing learner profiles, built with Go. It allows for creating, updating, retrieving, and deleting learner records, as well as generating a summary for each learner.

## Features

- **Create Learner**: Register a new learner with basic information.
- **Retrieve All Learners**: View a list of all registered learners.
- **Retrieve Learner by ID**: View detailed information for a specific learner by ID.
- **Update Learner**: Update details for an existing learner.
- **Delete Learner**: Remove a learner by their ID.
- **Generate Learner Summary**: Generate a brief profile summary for a learner by ID.

## Installation

1. **Clone the repository**:
    ```bash
    git clone <repository-url>
    cd <repository-folder>
    ```

2. **Run the application**:
    ```bash
    go run main.go
    ```
   The server will start on port `8080`.

## API Endpoints

### 1. Register a Learner

- **URL**: `/learners`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "name": "John Doe",
    "age": 25,
    "email": "john.doe@example.com"
  }
