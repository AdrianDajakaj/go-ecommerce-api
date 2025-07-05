# 🚀 Deployment Instructions

## Struktura głównego repozytorium

```
main-project/
├── docker-compose.yml
├── go-ecommerce-api/          # Submoduł API
│   ├── assets/               # Lokalne assety (nie w git)
│   ├── ecommerce.db         # Lokalna baza danych (nie w git)
│   ├── Dockerfile
│   └── ...
└── react-frontend/           # Submoduł frontendu
    ├── Dockerfile
    └── ...
```

## 📋 Kroki wdrożenia

### 1. Przygotowanie głównego repozytorium

```bash
# W głównym repozytorium stwórz docker-compose.yml
# Przykład konfiguracji:

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

### 2. Przygotowanie submodułów

```bash
# Upewnij się, że assets i baza danych są lokalnie dostępne
cd go-ecommerce-api/
ls -la assets/      # Sprawdź czy assety istnieją
ls -la ecommerce.db # Sprawdź czy baza danych istnieje
```

### 3. Buildowanie i uruchamianie

```bash
# Z głównego katalogu
docker-compose up --build
```

### 4. Testowanie

```bash
# Sprawdź czy usługi działają
curl http://localhost:8080/health  # API health check
curl http://localhost:3000/health  # Frontend health check

# Sprawdź API endpoints
curl http://localhost:8080/api/categories
curl http://localhost:8080/api/products
```

## 🔧 Konfiguracja środowisk

### Development

```yaml
# docker-compose.override.yml (dla development)
version: '3.8'

services:
  api:
    volumes:
      # Mount source code dla hot reload (opcjonalnie)
      - ./go-ecommerce-api:/app
    environment:
      - DEBUG=true
    command: go run cmd/server.go

  frontend:
    volumes:
      # Mount source code dla hot reload (opcjonalnie)
      - ./react-frontend/src:/app/src
```

### Production

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  api:
    environment:
      - JWT_SECRET=${JWT_SECRET}  # Z .env
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

## 🗂️ Zarządzanie volumes

### Assets (statyczne pliki)

```bash
# Assets są montowane jako read-only
# Aby zaktualizować assets:
docker-compose down
# Zaktualizuj pliki w go-ecommerce-api/assets/
docker-compose up
```

### Baza danych

```bash
# Backup bazy danych
docker-compose exec api sqlite3 /app/data/ecommerce.db ".backup /tmp/backup.db"
docker cp $(docker-compose ps -q api):/tmp/backup.db ./backup.db

# Restore bazy danych
docker cp ./backup.db $(docker-compose ps -q api):/tmp/restore.db
docker-compose exec api sqlite3 /app/data/ecommerce.db ".restore /tmp/restore.db"
```

## 🛠️ Rozwiązywanie problemów

### Permission issues

```bash
# Upewnij się, że Docker ma dostęp do plików
chmod 644 go-ecommerce-api/ecommerce.db
chmod -R 755 go-ecommerce-api/assets/
```

### Port conflicts

```bash
# Zmień porty w docker-compose.yml jeśli są zajęte
ports:
  - "8081:8080"  # API na porcie 8081
  - "3001:80"    # Frontend na porcie 3001
```

### Build issues

```bash
# Rebuild bez cache
docker-compose build --no-cache

# Restart z pełnym rebuild
docker-compose down -v
docker-compose up --build --force-recreate
```

## 📊 Monitoring

### Logi

```bash
# Wszystkie logi
docker-compose logs -f

# Tylko API
docker-compose logs -f api

# Tylko frontend
docker-compose logs -f frontend
```

### Health checks

```bash
# Status kontenerów
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

1. **Secrets**: Używaj `.env` plików dla wrażliwych danych
2. **Networks**: Usługi komunikują się przez prywatną sieć Docker
3. **Volumes**: Assets są read-only, baza danych ma ograniczone uprawnienia
4. **Users**: Kontenery działają jako non-root users

## ⚡ Performance Tips

1. **Multi-stage builds**: Zmniejszają rozmiar obrazów
2. **Health checks**: Zapewniają niezawodność
3. **Resource limits**: Zapobiegają przeciążeniu systemu
4. **Layer caching**: Wykorzystuj .dockerignore dla szybszych buildów
