# Feature #39 — Docker Deployment: Design

## Dockerfile

Multi-stage build:
- **Build stage**: `golang:1.22-alpine`, copies source, builds static binary.
- **Runtime stage**: `alpine:3.19`, copies binary + resources + .env, exposes 8080, uses HEALTHCHECK.

## docker-compose.yml

Three services:
- **app** — built from Dockerfile, port 8080, depends on db health.
- **db** — `postgres:16-alpine` with healthcheck.
- **redis** — `redis:7-alpine`.

## .dockerignore

Excludes `.git`, `reference/`, `docs/`, test files, IDE configs.

## No Go Code Changes

Configuration files only.
