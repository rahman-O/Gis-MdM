# Docker Start Script for Windows PowerShell
# استخدام: powershell -ExecutionPolicy Bypass -File docker-start.ps1
# أو: .\docker-start.ps1

# Color functions
function Write-Success {
    param([string]$Message)
    Write-Host "✓ $Message" -ForegroundColor Green
}

function Write-Error-Custom {
    param([string]$Message)
    Write-Host "✗ $Message" -ForegroundColor Red
}

function Write-Info {
    param([string]$Message)
    Write-Host "ℹ $Message" -ForegroundColor Cyan
}

function Write-Warning-Custom {
    param([string]$Message)
    Write-Host "⚠ $Message" -ForegroundColor Yellow
}

# Header
Write-Host ""
Write-Host "╔════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║   MDM Project - Docker Start Script    ║" -ForegroundColor Cyan
Write-Host "║   تشغيل المشروع من Docker             ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# Check Docker installation
Write-Info "Checking Docker installation..."
$dockerCheck = docker --version 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Error-Custom "Docker is not installed or not in PATH"
    Write-Info "Please install Docker Desktop from: https://www.docker.com/products/docker-desktop"
    Read-Host "Press Enter to exit"
    exit 1
}
Write-Success "Docker found: $dockerCheck"

# Check Docker Compose installation
Write-Info "Checking Docker Compose installation..."
$composeCheck = docker-compose --version 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Error-Custom "Docker Compose is not available"
    Write-Info "Please ensure Docker Desktop is installed"
    Read-Host "Press Enter to exit"
    exit 1
}
Write-Success "Docker Compose found: $composeCheck"

# Check if Docker daemon is running
Write-Info "Checking if Docker is running..."
$dockerRunning = docker ps 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Error-Custom "Docker daemon is not running"
    Write-Info "Please start Docker Desktop"
    Read-Host "Press Enter to exit"
    exit 1
}
Write-Success "Docker daemon is running"

Write-Host ""
Write-Host "════════════════════════════════════════" -ForegroundColor Cyan
Write-Host ""

# Ask user if they want to rebuild
Write-Info "Would you like to rebuild images? (This may take a few minutes)"
$rebuild = Read-Host "Enter 'yes' to rebuild, or 'no' to use existing images (default: no)"

if ($rebuild.ToLower() -eq "yes") {
    Write-Warning-Custom "Building images with --no-cache..."
    Write-Host ""
    docker-compose build --no-cache
    if ($LASTEXITCODE -ne 0) {
        Write-Error-Custom "Build failed"
        Read-Host "Press Enter to exit"
        exit 1
    }
    Write-Success "Build completed successfully"
} else {
    Write-Info "Using existing images (or building if missing)"
}

Write-Host ""
Write-Host "════════════════════════════════════════" -ForegroundColor Cyan
Write-Host ""

# Start containers
Write-Info "Starting containers..."
Write-Host ""
docker-compose up

# If user stops the containers
Write-Host ""
Write-Host "════════════════════════════════════════" -ForegroundColor Cyan
Write-Host ""
Write-Info "Containers stopped"
Write-Host ""

# Ask if user wants to view logs
$viewLogs = Read-Host "View logs? (yes/no, default: no)"
if ($viewLogs.ToLower() -eq "yes") {
    docker-compose logs
}

# Ask if user wants to keep containers running
$keepRunning = Read-Host "Keep containers running in background? (yes/no, default: no)"
if ($keepRunning.ToLower() -eq "yes") {
    Write-Info "Starting containers in background..."
    docker-compose up -d
    Write-Success "Containers are running in background"
    Write-Info "To stop, run: docker-compose down"
} else {
    Write-Info "Stopping containers..."
    docker-compose down
    Write-Success "Containers stopped"
}

Write-Host ""
Write-Host "════════════════════════════════════════" -ForegroundColor Cyan
Write-Host ""
Write-Info "Access Points:"
Write-Host "  Frontend:  http://localhost:3000" -ForegroundColor Yellow
Write-Host "  Backend:   http://localhost:8081" -ForegroundColor Yellow
Write-Host "  WebSocket: ws://localhost:8081/rest/ws/connect" -ForegroundColor Yellow
Write-Host ""
Write-Host "════════════════════════════════════════" -ForegroundColor Cyan
Write-Host ""

Read-Host "Press Enter to exit"
