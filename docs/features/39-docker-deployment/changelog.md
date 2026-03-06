# Feature #39 — Docker Deployment: Changelog

## [Unreleased]

### Added
- `Dockerfile` — multi-stage build (golang:alpine → alpine runtime).
- `docker-compose.yml` — app, postgres, redis services with healthchecks.
- `.dockerignore` — excludes unnecessary files from Docker build context.
