#!/bin/bash

set -e

echo "🚀 Rotki Demo - Quick Start Script"
echo "=================================="
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    echo "   Visit: https://docs.docker.com/get-docker/"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Check if config file exists
if [ ! -f "config.yaml" ]; then
    echo "📝 Creating config.yaml from template..."
    cp config.docker.yaml config.yaml
    echo ""
    echo "⚠️  IMPORTANT: Please edit config.yaml and add your DeBank API key!"
    echo "   Get your API key from: https://docs.cloud.debank.com"
    echo ""
    read -p "Press Enter to continue after you've added your API key..."
fi

echo ""
echo "1️⃣  Starting MySQL and Redis..."
docker-compose up -d mysql redis

echo ""
echo "⏳ Waiting for MySQL to be ready..."
for i in {1..30}; do
    if docker-compose exec -T mysql mysqladmin ping -h localhost -uroot -protki123 &> /dev/null; then
        echo "✅ MySQL is ready!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "❌ MySQL failed to start. Check logs with: docker-compose logs mysql"
        exit 1
    fi
    echo -n "."
    sleep 2
done

echo ""
echo "⏳ Waiting for Redis to be ready..."
for i in {1..10}; do
    if docker-compose exec -T redis redis-cli ping &> /dev/null; then
        echo "✅ Redis is ready!"
        break
    fi
    if [ $i -eq 10 ]; then
        echo "❌ Redis failed to start. Check logs with: docker-compose logs redis"
        exit 1
    fi
    echo -n "."
    sleep 1
done

echo ""
echo "2️⃣  Starting backend server..."
docker-compose up -d backend

echo ""
echo "⏳ Waiting for backend to be ready..."
for i in {1..20}; do
    if curl -s http://localhost:8080/health > /dev/null; then
        echo "✅ Backend is ready!"
        break
    fi
    if [ $i -eq 20 ]; then
        echo "⚠️  Backend might need more time. Check logs with: docker-compose logs backend"
    fi
    echo -n "."
    sleep 2
done

echo ""
echo "3️⃣  Setting up frontend..."
if [ ! -d "frontend/node_modules" ]; then
    echo "📦 Installing frontend dependencies..."
    cd frontend
    npm install
    cd ..
else
    echo "✅ Frontend dependencies already installed"
fi

echo ""
echo "✅ All services are running!"
echo ""
echo "📊 Service Status:"
docker-compose ps

echo ""
echo "🌐 Access points:"
echo "   - Backend API:  http://localhost:8080"
echo "   - Health check: http://localhost:8080/health"
echo ""
echo "🎯 Next steps:"
echo "   1. Start the frontend:"
echo "      cd frontend && npm run dev"
echo ""
echo "   2. Open your browser:"
echo "      http://localhost:3000"
echo ""
echo "📝 Useful commands:"
echo "   - View logs:        docker-compose logs -f"
echo "   - Stop services:    docker-compose down"
echo "   - Restart:          docker-compose restart"
echo "   - Check status:     docker-compose ps"
echo ""
echo "📚 Documentation:"
echo "   - Quick setup:      docs/SETUP.md"
echo "   - Docker guide:     docs/DOCKER.md"
echo "   - Architecture:     docs/ARCHITECTURE.md"
echo ""
echo "Happy tracking! 🎉"
