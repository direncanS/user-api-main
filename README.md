# User API

> Client-facing REST API for querying webcam image metadata and video sequences, with Redis caching for low-latency reads.

Part of the **Weather Archive** platform -- see also [webcam-api](https://github.com/direncanS/webcam-api-main) and [video-service](https://github.com/direncanS/video-service-main).

Built with **Go 1.19** | **AWS Lambda** | **PostgreSQL** | **Redis** | **S3**

---

## Overview

The User API is the read layer of the Weather Archive system. It exposes REST endpoints for browsing available webcam topics, querying image metadata with optional date-range filtering, and retrieving compiled video sequences. All database and S3 query results are cached in Redis with a 5-minute TTL to minimize latency on repeated requests.

## Features

- **Topic Discovery** -- list all available webcam topics
- **Date-Range Filtering** -- query images by topic with optional RFC 3339 start/end parameters
- **Video Retrieval** -- list all image objects within a video sequence folder with direct S3 URLs
- **Redis Caching** -- 5-minute TTL cache layer for all read operations
- **Dual Deployment** -- runs locally on port 8080 or as an AWS Lambda function (auto-detected)

## API Endpoints

| Method | Path | Query Params | Description |
|--------|------|-------------|-------------|
| GET | `/user/images` | -- | List all distinct webcam topics |
| GET | `/user/images/{topic}` | `start`, `end` (RFC 3339, optional) | Get image metadata for a topic |
| GET | `/user/videos` | `topic` (required) | Get all video sequences for a topic |
| GET | `/user/videos/{folder}` | -- | List all objects in a video folder with S3 URLs |

### Example Requests

```bash
# List all topics
curl http://localhost:8080/user/images

# Get images for a topic within a date range
curl "http://localhost:8080/user/images/vienna-skyline?start=2024-01-01T00:00:00Z&end=2024-03-01T00:00:00Z"

# Get video sequences
curl "http://localhost:8080/user/videos?topic=vienna-skyline"

# Get video objects with S3 URLs
curl http://localhost:8080/user/videos/550e8400-e29b-41d4-a716-446655440000
```

## Architecture

```
Client Application
        │
        ▼
   User API (Lambda / HTTP)
        │
        ├── Check Redis cache
        │   ├── Hit → return cached JSON
        │   └── Miss ↓
        │
        ├── Query PostgreSQL (topics, images, videos)
        │   or
        ├── List S3 objects (video folder contents)
        │
        ├── Cache result in Redis (5 min TTL)
        └── Return JSON response
```

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.19 |
| Runtime | AWS Lambda (ARM64) / standalone HTTP server |
| Router | Gorilla Mux with Lambda proxy adapter |
| Database | PostgreSQL (pgx v5) |
| Cache | Redis (go-redis v9) |
| Object Storage | AWS S3 (SDK v2) |
| Logging | Uber Zap |

## Response Format

All endpoints return a consistent JSON envelope:

```json
{
  "result": { },
  "error": ""
}
```

## Getting Started

### Prerequisites

- Go 1.19+
- PostgreSQL instance (shared with webcam-api and video-service)
- Redis instance
- AWS credentials with S3 read access

### Environment Variables

```bash
host=<postgres-host>
password=<postgres-password>
endpoint=<redis-endpoint>
REDIS_PASSWORD=<redis-password>
SECRET_KEY=<app-secret>
S3_BUCKET=<s3-bucket-name>
```

### Run Locally

```bash
make build && make run
# Server starts on :8080
```

### Deploy to AWS Lambda

```bash
make buildAWS   # ARM64 binary
make zip         # Creates bootstrap.zip for Lambda upload
```

## Related Services

| Service | Purpose |
|---------|---------|
| [webcam-api](https://github.com/direncanS/webcam-api-main) | Ingests webcam images via REST API into S3 and PostgreSQL |
| [video-service](https://github.com/direncanS/video-service-main) | Compiles stored webcam images into video sequences per topic |
