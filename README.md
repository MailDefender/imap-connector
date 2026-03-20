# **imap-connector**

A Go-based IMAP connector providing a RESTful API for mailbox management, message retrieval, and email operations.

> Note: This documentation is currently being drafted and will be completed in a future version.

## 📋 Table of Contents

- [Features](#-features)
- [Prerequisites](#-prerequisites)
- [Installation](#-installation)
- [Usage](#-usage)
- [API Documentation](#-api-documentation)
- [Configuration](#-configuration)
- [Testing](#-testing)
- [Contributing](#-contributing)
- [License](#-license)

## ✨ Features

- **Mailbox Management**: List mailboxes, get status, create new mailboxes.
- **Message Operations**: Fetch messages with pagination, move messages between mailboxes.
- **Header Parsing**: Robust parsing of email headers.
- **Hashing**: Generate unique hashes for messages without Message-ID which allow us the identify them.
- **Thread-Safe**: Mutex-based locking for concurrent IMAP operations.
- **Docker Support**: Easy deployment with Docker.

## 📦 Prerequisites

- **Go**: Version 1.18 or later.
- **Docker**: Optional, for containerized deployment.
- **IMAP Server**: Access to an IMAP server (e.g., Gmail, Dovecot).

## 🚀 Installation

### With Docker (Recommended)

1. Build the image:

   ```bash
   docker build -t maildefender/imap-connector .
   ```

2. Run the container:
   ```bash
   docker run -p 8080:8080 --env-file .env maildefender/imap-connector
   ```

### Without Docker

1. Clone the repository:

   ```bash
   git clone https://github.com/MailDefender/imap-connector.git
   cd imap-connector
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Build the application:

   ```bash
   go build -o imap-connector ./cmd/imap-connector
   ```

4. Set environment variables (see [Configuration](#-configuration)):

   ```bash
   source .env
   ```

5. Run the application:
   ```bash
   ./imap-connector
   ```

## 🏃‍♂️ Usage

The application starts an HTTP server on port 8080 (configurable via environment). It provides REST endpoints for IMAP operations.

### Example API Calls

- List mailboxes:

  ```bash
  curl -X GET "http://localhost:8080/v1/mailbox"
  ```

- Fetch messages:

  ```bash
  curl -X GET "http://localhost:8080/v1/mailbox/INBOX/message?limit=10&offset=0"
  ```

- Move a message:
  ```bash
  curl -X POST "http://localhost:8080/v1/mailbox/INBOX/message/123/move" \
       -H "Content-Type: application/json" \
       -d '{"destination": "ARCHIVE"}'
  ```

For full API details, see [API Documentation](#-api-documentation).

## 📖 API Documentation

The API uses the [Shaker](https://github.com/kessaro/shaker) framework. Swagger documentation is available at `/swagger` endpoint once the server is running.

Key endpoints:

- `GET /v1/mailbox` - List all mailboxes
- `GET /v1/mailbox/status` - Get mailbox status
- `GET /v1/mailbox/:mailbox/message` - Fetch messages from a mailbox
- `POST /v1/mailbox/:mailbox/message/:messageId/move` - Move a message
- `POST /v1/mailbox` - Create a new mailbox

## 🛠 Configuration

Create a `.env` file in the project root with the following variables:

```env
IMAP_HOST=imap.example.com
IMAP_PORT=993
IMAP_TLS=true
IMAP_USERNAME=your-email@example.com
IMAP_PASSWORD=your-password
```

- `IMAP_HOST`: IMAP server hostname.
- `IMAP_PORT`: IMAP server port (993 for TLS, 143 for non-TLS).
- `IMAP_TLS`: Enable TLS connection (default: true).
- `IMAP_USERNAME`: IMAP username.
- `IMAP_PASSWORD`: IMAP password.

**Security Note**: Never commit `.env` files to version control. Use environment-specific secrets management.

## 🧪 Testing

### Unit Tests

Run unit tests with:

```bash
go test ./...
```

### Integration Testing with Mock Server

The project includes a mock IMAP server using Dovecot for testing:

1. Start the mock server:

   ```bash
   cd mock
   docker-compose up -d
   ```

2. Run tests or the application against the mock.

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository.
2. Create a feature branch: `git checkout -b feature/your-feature`.
3. Make your changes and add tests.
4. Run tests: `go test ./...`.
5. Commit your changes: `git commit -am 'Add some feature'`.
6. Push to the branch: `git push origin feature/your-feature`.
7. Submit a pull request.

### Code Style

- Follow Go conventions (use `gofmt`, `go vet`).
- Add unit tests for new features.
- Update documentation as needed.

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you encounter issues or have questions:

- Open an issue on GitHub.
- Check the [Troubleshooting](#-troubleshooting) section.

## 🔧 Troubleshooting

- **Connection Issues**: Verify IMAP credentials and server settings.
- **Build Errors**: Ensure Go 1.18+ is installed.
- **Test Failures**: Check mock server is running for integration tests.

## 📈 Roadmap

- [ ] Complete API documentation.
- [ ] Add authentication middleware.
- [ ] Support for IMAP search queries.
- [ ] Webhook notifications for new messages.
