# Distributed Job Platform

A complete production-grade distributed job processing platform similar to Temporal, Celery, or Sidekiq. Built with Go (Gin), Next.js, PostgreSQL, and Redis.

## Architecture

![Architecture Diagram](architecture.md)

```mermaid
graph TD
    Client((Client)) -->|HTTP Rest| API[API Service Gin]
    Dashboard[Dashboard Next.js] -->|HTTP Rest| API
    
    API -->|CRUD/Query| DB[(PostgreSQL)]
    API -->|Enqueue Job| Queue[(Redis Queue)]
    
    Scheduler[Job Scheduler Cron] -->|Enqueue Scheduled Jobs| Queue
    Scheduler -->|Sync Logs| DB
    
    Queue -->|Dequeue| Worker1[Worker Node 1]
    Queue -->|Dequeue| Worker2[Worker Node 2]
    
    Worker1 -->|Update Status/Logs| DB
    Worker2 -->|Update Status/Logs| DB
    
    Worker1 -.->|Failed Jobs| DLQ[(Dead Letter Queue Redis)]
    Worker2 -.->|Failed Jobs| DLQ
    
    API -.->|Metrics| Prometheus[Prometheus]
    Worker1 -.->|Metrics/Health| Prometheus
    Worker2 -.->|Metrics/Health| Prometheus
    Scheduler -.->|Metrics| Prometheus
    
    Prometheus --> Grafana[Grafana Dashboard]
```

## Repository Structure

- `apps/api`: REST API for job submission and status tracking
- `apps/dashboard`: Next.js Web UI
- `services/worker`: Job processing engine
- `services/scheduler`: Cron job scheduling
- `packages/*`: Shared database, config, logger, and queue abstractions
- `infra/*`: Docker and observability configurations

## Tech Stack

- **Backend**: Go (Gin)
- **Frontend**: Next.js + TypeScript
- **Database**: PostgreSQL (GORM)
- **Message Broker/Queue**: Redis
- **Infra**: Docker Compose
- **Monitoring**: Prometheus + Grafana

## License
MIT
