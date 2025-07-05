# 🚀 Deployment Instructions

## Main Repository Structure

```
main-project/
├── docker-compose.yml
├── go-ecommerce-api/          # API submodule
│   ├── assets/               # Local assets (not in git)
│   ├── ecommerce.db         # Local database (not in git)
│   ├── Dockerfile
│   └── ...
└── react-frontend/           # Frontend submodule
    ├── Dockerfile
    └── ...
```

## 📋 Deployment Steps

### 1. Main Repository Setup

```bash
# Create docker-compose.yml in the main repository
# Example configuration:

version: '3.8'

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
      - JWT_SECRET=your-jwt-secret
      - DB_PATH=/app/data/ecommerce.db
    networks:
      - ecommerce-network

  frontend:
    build:
      context: ./react-frontend
      dockerfile: Dockerfile
      args:
        - REACT_APP_API_URL=http://localhost:8080
    ports:
      - "3000:80"
    depends_on:
      api:
        condition: service_healthy
    networks:
      - ecommerce-network

networks:
  ecommerce-network:
    driver: bridge
```

### 2. Submodules Setup

```bash
# Make sure assets and database are locally available
cd go-ecommerce-api/
ls -la assets/      # Check if assets exist
ls -la ecommerce.db # Check if database exists
```

### 3. Building and Running

```bash
# From the main directory
docker-compose up --build
```

### 4. Testing

```bash
# Check if services are running
curl http://localhost:8080/health  # API health check
curl http://localhost:3000/health  # Frontend health check

# Check API endpoints
curl http://localhost:8080/api/categories
curl http://localhost:8080/api/products
```

## 🔧 Environment Configuration

### Development

```yaml
# docker-compose.override.yml (for development)
version: '3.8'

services:
  api:
    volumes:
      # Mount source code for hot reload (optional)
      - ./go-ecommerce-api:/app
    environment:
      - DEBUG=true
    command: go run cmd/server.go

  frontend:
    volumes:
      # Mount source code for hot reload (optional)
      - ./react-frontend/src:/app/src
```

### Production

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  api:
    environment:
      - JWT_SECRET=${JWT_SECRET}  # From .env
      - DB_PATH=/app/data/ecommerce.db
    deploy:
      resources:
        limits:
          memory: 256M
          cpus: '0.5'

  frontend:
    environment:
      - NODE_ENV=production
    deploy:
      resources:
        limits:
          memory: 128M
          cpus: '0.25'
```

## 🗂️ Volume Management

### Assets (static files)

```bash
# Assets are mounted as read-only
# To update assets:
docker-compose down
# Update files in go-ecommerce-api/assets/
docker-compose up
```

### Database

```bash
# Database backup
docker-compose exec api sqlite3 /app/data/ecommerce.db ".backup /tmp/backup.db"
docker cp $(docker-compose ps -q api):/tmp/backup.db ./backup.db

# Database restore
docker cp ./backup.db $(docker-compose ps -q api):/tmp/restore.db
docker-compose exec api sqlite3 /app/data/ecommerce.db ".restore /tmp/restore.db"
```

## 🛠️ Troubleshooting

### Permission issues

```bash
# Make sure Docker has access to files
chmod 644 go-ecommerce-api/ecommerce.db
chmod -R 755 go-ecommerce-api/assets/
```

### Port conflicts

```bash
# Change ports in docker-compose.yml if they are occupied
ports:
  - "8081:8080"  # API on port 8081
  - "3001:80"    # Frontend on port 3001
```

### Build issues

```bash
# Rebuild without cache
docker-compose build --no-cache

# Restart with full rebuild
docker-compose down -v
docker-compose up --build --force-recreate
```

## 📊 Monitoring

### Logs

```bash
# All logs
docker-compose logs -f

# API only
docker-compose logs -f api

# Frontend only
docker-compose logs -f frontend
```

### Health checks

```bash
# Container status
docker-compose ps

# Detailed health info
docker inspect $(docker-compose ps -q api) | grep -A5 Health
```

## 🔄 CI/CD Integration

### GitHub Actions example

```yaml
# .github/workflows/deploy.yml
name: Deploy
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          submodules: true
      
      - name: Build and deploy
        run: |
          docker-compose build
          docker-compose up -d
```

## 🔒 Security Notes

1. **Secrets**: Use `.env` files for sensitive data
2. **Networks**: Services communicate through private Docker network
3. **Volumes**: Assets are read-only, database has limited permissions
4. **Users**: Containers run as non-root users

## ⚡ Performance Tips

1. **Multi-stage builds**: Reduce image sizes
2. **Health checks**: Ensure reliability
3. **Resource limits**: Prevent system overload
4. **Layer caching**: Use .dockerignore for faster builds
