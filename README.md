# Enterprise Monitoring & Alerting

## English | [中文](#中文)

A distributed tracing, log aggregation, and metrics alerting platform built with Go, Gin, PostgreSQL, and Redis.

### Features
- **Distributed tracing** OpenTelemetry-compatible span ingest and trace reconstruction
- **Log aggregation** structured log ingest, full-text and field-level search
- **Metrics** counter/gauge/histogram with multi-dimensional labels
- **Alerting** rule-based alerts with severity, routing, ack/close workflow
- **Services catalog** service registration with environment, owner, tags
- **Dashboards** customizable multi-panel dashboards
- **SLO tracking** error budget calculation, MTTR, MTTD analytics
- **Multi-tenant** data isolation by tenant

### Tech Stack
- Go 1.22 + Gin + PostgreSQL 16 + Redis 7
- Multi-tenant, RBAC, JWT auth
- Docker Compose

### Quick Start
```bash
docker-compose up -d
go run cmd/api/main.go
```

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/v1/auth/login | Login |
| POST | /api/v1/traces/ingest | Ingest trace |
| GET | /api/v1/traces | Query traces |
| GET | /api/v1/traces/:id | Get trace detail |
| POST | /api/v1/logs/ingest | Ingest log |
| GET | /api/v1/logs | Query logs |
| POST | /api/v1/metrics/ingest | Ingest metric |
| POST | /api/v1/metrics/query | Query metrics |
| GET | /api/v1/alerts/rules | List alert rules |
| GET | /api/v1/alerts/events | List alert events |

---

<a id="中文"></a>
# 企业级监控告警平台

基于 Go + Gin + PostgreSQL + Redis 构建的分布式追踪、日志聚合与指标告警平台。

### 功能特性
- **分布式追踪** 兼容 OpenTelemetry 的 span 摄取与 trace 重建
- **日志聚合** 结构化日志摄取、全文与字段级查询
- **指标** counter/gauge/histogram、多维标签
- **告警** 基于规则的告警、严重级别、路由、确认/关闭
- **服务目录** 服务注册、环境、负责人、标签
- **仪表板** 可定制的多面板仪表板
- **SLO 跟踪** 错误预算计算、MTTR、MTTD 分析
- **多租户** 按租户数据隔离

### 技术栈
- Go 1.22 + Gin + PostgreSQL 16 + Redis 7
- 多租户、RBAC、JWT 认证
- Docker Compose

### 快速开始
```bash
docker-compose up -d
go run cmd/api/main.go
```
