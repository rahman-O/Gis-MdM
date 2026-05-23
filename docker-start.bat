@echo off
REM Docker Start Script for Windows
REM استخدام: docker-start.bat أو double-click على الملف

color 0A
echo.
echo =========================================
echo   Docker Start Script
echo   تشغيل المشروع
echo =========================================
echo.

REM Check if Docker is installed
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Docker is not installed or not in PATH
    echo Please install Docker Desktop from https://www.docker.com/products/docker-desktop
    pause
    exit /b 1
)

REM Check if Docker Compose is available
docker-compose --version >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Docker Compose is not available
    echo Please ensure Docker Desktop is installed
    pause
    exit /b 1
)

echo [OK] Docker and Docker Compose found
echo.

echo Starting containers...
echo.

REM Start Docker Compose
docker-compose up --build

pause
