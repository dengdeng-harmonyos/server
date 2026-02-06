# Dengdeng Push Server

[![GitHub release](https://img.shields.io/github/v/release/dengdeng-harmonyos/server)](https://github.com/dengdeng-harmonyos/server/releases)
[![GitHub stars](https://img.shields.io/github/stars/dengdeng-harmonyos/server?style=social)](https://github.com/dengdeng-harmonyos/server)
[![License](https://img.shields.io/github/license/dengdeng-harmonyos/server)](LICENSE)
[![Docker Pulls](https://img.shields.io/docker/pulls/ricwang/dengdeng-server)](https://hub.docker.com/r/ricwang/dengdeng-server)
[![Go Version](https://img.shields.io/github/go-mod/go-version/dengdeng-harmonyos/server)](go.mod)

English | [简体中文](README.md)

## 📖 Introduction

Dengdeng Push Server is a **secure, privacy-friendly** push notification service solution designed specifically for **HarmonyOS Next**. This project is fully open-source and aims to provide developers with a trustworthy and easy-to-deploy push service infrastructure.

> 🎯 **v1.0 Official Release**: Production-ready with complete push functionality and automated deployment

### ✨ Highlights

- **🚀 One-Click Deployment**: Single Docker container with built-in PostgreSQL
- **🔐 Security First**: Configuration embedded at compile-time, supports AES-256-GCM encryption
- **📦 Zero Dependencies**: No external config files required, ready to use out of the box
- **🤖 CI/CD Automation**: GitHub Actions for automated build and deployment
- **🌍 Production Ready**: Supports separate test and production environment deployments

### 🔒 Security & Privacy Commitment

- **🚫 Zero Message Storage**: No push message content is stored, only anonymous statistics
- **🔐 End-to-End Encryption**: Push tokens stored with AES-256-GCM encryption
- **🎭 Anonymous Design**: Uses randomly generated Device Keys unlinked to real devices
- **📊 Anonymized Statistics**: Only success/failure counts recorded, no specific content
- **🔑 Compile-Time Config Embedding**: Sensitive configs embedded in binary during build
- **🛡️ Open Source Transparency**: All code is public and open to community review

## ✨ Core Features

### 🚀 Deployment & Operations

- **📦 Single Container Deployment**: PostgreSQL + Push service, ready out of the box
- **🔧 Embedded Configuration**: Huawei Push configs embedded at compile-time, no external files needed
- **🤖 Automated CI/CD**: GitHub Actions for automated build, test, and deployment
- **🏥 Health Checks**: Built-in health check endpoints for monitoring
- **🐳 Docker Support**: Official images hosted on Docker Hub
- **🔄 Auto-Restart**: Automatic recovery from container crashes

### 🔐 Security

- **🔒 AES-256-GCM Encryption**: Protects Push Token storage
- **🎲 Cryptographically Secure Randomness**: Uses crypto/rand for Device Key generation
- **🔑 RSA Public Key Support**: Optional end-to-end message encryption
- **⏱️ Auto-Expiration Mechanism**: Device Key validity management
- **🚦 Rate Limiting**: Prevents push abuse (daily limit per device)
- **🛡️ Compile-Time Key Injection**: Sensitive configs embedded via ldflags

### 🎯 Privacy Protection

- **📝 Zero Message Storage**: No push message content is saved
- **🎭 Fully Anonymous**: Device identifiers cannot be traced to real devices
- **📊 Aggregated Statistics**: Only statistical data recorded, cannot trace specific devices
- **🗑️ Auto-Cleanup**: Periodic cleanup of expired device records
- **🔍 Minimization Principle**: Database fields follow minimum necessary principle

### 📡 Functional Features

- **📬 Notification Push**: Supports notification bar messages (with title, content, custom data)
- **🃏 Card Refresh**: Supports HarmonyOS card updates
- **🔄 Background Push**: Supports background data push
- **📦 Batch Push**: Send messages to multiple devices at once
- **📊 Push Statistics**: View push success rate and historical data
- **🏥 Health Monitoring**: Built-in health check and service status endpoints
- **🌐 RESTful API**: Simple HTTP GET interface, easy to integrate

## 🚀 Quick Start

### Prerequisites

Before starting, you need:

1. **Huawei Developer Account**: [Huawei Developer Alliance](https://developer.huawei.com/)
2. **HarmonyOS Application**: Created HarmonyOS Next app
3. **Push Service Configuration**:
   - `agconnect-services.json` - Download from AppGallery Connect
   - `private.json` - Huawei Push Service account private key

### Method 1: Using Docker Hub Image (Recommended)

This is the simplest and fastest deployment method:

#### 1. Generate Encryption Key

```bash
# Generate 32-byte random key (Base64 encoded)
openssl rand -base64 32
```

Save the generated key to `.env` file:

```bash
echo "PUSH_TOKEN_ENCRYPTION_KEY=your-generated-key" > .env
```

#### 2. Start Service

```bash
# Pull latest image
docker pull ricwang/dengdeng-server:latest

# Start service
docker run -d \
  --name push-server \
  -p 8080:8080 \
  -e PUSH_TOKEN_ENCRYPTION_KEY=your-encryption-key \
  -e SERVER_NAME=Dengdeng\ Push\ Server \
  -v push-data:/var/lib/postgresql/data \
  --restart unless-stopped \
  ricwang/dengdeng-server:latest
```

> ⚠️ **Note**: Docker Hub image uses compile-time embedded Huawei Push configuration, suitable for public demos only. For production, use Method 2 to build your own.

#### 3. Verify Service

```bash
# Check health status
curl http://localhost:8080/health

# View logs
docker logs -f push-server
```

### Method 2: Build with Your Own Configuration (Production Recommended)

If deploying to production, it's recommended to use your own Huawei Push configuration:

#### 1. Prepare Configuration Files

Save the config files downloaded from Huawei Developer backend to GitHub Secrets:

- `AGCONNECT_JSON` - Complete content of `agconnect-services.json`
- `PRIVATE_JSON` - Complete content of `private.json`
- `PUSH_TOKEN_ENCRYPTION_KEY` - Generated using `openssl rand -base64 32`

#### 2. Fork Repository and Configure Secrets

1. Fork this repository to your GitHub account
2. Add the above Secrets in repository settings
3. Push code to `main` branch (test environment) or `release` branch (production environment)

#### 3. Automatic Build and Deployment

GitHub Actions will automatically:
- ✅ Embed your Huawei Push configuration at compile-time
- ✅ Build optimized statically-linked binary
- ✅ Build Docker image
- ✅ Deploy to your configured server

### Method 3: Local Development Build

```bash
# Clone repository
git clone https://github.com/dengdeng-harmonyos/server.git
cd server

# Prepare config files (place in project root)
# - agconnect-services.json
# - private.json

# Generate encryption key
echo "PUSH_TOKEN_ENCRYPTION_KEY=$(openssl rand -base64 32)" > .env

# Method A: Use Docker Compose
docker-compose up -d --build

# Method B: Local compile and run
go mod download
go build -o bin/server cmd/server/main.go

# Start database
docker run -d --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=push_server \
  -p 5432:5432 \
  postgres:15-alpine

# Run server
export AGCONNECT_SERVICES_FILE=agconnect-services.json
export PRIVATE_KEY_FILE=private.json
./bin/server
```

## 📡 API Interface

### Quick Overview

All endpoints use simple HTTP GET requests, no complex authentication required.

| Function | Endpoint | Description |
|---------|---------|-------------|
| Health Check | `GET /health` | Check service status |
| Device Registration | `GET /api/v1/device/register` | Register device and get Device Key |
| Notification Push | `GET /api/v1/push/notification` | Send notification bar message |
| Card Refresh | `GET /api/v1/push/form` | Update HarmonyOS card |
| Background Push | `GET /api/v1/push/background` | Send background data |
| Batch Push | `GET /api/v1/push/batch` | Batch send notifications |
| Push Statistics | `GET /api/v1/push/statistics` | View push statistics |

### Example: Send Notification

```bash
curl "http://your-server:8080/api/v1/push/notification?device_id=YOUR_DEVICE_KEY&title=Test%20Message&content=This%20is%20a%20test%20push"
```

### Example: Batch Push

```bash
curl "http://your-server:8080/api/v1/push/batch?device_ids=key1,key2,key3&title=Batch%20Notification&body=Send%20to%20multiple%20devices"
```

### Complete Documentation

For detailed API documentation and parameters:

- 📚 **API Documentation**: See API usage examples in the repository
- 🔍 **Source Reference**: [internal/handler](internal/handler) directory
- 💡 **Integration Examples**: Check HarmonyOS client project

## 🔧 Configuration

### Environment Variables

| Variable | Description | Required | Default |
|---------|-------------|:--------:|---------|
| `PUSH_TOKEN_ENCRYPTION_KEY` | Push Token encryption key (32 bytes, Base64) | ✅ | - |
| `SERVER_NAME` | Server identifier name | ❌ | `Dengdeng Push Server` |
| `PORT` | HTTP service port | ❌ | `8080` |
| `GIN_MODE` | Run mode (debug/release) | ❌ | `release` |
| `DEVICE_ID_TTL` | Device Key validity period (seconds) | ❌ | `31536000` (1 year) |
| `MAX_DAILY_PUSH_PER_DEVICE` | Daily push limit per device | ❌ | `1000` |

### GitHub Secrets Configuration (CI/CD)

If using GitHub Actions for automated builds, configure these Secrets:

| Secret Name | Description | How to Obtain |
|------------|-------------|---------------|
| `AGCONNECT_JSON` | AppGallery Connect config | Download `agconnect-services.json` from Huawei developer backend |
| `PRIVATE_JSON` | Push service account private key | Generate and download from Huawei developer backend |
| `PUSH_TOKEN_ENCRYPTION_KEY` | Encryption key | Generate with `openssl rand -base64 32` |
| `TEST_SERVER_HOST` | Test server address (optional) | For auto-deployment |
| `TEST_SERVER_USER` | Test server SSH user (optional) | For auto-deployment |
| `TEST_SERVER_SSH_KEY` | Test server SSH private key (optional) | For auto-deployment |
| `TEST_SERVER_PORT` | Test server SSH port (optional) | For auto-deployment |

### Getting Huawei Push Configuration

1. Visit [AppGallery Connect](https://developer.huawei.com/consumer/cn/service/josp/agc/index.html)
2. Select your app
3. Go to "Push Service" → "Configuration"
4. Download `agconnect-services.json`
5. Generate service account key (`private.json`) in "Project Settings" → "API Management"

### Data Persistence

Docker container uses named volumes to store PostgreSQL data:

```bash
# View data volumes
docker volume ls | grep push-data

# Backup data
docker run --rm -v push-data:/data -v $(pwd):/backup alpine \
  tar czf /backup/push-data-backup.tar.gz /data

# Restore data
docker run --rm -v push-data:/data -v $(pwd):/backup alpine \
  tar xzf /backup/push-data-backup.tar.gz -C /
```

## 📊 Data Storage

### Stored Data

1. **Device Information** (anonymized)
   - Device Key (randomly generated)
   - Push Token (AES-256-GCM encrypted)
   - Device metadata (type, version, etc.)
   - RSA public key (optional)

2. **Statistical Data** (aggregated)
   - Daily push count
   - Success/failure counts
   - Push type distribution

### Data Not Stored

- ❌ Push message content
- ❌ User identity information
- ❌ Device hardware identifiers
- ❌ Geolocation information
- ❌ IP addresses
- ❌ Any user-traceable information

## 🏗️ Architecture Design

### System Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Docker Container                      │
│  ┌─────────────────┐         ┌──────────────────┐      │
│  │  PostgreSQL 15  │ ←────→ │  Push Service (Go)│      │
│  │  - Device info  │         │  - Gin Framework  │      │
│  │  - Encrypted    │         │  - AES-256        │      │
│  │    Token        │         │  - Huawei Push    │      │
│  │  - Statistics   │         │    API            │      │
│  └─────────────────┘         └──────────────────┘      │
│         ↓                             ↑                 │
│    Persistent Volume                Port 8080           │
└───────────────────────────────────────┬─────────────────┘
                                        │
                                   HTTP API
                                        │
                    ┌───────────────────┼───────────────────┐
                    ↓                   ↓                   ↓
           HarmonyOS App 1      HarmonyOS App 2     Other Clients
```

### Data Flow

#### 1. Device Registration Flow

```
Client                 Push Service            Database
  │                      │                      │
  │─ Registration  ────→ │                      │
  │                      │                      │
  │                      │─ Generate Key  ────→ │
  │                      │  (crypto/rand)       │
  │                      │                      │
  │                      │─ Encrypt Token ────→ │
  │                      │  (AES-256-GCM)       │
  │                      │                      │
  │← Return Key  ───────  │                      │
```

#### 2. Push Message Flow

```
App Backend           Push Service           Huawei Push
  │                      │                      │
  │─ Push Request  ────→ │                      │
  │  (Device Key)        │                      │
  │                      │─ Decrypt Token ────→ │
  │                      │                      │
  │                      │                      │─ Send Push ──→ User Device
  │                      │← Result  ───────────  │
  │                      │                      │
  │← Success  ──────────  │                      │
  │                      │                      │
  │                      │─ Record Stats ─────→ Database
  │                      │  (no message content)
```

### Security Mechanisms

1. **Compile-Time Config Embedding**
   ```
   Source + Secrets → GitHub Actions
          ↓
   ldflags compile-time injection (Base64)
          ↓
   Statically-linked binary (no external deps)
          ↓
   Docker image (config embedded)
   ```

2. **Push Token Encrypted Storage**
   ```
   Plain Token → AES-256-GCM Encrypt → Database
   (Random Nonce)     (32-byte key)
   ```

3. **Device Key Generation**
   ```
   crypto/rand → Base64 URL Safe → Storage
   (32-byte random)  (no special chars)
   ```

## 🔐 Security Best Practices

### 1. Key Management

**Generate Strong Keys**
```bash
# Recommended: Use OpenSSL to generate 32-byte random key
openssl rand -base64 32

# Or use /dev/urandom (Linux/macOS)
head -c 32 /dev/urandom | base64
```

**Key Rotation**
```bash
# 1. Generate new key
NEW_KEY=$(openssl rand -base64 32)

# 2. Update GitHub Secrets or environment variables

# 3. Rebuild and redeploy
# GitHub Actions will automatically compile with new key

# 4. Old devices need to re-register
```

**Storage Security**
- ✅ Use environment variables or key management services
- ✅ Use GitHub Secrets for sensitive configs
- ✅ Compile-time embedding to avoid config file exposure
- ❌ Don't hardcode in source code
- ❌ Don't commit to Git repository
- ❌ Don't output through logs

### 2. Network Security

**Use HTTPS**
```nginx
# Nginx reverse proxy example
server {
    listen 443 ssl http2;
    server_name push.yourdomain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

**Firewall Rules**
```bash
# Allow only specific IPs (optional)
sudo ufw allow from YOUR_IP to any port 8080

# Or allow only internal network
sudo ufw allow from 10.0.0.0/8 to any port 8080
```

**Rate Limiting**
```bash
# Built-in application-level rate limiting:
# - Maximum 1000 pushes per device per day
# - Adjustable via MAX_DAILY_PUSH_PER_DEVICE environment variable
```

### 3. Database Security

**Regular Backups**
```bash
# Create backup
docker exec push-server pg_dump -U postgres push_server > backup-$(date +%Y%m%d).sql

# Automated backup script (add to crontab)
0 2 * * * docker exec push-server pg_dump -U postgres push_server | gzip > /backup/push-$(date +\%Y\%m\%d).sql.gz
```

**Clean Expired Data**
```sql
-- Clean expired devices from 30 days ago
DELETE FROM devices WHERE expired_at < NOW() - INTERVAL '30 days';

-- Clean push statistics from 90 days ago
DELETE FROM push_statistics WHERE date < NOW() - INTERVAL '90 days';
```

### 4. Monitoring & Auditing

**Health Monitoring**
```bash
# Basic health check
curl http://localhost:8080/health

# With monitoring system (e.g., Prometheus)
# Periodically check health status and alert
```

**Log Auditing**
```bash
# View push logs
docker logs push-server | grep "Push"

# View error logs
docker logs push-server | grep "ERROR"

# Real-time monitoring
docker logs -f push-server
```

**Anomaly Detection**
```bash
# Check abnormally high push frequency
# View push statistics API
curl "http://localhost:8080/api/v1/push/statistics?date=$(date +%Y-%m-%d)"
```

### 5. Deployment Security Checklist

Confirm before deployment:

- [ ] ✅ Generated strong random encryption key
- [ ] ✅ Configured HTTPS/TLS
- [ ] ✅ Set up firewall rules
- [ ] ✅ Configured data backup strategy
- [ ] ✅ Enabled health check monitoring
- [ ] ✅ Reviewed log output (no sensitive info)
- [ ] ✅ Limited server access permissions
- [ ] ✅ Updated all dependencies to latest versions
- [ ] ✅ Configured auto-restart policy
- [ ] ✅ Tested push functionality works normally

## 📦 Docker Images

### Official Images

🐳 **Docker Hub**: [ricwang/dengdeng-server](https://hub.docker.com/r/ricwang/dengdeng-server)

### Available Tags

| Tag | Description | Update Frequency |
|-----|-------------|------------------|
| `latest` | Latest stable version (main branch) | Every commit to main |
| `v1.0.0`, `v1.0.x` | Specific version number | Created on release |
| `release` | Production release version | Commits to release branch |

### Image Details

- **Base Image**: `postgres:15-alpine`
- **Included Components**: PostgreSQL 15 + Go Push Service
- **Image Size**: ~300MB
- **Supported Architecture**: `linux/amd64`
- **Configuration Method**: Compile-time embedding (Docker Hub image uses demo config)

### Image Build

All images are automatically built via GitHub Actions, ensuring:

- ✅ **Reproducible Builds**: Same code generates same image
- ✅ **Security Scanning**: No sensitive info leaked during build
- ✅ **Static Linking**: No external dependencies, runs directly
- ✅ **Minimized Size**: Uses Alpine base image

### Self-Built Images

```bash
# Clone repository
git clone https://github.com/dengdeng-harmonyos/server.git
cd server

# Prepare config files
# - agconnect-services.json
# - private.json

# Build image
docker build -t my-dengdeng-server .

# Run
docker run -d \
  --name push-server \
  -p 8080:8080 \
  -e PUSH_TOKEN_ENCRYPTION_KEY=$(openssl rand -base64 32) \
  -v push-data:/var/lib/postgresql/data \
  my-dengdeng-server
```

## 🛠️ Development Guide

### Requirements

- **Go**: 1.21 or higher
- **PostgreSQL**: 15 or higher
- **Docker**: 20.10 or higher (optional)
- **Git**: 2.x

### Local Development Setup

#### 1. Clone Code

```bash
git clone https://github.com/dengdeng-harmonyos/server.git
cd server
```

#### 2. Install Dependencies

```bash
go mod download
```

#### 3. Prepare Config Files

Place config files downloaded from Huawei developer backend in project root:
- `agconnect-services.json`
- `private.json`

#### 4. Start Database

```bash
# Start PostgreSQL using Docker
docker run -d \
  --name postgres-dev \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=push_server \
  -p 5432:5432 \
  postgres:15-alpine

# Wait for database to start
sleep 5

# Run database migrations
cd database
./migrate.sh
cd ..
```

#### 5. Run Development Server

```bash
# Set environment variables
export PUSH_TOKEN_ENCRYPTION_KEY=$(openssl rand -base64 32)
export GIN_MODE=debug
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=push_server

# Run server
go run cmd/server/main.go
```

Server will start at `http://localhost:8080`.

### Build

#### Local Build

```bash
# Compile binary
go build -o bin/server cmd/server/main.go

# Run
./bin/server
```

#### Build with Embedded Config

```bash
# Base64 encode config files
AGCONNECT_BASE64=$(cat agconnect-services.json | base64)
PRIVATE_BASE64=$(cat private.json | base64)
ENCRYPTION_KEY_BASE64=$(echo "your-encryption-key" | base64)

# Compile with injected config
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags "\
    -X 'github.com/dengdeng-harmonyos/server/internal/config.embeddedAgConnectJSON=$AGCONNECT_BASE64' \
    -X 'github.com/dengdeng-harmonyos/server/internal/config.embeddedPrivateJSON=$PRIVATE_BASE64' \
    -X 'github.com/dengdeng-harmonyos/server/internal/config.embeddedEncryptionKey=$ENCRYPTION_KEY_BASE64' \
    -s -w" \
  -o bin/dengdeng-server \
  cmd/server/main.go
```

#### Build Docker Image

```bash
# Local build
docker build -t dengdeng-server:dev .

# Use CI Dockerfile (requires pre-compiled binary)
docker build -f Dockerfile.ci -t dengdeng-server:ci .
```

### Database Management

#### Create Migration File

```bash
cd database
./create-migration.sh add_new_feature
```

This creates two files:
- `migrations/YYYYMMDDHHMMSS_add_new_feature.up.sql` - Forward migration
- `migrations/YYYYMMDDHHMMSS_add_new_feature.down.sql` - Rollback migration

#### Run Migrations

```bash
cd database
./migrate.sh
```

#### Rollback Migration

```bash
cd database
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/push_server?sslmode=disable" down 1
```

### Project Structure

```
server/
├── cmd/
│   └── server/
│       └── main.go              # Application entry
├── internal/
│   ├── config/
│   │   ├── config.go            # Config loading
│   │   └── embedded_secrets.go  # Embedded config
│   ├── database/
│   │   └── database.go          # Database operations
│   ├── handler/
│   │   ├── device.go            # Device management
│   │   ├── message.go           # Message processing
│   │   ├── push.go              # Push logic
│   │   └── response.go          # Response wrapper
│   ├── logger/
│   │   └── logger.go            # Logging system
│   ├── middleware/
│   │   └── middleware.go        # HTTP middleware
│   ├── models/
│   │   └── models.go            # Data models
│   └── service/
│       ├── crypto.go            # Encryption service
│       ├── encryption.go        # Token encryption
│       └── huawei_push.go       # Huawei Push API
├── database/
│   ├── migrations/              # Database migration files
│   ├── migrate.sh              # Migration script
│   └── 001_initial_schema.sql  # Initial DB structure
├── .github/
│   └── workflows/
│       └── build.yml           # CI/CD config
├── Dockerfile                  # Standard Dockerfile
├── Dockerfile.ci               # CI-specific Dockerfile
├── docker-compose.yml          # Docker Compose config
└── README.md
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./internal/service/...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Standards

```bash
# Format code
go fmt ./...

# Static analysis
go vet ./...

# Use golangci-lint (recommended)
golangci-lint run
```

## 🤝 Contributing

We welcome all forms of contributions! Whether reporting bugs, suggesting new features, or submitting code, everything helps make this project better.

### How to Contribute

1. **Fork the Repository**
   ```bash
   # Click Fork button on GitHub
   ```

2. **Clone Your Fork**
   ```bash
   git clone https://github.com/your-username/server.git
   cd server
   ```

3. **Create Feature Branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```

4. **Make Changes and Commit**
   ```bash
   git add .
   git commit -m "Add some amazing feature"
   ```

5. **Push to Your Fork**
   ```bash
   git push origin feature/amazing-feature
   ```

6. **Create Pull Request**
   - Open your Fork on GitHub
   - Click "New Pull Request"
   - Describe your changes

### Focus Areas for Contributions

We especially welcome contributions in:

- 🔒 **Security Improvements**: Encryption algorithm optimization, security vulnerability fixes
- 🔐 **Privacy Protection**: Better data anonymization solutions
- 📝 **Documentation**: API docs, tutorials, best practices
- 🐛 **Bug Fixes**: Finding and fixing issues
- ✨ **New Features**: New push types, management features, etc.
- 🧪 **Test Coverage**: Unit tests, integration tests
- 🌍 **Internationalization**: Multi-language support
- 🎨 **UI/UX**: Management interface improvements

### Commit Message Convention

Please follow this commit message format:

```
<type>: <short description>

<detailed description>

<related issue>
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation update
- `style`: Code formatting
- `refactor`: Code refactoring
- `test`: Test-related
- `chore`: Build/tools-related

**Example**:
```
feat: add batch push API

Implemented functionality to send pushes to multiple devices
simultaneously, supporting up to 100 devices in batch operation.

Closes #123
```

### Development Workflow

1. **Ensure code passes tests**
   ```bash
   go test ./...
   ```

2. **Format code**
   ```bash
   go fmt ./...
   go vet ./...
   ```

3. **Update documentation**
   - If adding new features, update README.md
   - If modifying API, update API docs

4. **Pre-commit checklist**
   - [ ] Code is formatted
   - [ ] Tests pass
   - [ ] Documentation updated
   - [ ] Commit message is clear

### Reporting Issues

Found a bug? Please [create an Issue](https://github.com/dengdeng-harmonyos/server/issues/new) and include:

- 🔍 **Problem Description**: Clearly describe the issue encountered
- 📋 **Reproduction Steps**: How to trigger the problem
- 💻 **Environment Info**: OS, Go version, Docker version, etc.
- 📸 **Screenshots/Logs**: If applicable

### Feature Suggestions

Have a new idea? Please [create a Feature Request](https://github.com/dengdeng-harmonyos/server/issues/new) and explain:

- 💡 **Feature Description**: What feature you want
- 🎯 **Use Case**: Why this feature is needed
- 📝 **Expected Behavior**: How the feature should work
- 🔄 **Alternatives**: Are there other solutions

### Code of Conduct

- ✅ Respect all contributors
- ✅ Stay friendly and professional
- ✅ Accept constructive criticism
- ✅ Focus on project's overall benefit
- ❌ No harassment or discriminatory language

### Getting Help

Need help? Get assistance through:

- 📖 Check [documentation](README.md)
- 💬 Ask in [Issues](https://github.com/dengdeng-harmonyos/server/issues)
- 🔍 Search existing Issues and Pull Requests

## 🎯 Roadmap

### v1.0 ✅ (Current Version)

- [x] Basic push functionality (notification, card, background)
- [x] Device registration and management
- [x] Push Token AES-256-GCM encryption
- [x] Batch push support
- [x] Docker single container deployment
- [x] Push statistics functionality
- [x] GitHub Actions CI/CD
- [x] Compile-time config embedding
- [x] Rate limiting and security protection

### v1.1 🚧 (Planned)

- [ ] Web management interface
  - Device management
  - Push history view
  - Real-time statistics charts
- [ ] API key management system
- [ ] More detailed push reports
- [ ] Push template management
- [ ] Scheduled push functionality

### v2.0 🔮 (Future)

- [ ] Multi-provider support
  - Firebase Cloud Messaging
  - Apple Push Notification
  - Other cloud push services
- [ ] Message priority queue
- [ ] Push retry mechanism optimization
- [ ] Multi-language SDKs
  - Go SDK
  - JavaScript/TypeScript SDK
  - Python SDK
- [ ] Webhook support
- [ ] Push rules engine

### Community Suggestions

Welcome to share your ideas in [Issues](https://github.com/dengdeng-harmonyos/server/issues)!

---

## 📄 License

This project is licensed under the **MIT License**. See [LICENSE](LICENSE) file for details.

### License Explanation

- ✅ Commercial use allowed
- ✅ Modification allowed
- ✅ Distribution allowed
- ✅ Private use allowed
- ⚠️ License and copyright notice must be included
- ⚠️ No liability warranty provided

---

## 🌟 Acknowledgments

Thanks to all contributors to this project!

### Technical Support

- [HarmonyOS Next](https://developer.harmonyos.com/) - HarmonyOS Operating System
- [Huawei Push Kit](https://developer.huawei.com/consumer/cn/hms/huawei-pushkit/) - Huawei Push Service
- [Gin Web Framework](https://gin-gonic.com/) - Go Web Framework
- [PostgreSQL](https://www.postgresql.org/) - Open Source Database

### Contributors

<a href="https://github.com/dengdeng-harmonyos/server/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=dengdeng-harmonyos/server" />
</a>

---

## 📞 Contact & Support

### Project Links

- 🏠 **Project Home**: [https://github.com/dengdeng-harmonyos/server](https://github.com/dengdeng-harmonyos/server)
- 🐛 **Issue Tracker**: [GitHub Issues](https://github.com/dengdeng-harmonyos/server/issues)
- 🐳 **Docker Image**: [Docker Hub](https://hub.docker.com/r/ricwang/dengdeng-server)
- 📖 **Documentation**: [README](README.md) | [简体中文](README_EN.md)

### Get Help

- 💬 Ask in [GitHub Issues](https://github.com/dengdeng-harmonyos/server/issues)
- 📧 Email project maintainers
- ⭐ Star the project to follow latest updates

---

## 📊 Project Status

### Statistics

[![GitHub stars](https://img.shields.io/github/stars/dengdeng-harmonyos/server?style=social)](https://github.com/dengdeng-harmonyos/server)
[![GitHub forks](https://img.shields.io/github/forks/dengdeng-harmonyos/server?style=social)](https://github.com/dengdeng-harmonyos/server/fork)
[![GitHub issues](https://img.shields.io/github/issues/dengdeng-harmonyos/server)](https://github.com/dengdeng-harmonyos/server/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/dengdeng-harmonyos/server)](https://github.com/dengdeng-harmonyos/server/pulls)
[![GitHub license](https://img.shields.io/github/license/dengdeng-harmonyos/server)](LICENSE)

### Star History

[![Star History Chart](https://api.star-history.com/svg?repos=dengdeng-harmonyos/server&type=Date)](https://star-history.com/#dengdeng-harmonyos/server&Date)

---

## ⚠️ Disclaimer

**This service provides push infrastructure and does not store any user data.**

- 🔒 Ensure your encryption key is secure, don't share with others
- 🔐 Keep Huawei Push Service config files safe
- 📝 Comply with local laws and privacy protection policies
- ⚖️ This project is not responsible for any consequences of using this service
- 🛡️ Regularly update dependencies and security patches

---

## 💡 Final Words

If this project helps you, we welcome you to:

- ⭐ Give the project a Star
- 🔄 Fork and contribute
- 📢 Share with more developers
- 💬 Provide feedback and suggestions

**Let's build a secure and reliable push service together!** 🚀
