#!/bin/bash

# Docker Start Script
# استخدام: ./docker-start.sh أو bash docker-start.sh

set -e

echo "🚀 Starting MDM Project with Docker..."
echo ""

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

echo -e "${BLUE}✓ Docker and Docker Compose found${NC}"
echo ""

# Build and start containers (db, frontend, and backend-go only)
echo -e "${YELLOW}Building and starting containers (db, frontend, backend-go)...${NC}"
docker-compose up --build db frontend backend-go

# Print access information when up
echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✓ Project Started Successfully!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "📱 Frontend:  ${BLUE}http://localhost:3000${NC}"
echo -e "🔌 Backend:   ${BLUE}http://localhost:8081${NC}"
echo -e "💾 Database:  ${BLUE}localhost:5433${NC}"
echo -e "🌐 WebSocket: ${BLUE}ws://localhost:8081/rest/ws/connect${NC}"
echo ""
echo -e "${YELLOW}To stop: Press Ctrl+C${NC}"
echo -e "${YELLOW}To view logs: docker-compose logs -f${NC}"
echo ""
