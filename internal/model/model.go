package model

import "time"

type User struct {
	ID        int64     `json:"id" db:"id"`
	TenantID  int64     `json:"tenant_id" db:"tenant_id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Role      string    `json:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Service struct {
	ID          int64     `json:"id" db:"id"`
	TenantID    int64     `json:"tenant_id" db:"tenant_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Environment string    `json:"environment" db:"environment"`
	Owner       string    `json:"owner" db:"owner"`
	Tags        string    `json:"tags" db:"tags"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Trace struct {
	ID            int64     `json:"id" db:"id"`
	TenantID      int64     `json:"tenant_id" db:"tenant_id"`
	TraceID       string    `json:"trace_id" db:"trace_id"`
	ServiceName   string    `json:"service_name" db:"service_name"`
	OperationName string    `json:"operation_name" db:"operation_name"`
	StartTime     time.Time `json:"start_time" db:"start_time"`
	Duration      int       `json:"duration" db:"duration"` // microseconds
	Status        string    `json:"status" db:"status"`
	SpanCount     int       `json:"span_count" db:"span_count"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type Span struct {
	ID            int64     `json:"id" db:"id"`
	TraceID       string    `json:"trace_id" db:"trace_id"`
	SpanID        string    `json:"span_id" db:"span_id"`
	ParentSpanID  string    `json:"parent_span_id" db:"parent_span_id"`
	ServiceName   string    `json:"service_name" db:"service_name"`
	OperationName string    `json:"operation_name" db:"operation_name"`
	StartTime     time.Time `json:"start_time" db:"start_time"`
	Duration      int       `json:"duration" db:"duration"`
	Tags          string    `json:"tags" db:"tags"`
	Logs          string    `json:"logs" db:"logs"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type LogEntry struct {
	ID         int64     `json:"id" db:"id"`
	TenantID   int64     `json:"tenant_id" db:"tenant_id"`
	Timestamp  time.Time `json:"timestamp" db:"timestamp"`
	Service    string    `json:"service" db:"service"`
	Level      string    `json:"level" db:"level"` // debug, info, warn, error, fatal
	Message    string    `json:"message" db:"message"`
	Logger     string    `json:"logger" db:"logger"`
	Thread     string    `json:"thread" db:"thread"`
	Hostname   string    `json:"hostname" db:"hostname"`
	TraceID    string    `json:"trace_id" db:"trace_id"`
	SpanID     string    `json:"span_id" db:"span_id"`
	Fields     string    `json:"fields" db:"fields"` // JSON
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type Metric struct {
	ID         int64     `json:"id" db:"id"`
	TenantID   int64     `json:"tenant_id" db:"tenant_id"`
	Name       string    `json:"name" db:"name"`
	Type       string    `json:"type" db:"type"` // counter, gauge, histogram, summary
	Value      float64   `json:"value" db:"value"`
	Labels     string    `json:"labels" db:"labels"`
	Timestamp  time.Time `json:"timestamp" db:"timestamp"`
	Service    string    `json:"service" db:"service"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type AlertRule struct {
	ID           int64     `json:"id" db:"id"`
	TenantID     int64     `json:"tenant_id" db:"tenant_id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	Metric       string    `json:"metric" db:"metric"`
	Condition    string    `json:"condition" db:"condition"` // >, <, ==, !=
	Threshold    float64   `json:"threshold" db:"threshold"`
	Duration     int       `json:"duration" db:"duration"`
	Severity     string    `json:"severity" db:"severity"` // info, warning, critical
	Channels     string    `json:"channels" db:"channels"`
	Enabled      bool      `json:"enabled" db:"enabled"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type AlertEvent struct {
	ID         int64      `json:"id" db:"id"`
	TenantID   int64      `json:"tenant_id" db:"tenant_id"`
	RuleID     int64      `json:"rule_id" db:"rule_id"`
	RuleName   string     `json:"rule_name" db:"rule_name"`
	Severity   string     `json:"severity" db:"severity"`
	Message    string     `json:"message" db:"message"`
	Value      float64    `json:"value" db:"value"`
	Status     string     `json:"status" db:"status"` // firing, acknowledged, resolved
	AckBy      *int64     `json:"ack_by" db:"ack_by"`
	AckAt      *time.Time `json:"ack_at" db:"ack_at"`
	ResolvedAt *time.Time `json:"resolved_at" db:"resolved_at"`
	FiredAt    time.Time  `json:"fired_at" db:"fired_at"`
}

type Dashboard struct {
	ID         int64     `json:"id" db:"id"`
	TenantID   int64     `json:"tenant_id" db:"tenant_id"`
	Name       string    `json:"name" db:"name"`
	Layout     string    `json:"layout" db:"layout"` // JSON
	Owner      int64     `json:"owner" db:"owner"`
	IsPublic   bool      `json:"is_public" db:"is_public"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
