# Mintok

Mintok is an AI-powered developer productivity platform.

It helps developers:

- Analyze repositories
- Improve code quality
- Detect security issues
- Generate documentation
- Optimize performance
- Explain code using AI
- Review Pull Requests
- Generate tests
- Improve architecture

## Tech Stack

Backend

- Go
- Gin
- PostgreSQL
- Redis

Frontend

- React
- TypeScript
- Tailwind

Infrastructure

- Docker
- GitHub Actions

Status

Under Active Development

## Project Structure

- `backend` - Go Gin API service
- `frontend` - React, TypeScript, and Tailwind application
- `.github/workflows` - CI build, lint, and test automation

## Development

Backend:

```sh
cd backend
go test ./...
go build ./cmd/api
```

Frontend:

```sh
cd frontend
npm install
npm run format
npm run lint
npm test
npm run build
```

## Docker

Run the local stack:

```sh
docker compose up --build
```

Services:

- Frontend: `http://localhost:3000`
- Backend API: `http://localhost:8080`
- PostgreSQL: `localhost:5432`
- Redis: `localhost:6379`
