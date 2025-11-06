#!/bin/bash

echo "🛑 Stopping Rotki Demo services..."
echo ""

docker-compose down

echo ""
echo "✅ All services stopped"
echo ""
echo "To remove data volumes as well, run:"
echo "   docker-compose down -v"
