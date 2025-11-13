# Gin Blog Backend

This project is a simple blog backend built with Gin and GORM in Go. It provides RESTful APIs for user authentication, post management, and comment management.

## Running Environment

*   **Go Version**: 1.25.0 or higher
*   **Database**: MySQL

## Installation

1.  **Clone the repository**:
    ```bash
    git clone <repository_url>
    cd web3-backend-learning-blog/docs/backend/gin/blog
    ```
    (Note: Replace `<repository_url>` with the actual repository URL if this project is part of a larger repository.)

2.  **Install dependencies**:
    ```bash
    go mod tidy
    ```

3.  **Database Setup**:
    Ensure you have a MySQL server running. The application will attempt to connect to a MySQL database with the following credentials (defined in `config/global.go`):
    *   **Username**: `testuser`
    *   **Password**: `123456test`
    *   **Host**: `127.0.0.1`
    *   **Port**: `3306`
    *   **Database Name**: `ginTestDB`

    The application will automatically create the `ginTestDB` database and migrate the necessary tables (`users`, `posts`, `comments`) on startup if they don't exist.

## Usage

1.  **Run the application**:
    ```bash
    go run main.go
    ```
    The server will start on `http://localhost:8080`.

2.  **API Endpoints**:

    ### Authentication
    *   **Register**: `POST /register`
        *   Request Body: `{ "username": "your_username", "password": "your_password" }`
    *   **Login**: `POST /login`
        *   Request Body: `{ "username": "your_username", "password": "your_password" }`
        *   Response: Returns a JWT token for authenticated requests.

    ### Posts (Requires Authentication)
    *   **Create Post**: `POST /posts`
        *   Headers: `Authorization: Bearer <your_jwt_token>`
        *   Request Body: `{ "title": "Your Post Title", "content": "Your post content" }`
    *   **Get All Posts**: `GET /posts`
        *   Headers: `Authorization: Bearer <your_jwt_token>`
    *   **Get Post by ID**: `GET /posts/:post_id`
        *   Headers: `Authorization: Bearer <your_jwt_token>`
    *   **Get Posts by User**: `GET /users/:user_id/posts`
        *   Headers: `Authorization: Bearer <your_jwt_token>`
    *   **Update Post**: `PUT /posts/:post_id`
        *   Headers: `Authorization: Bearer <your_jwt_token>`
        *   Request Body: `{ "title": "Updated Title", "content": "Updated content" }`
    *   **Delete Post**: `DELETE /posts/:post_id`
        *   Headers: `Authorization: Bearer <your_jwt_token>`

    ### Comments (Requires Authentication)
    *   **Create Comment**: `POST /comments`
        *   Headers: `Authorization: Bearer <your_jwt_token>`
        *   Request Body: `{ "post_id": 1, "content": "Your comment content" }`
    *   **Get Comments by Post**: `GET /posts/:post_id/comments`
        *   Headers: `Authorization: Bearer <your_jwt_token>`
    *   **Delete Comment**: `DELETE /comments/:comment_id`
        *   Headers: `Authorization: Bearer <your_jwt_token>`

## Logging

Application logs are output to the console and also saved to `logs/app.log`. The log file is automatically rotated.
