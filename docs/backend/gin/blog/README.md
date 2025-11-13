# Gin 博客后端

该项目是一个使用 Go 语言的 Gin 和 GORM 构建的简单博客后端。它为用户认证、文章管理和评论管理提供了 RESTful API。

## 运行环境

*   **Go 版本**: 1.25.0 或更高
*   **数据库**: MySQL

## 安装

1.  **克隆仓库**:
    ```bash
    git clone <repository_url>
    cd web3-backend-learning-blog/docs/backend/gin/blog
    ```
    （注意：如果该项目是更大仓库的一部分，请将 `<repository_url>` 替换为实际的仓库 URL。）

2.  **安装依赖**:
    ```bash
    go mod tidy
    ```

3.  **数据库设置**:
    确保您有一个正在运行的 MySQL 服务器。应用程序将尝试使用以下凭据连接到 MySQL 数据库（在 `config/global.go` 中定义）：
    *   **用户名**: `testuser`
    *   **密码**: `123456test`
    *   **主机**: `127.0.0.1`
    *   **端口**: `3306`
    *   **数据库名**: `ginTestDB`

    如果 `ginTestDB` 数据库和必要的表（`users`、`posts`、`comments`）不存在，应用程序将在启动时自动创建它们。

## 使用

1.  **运行应用程序**:
    ```bash
    go run main.go
    ```
    服务器将在 `http://localhost:8080` 上启动。

2.  **API 端点**:

    ### 认证
    *   **注册**: `POST /register`
        *   请求体: `{ "username": "your_username", "password": "your_password" }`
    *   **登录**: `POST /login`
        *   请求体: `{ "username": "your_username", "password": "your_password" }`
        *   响应: 返回用于认证请求的 JWT 令牌。

    ### 文章 (需要认证)
    *   **创建文章**: `POST /posts`
        *   请求头: `Authorization: Bearer <your_jwt_token>`
        *   请求体: `{ "title": "Your Post Title", "content": "Your post content" }`
    *   **获取所有文章**: `GET /posts`
        *   请求头: `Authorization: Bearer <your_jwt_token>`
    *   **根据 ID 获取文章**: `GET /posts/:post_id`
        *   请求头: `Authorization: Bearer <your_jwt_token>`
    *   **根据用户获取文章**: `GET /users/:user_id/posts`
        *   请求头: `Authorization: Bearer <your_jwt_token>`
    *   **更新文章**: `PUT /posts/:post_id`
        *   请求头: `Authorization: Bearer <your_jwt_token>`
        *   请求体: `{ "title": "Updated Title", "content": "Updated content" }`
    *   **删除文章**: `DELETE /posts/:post_id`
        *   请求头: `Authorization: Bearer <your_jwt_token>`

    ### 评论 (需要认证)
    *   **创建评论**: `POST /comments`
        *   请求头: `Authorization: Bearer <your_jwt_token>`
        *   请求体: `{ "post_id": 1, "content": "Your comment content" }`
    *   **根据文章获取评论**: `GET /posts/:post_id/comments`
        *   请求头: `Authorization: Bearer <your_jwt_token>`
    *   **删除评论**: `DELETE /comments/:comment_id`
        *   请求头: `Authorization: Bearer <your_jwt_token>`

## 日志

应用程序日志会输出到控制台，并保存到 `logs/app.log` 文件中。日志文件会自动轮转。
