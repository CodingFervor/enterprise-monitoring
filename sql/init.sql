-- Enterprise Monitoring & Alerting - Database Schema
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL DEFAULT 0,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(100), email VARCHAR(100), role VARCHAR(50) DEFAULT 'viewer',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE services (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL DEFAULT 0,
    name VARCHAR(100) NOT NULL, description TEXT,
    environment VARCHAR(50) DEFAULT 'production', owner VARCHAR(100),
    tags JSONB, status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(tenant_id, name, environment)
);

CREATE TABLE traces (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL DEFAULT 0,
    trace_id VARCHAR(64) NOT NULL, service_name VARCHAR(100) NOT NULL,
    operation_name VARCHAR(200) NOT NULL, start_time TIMESTAMP NOT NULL,
    duration BIGINT NOT NULL, status VARCHAR(20), span_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_traces_trace_id ON traces(trace_id);
CREATE INDEX idx_traces_service ON traces(service_name);
CREATE INDEX idx_traces_start ON traces(start_time);

CREATE TABLE spans (
    id BIGSERIAL PRIMARY KEY,
    trace_id VARCHAR(64) NOT NULL, span_id VARCHAR(32) NOT NULL,
    parent_span_id VARCHAR(32), service_name VARCHAR(100) NOT NULL,
    operation_name VARCHAR(200) NOT NULL, start_time TIMESTAMP NOT NULL,
    duration BIGINT NOT NULL, tags JSONB, logs JSONB, status VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_spans_trace ON spans(trace_id);

CREATE TABLE logs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL DEFAULT 0,
    "timestamp" TIMESTAMP NOT NULL, service VARCHAR(100) NOT NULL,
    level VARCHAR(20) NOT NULL, message TEXT, logger VARCHAR(200),
    thread VARCHAR(100), hostname VARCHAR(100),
    trace_id VARCHAR(64), span_id VARCHAR(32), fields JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_logs_ts ON logs("timestamp");
CREATE INDEX idx_logs_service ON logs(service);
CREATE INDEX idx_logs_level ON logs(level);

CREATE TABLE metrics (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL DEFAULT 0,
    name VARCHAR(200) NOT NULL, type VARCHAR(20) NOT NULL,
    value DOUBLE PRECISION NOT NULL, labels JSONB,
    "timestamp" TIMESTAMP NOT NULL, service VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_metrics_name ON metrics(name);
CREATE INDEX idx_metrics_ts ON metrics("timestamp");

CREATE TABLE alert_rules (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL DEFAULT 0,
    name VARCHAR(200) NOT NULL, description TEXT,
    metric VARCHAR(200) NOT NULL, condition VARCHAR(10) NOT NULL,
    threshold DOUBLE PRECISION NOT NULL, duration INT DEFAULT 0,
    severity VARCHAR(20) DEFAULT 'warning', channels JSONB,
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE alert_events (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL DEFAULT 0,
    rule_id BIGINT NOT NULL REFERENCES alert_rules(id),
    rule_name VARCHAR(200), severity VARCHAR(20) NOT NULL,
    message TEXT, value DOUBLE PRECISION,
    status VARCHAR(20) DEFAULT 'firing',
    ack_by BIGINT, ack_at TIMESTAMP, resolved_at TIMESTAMP,
    fired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_alert_events_status ON alert_events(status);

CREATE TABLE dashboards (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL DEFAULT 0,
    name VARCHAR(200) NOT NULL, layout JSONB,
    owner BIGINT, is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
