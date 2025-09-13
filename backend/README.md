# Backend Development Setup
This project provides a Go backend with MySQL, MinIO, and phpMyAdmin for development.

## Prerequisites
Before starting, make sure you have the following installed:
- [Go](https://golang.org/dl/) (>= 1.22 recommended)
- [Docker](https://docs.docker.com/get-docker/)

## Setup Instructions

1. **Clone Repository**
```bash
git clone https://github.com/nsza521/softdevKMITL.git
```

2. **CD Backend Directory**
```bash 
cd backend
```

3. **Install Go dependencies**
```bash
go mod download
``` 

4. **Build and Run All Services**
```bash
docker-compose -f docker-compose.dev.yml up --build
```

## Stop and remove containers, networks, and volumes
Use this command when you want to completely reset all services including databases.
```bash
docker-compose -f docker-compose.dev.yml down -v
```