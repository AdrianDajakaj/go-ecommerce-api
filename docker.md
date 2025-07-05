# Go Echo E-Commerce API - Docker Setup

This directory contains the Docker configuration for the Go Echo e-commerce API backend.

## 🚀 Quick Start

### Build the image
```bash
docker build -t go-ecommerce-api .
```

### Run the container
```bash
docker run -d -p 8080:8080 --name api-backend go-ecommerce-api
```

### Test the deployment
```bash
./docker-test.sh
```

## 🔧 Configuration

### Environment Variables

- `JWT_SECRET` - JWT secret key for authentication
- `DB_PATH` - Database file path (default: `/app/data/ecommerce.db`)
- `ASSETS_PATH` - Static assets path (default: `/app/assets`)
- `PORT` - Server port (default: `8080`)

### Build Arguments

You can customize the build:

```bash
docker build --build-arg GO_VERSION=1.22 -t go-ecommerce-api .
```

## 📁 Files

- `Dockerfile` - Multi-stage build configuration
- `.dockerignore` - Files to exclude from Docker build context
- `docker-test.sh` - Local testing script
- `docker.md` - This documentation

## 🏗️ Build Process

1. **Build Stage**: Uses Go Alpine to compile the application with CGO support for SQLite
2. **Production Stage**: Uses Alpine Linux with minimal dependencies for security

## ✨ Features

- ✅ Multi-stage build for smaller image size (~20MB)
- ✅ CGO enabled for SQLite support
- ✅ Static binary compilation
- ✅ Non-root user for security
- ✅ Health check endpoint
- ✅ Proper file permissions
- ✅ Security hardening
- ✅ Dockerfile best practices

## 🔒 Security Features

- Non-root user execution
- Minimal base image (Alpine)
- Static binary (no dynamic dependencies)
- Read-only filesystem support
- Security scanning compatible

## 🐳 Docker Compose Usage

This Dockerfile is designed to work with the main project's `docker-compose.yml` located in the parent repository:

```yaml
# In your main repository's docker-compose.yml
services:
  api:
    build:
      context: ./go-ecommerce-api
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    volumes:
      - ./go-ecommerce-api/assets:/app/assets:ro
      - ./go-ecommerce-api:/app/data
    environment:
      - JWT_SECRET=your-jwt-secret-here
      - DB_PATH=/app/data/ecommerce.db
      - ASSETS_PATH=/app/assets
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

## 📊 Image Information

- **Base Images**: 
  - Build: `golang:1.22-alpine` 
  - Production: `alpine:latest`
- **Final Size**: ~20MB (approximate)
- **Exposed Port**: 8080
- **Health Check**: Enabled
- **User**: Non-root (appuser:1001)

## 🧪 Testing

The container includes a health check that verifies the service is responding:

```bash
# Check container health
docker ps

# Test health endpoint
curl http://localhost:8080/health

# Test API endpoints
curl http://localhost:8080/api/products
curl http://localhost:8080/api/categories
```

## 🛠️ Development

For local development with hot reload:

```bash
# Development mode (mount source code)
docker run -v $(pwd):/app -p 8080:8080 -w /app golang:1.22-alpine go run cmd/server.go
```

## 🔍 Troubleshooting

### Build Issues

1. **CGO errors**: Make sure SQLite development headers are installed
2. **Module download fails**: Check internet connectivity and Go proxy settings
3. **Binary won't run**: Verify CGO is enabled and target architecture matches

### Runtime Issues

1. **Database connection**: Check if database file is mounted correctly
2. **Assets not found**: Verify assets volume mount
3. **Permission denied**: Ensure volumes have correct permissions

### Database Issues

```bash
# Check database file permissions
ls -la /path/to/ecommerce.db

# Fix permissions if needed
chmod 644 /path/to/ecommerce.db
```

### Logs

```bash
# View container logs
docker logs <container-id>

# Follow logs in real-time
docker logs -f <container-id>

# Debug mode
docker run --rm -it go-ecommerce-api sh
```

## 🔧 Advanced Configuration

### Custom Database Path

```bash
docker run -e DB_PATH=/custom/path/db.sqlite go-ecommerce-api
```

### Read-only Root Filesystem

```bash
docker run --read-only --tmpfs /tmp go-ecommerce-api
```

### Resource Limits

```bash
docker run --memory=128m --cpus=0.5 go-ecommerce-api
```

## 🏷️ Tags and Versions

- `latest` - Latest stable version
- `v1.0.0` - Specific version tag
- `dev` - Development version

## 📈 Monitoring

The API exposes metrics and health endpoints:

- `/health` - Health check endpoint
- `/metrics` - Prometheus metrics (if implemented)
- `/api/status` - API status endpoint

## 🔄 Updates

To update the container:

```bash
# Pull latest changes
git pull

# Rebuild image
docker build -t go-ecommerce-api .

# Restart container
docker-compose restart api
```
