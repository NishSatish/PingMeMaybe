
# PingMeMaybe

## Description
PingMeMaybe is my experiment of building a high scale, low latency bulk notifications handling system.
Inspiration for the architecture is drawn from this [medium article](https://amanmadhukar.medium.com/building-a-scalable-cloud-native-notification-system-a4e64b42d671) by @madhukaraman.

## Features WIP

This is a monorepo based microservice system, which has a gateway that exposes client facing APIs. The gateway accepts requests and payloads for quick transactional notifications as well as for bulk broadcasting. There is a redis based queueing system built with Asynq, from which the "processor" service picks up "jobs" (sending notifications) and processes them. Current WIP focuses on building a cohort querying mechanism and computing/distributing the relevant batches to workers.

---

## Folder Structure

- `gateway/` — HTTP API server (Gin), routes, and notification queuing.
- `processor/` — Asynq worker, cron jobs, and notification processing.
- `libs/` — Shared code: config, DB models, DTOs, message patterns, utilities.
- `k8s/` — Kubernetes manifests for deployment.
- `sql/` — Contains migrations and seed scripts for the DB setup.

---

## Prerequisites

- Go 1.23+
- PostgreSQL 
- Redis (default port 6379)

---

## Setup

1. **Clone the repository**
   ```
   git clone <repo-url>
   cd PingMeMaybe
   ```

2. **Install dependencies**
   ```
   go mod tidy
   ```

3. **Create the database table**

   Feel free to use the scripts in the `sql/` directory to create the necessary tables. You can run the SQL scripts manually or use a migration tool like `golang-migrate`.


4. **Configure environment variables**
   Create a file named `local.env` in the project root:
   ```
   DATABASE_SESSION_POOLING_MODE_URL=
   REDIS_CLUSTER=
   REDIS_USERNAME=
   REDIS_PASSWORD=
   ```
5. **Make sure the redis and postgres servers are up**

---

## Running the Project

**Start the API Gateway:**
```
go run gateway/main.go
```

**Start the Processor:**
```
go run processor/main.go
```

---

## Usage

**Queue a notification:**
```
curl -X POST http://localhost:8080/notification \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Notification",
    "description": "This is a test notification.",
    "link": "https://example.com"
  }'
```

---

## Deployment

- Docker and Kubernetes manifests are available in the `k8s/` folder for containerized deployment.

---
