#!/bin/sh
set -e

# Health check for the Go ecommerce API
if curl -f --silent --max-time 3 http://localhost:8080/health > /dev/null 2>&1; then
    exit 0
else
    exit 1
fi
