#!/bin/bash


set -e

echo "🐳 Building Go Echo API Docker Image..."

docker build -t go-ecommerce-api:latest .

echo "✅ Build completed successfully!"

echo "🧪 Testing the container..."

mkdir -p test_data

if [ -f "ecommerce.db" ]; then
    echo "📋 Making a copy of existing database for testing..."
    cp ecommerce.db test_data/ecommerce.db
else
    echo "⚠️  No existing database found, creating empty test database..."
    touch test_data/ecommerce.db
fi

CONTAINER_ID=$(docker run -d \
    -p 8080:8080 \
    -v $(pwd)/test_data:/app/data \
    -v $(pwd)/assets:/app/assets:ro \
    -e JWT_SECRET=test-secret-key \
    -e DB_PATH=/app/data/ecommerce.db \
    --name go-api-test \
    go-ecommerce-api:latest)

echo "⏳ Waiting for container to start..."
sleep 10

if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Container is responding on http://localhost:8080"
    echo "🌐 You can now test the API at: http://localhost:8080"
    echo ""
    echo "📝 Test endpoints:"
    echo "  Health: curl http://localhost:8080/health"
    echo "  Categories: curl http://localhost:8080/api/categories"
    echo "  Products: curl http://localhost:8080/api/products"
    echo ""
    echo "🔧 Commands:"
    echo "  View logs: docker logs $CONTAINER_ID"
    echo "  Stop container: docker stop $CONTAINER_ID"
    echo "  Remove container: docker rm $CONTAINER_ID"
    echo "  Remove image: docker rmi go-ecommerce-api:latest"
    echo "  Cleanup test data: rm -rf test_data"
else
    echo "❌ Container is not responding"
    echo "📋 Container logs:"
    docker logs $CONTAINER_ID
    docker stop $CONTAINER_ID
    docker rm $CONTAINER_ID
    rm -rf test_data
    exit 1
fi

echo ""
echo "🎉 API container is running successfully!"
echo "🔒 Your original database is SAFE - using test copy only!"
echo "Press Ctrl+C to stop watching logs, container will keep running"
echo ""

docker logs -f $CONTAINER_ID
