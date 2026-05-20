import pino from "pino";
import pinoHttp from "pino-http";
import * as prometheus from "prom-client";

// Create Prometheus HTTP metrics
export const httpRequestDuration = new prometheus.Histogram({
  name: "http_request_duration_seconds",
  help: "Duration of HTTP requests in seconds",
  labelNames: ["method", "route", "status_code"],
  buckets: [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10],
});

export const httpRequestTotal = new prometheus.Counter({
  name: "http_requests_total",
  help: "Total number of HTTP requests",
  labelNames: ["method", "route", "status_code"],
});

export const httpRequestSize = new prometheus.Histogram({
  name: "http_request_size_bytes",
  help: "Size of HTTP request payloads in bytes",
  labelNames: ["method", "route"],
  buckets: [100, 500, 1000, 5000, 10000, 50000, 100000, 500000],
});

export const httpResponseSize = new prometheus.Histogram({
  name: "http_response_size_bytes",
  help: "Size of HTTP response payloads in bytes",
  labelNames: ["method", "route", "status_code"],
  buckets: [100, 500, 1000, 5000, 10000, 50000, 100000, 500000],
});

// Middleware to collect HTTP metrics
export function metricsMiddleware(req: any, res: any, next: any) {
  const startTime = Date.now();
  const startHrTime = process.hrtime();

  // Record request size
  const contentLength = req.headers["content-length"];
  if (contentLength) {
    httpRequestSize.labels(req.method, req.route?.path || req.path || "/").observe(parseInt(contentLength, 10));
  }

  // Capture response finish
  const originalSend = res.send;
  res.send = function (data: any) {
    const route = req.route?.path || req.path || "/";
    const statusCode = res.statusCode;

    // Record response size
    if (data) {
      const responseSize = typeof data === "string" ? Buffer.byteLength(data) : JSON.stringify(data).length;
      httpResponseSize.labels(req.method, route, statusCode.toString()).observe(responseSize);
    }

    // Record duration and request count
    const hrTime = process.hrtime(startHrTime);
    const duration = hrTime[0] + hrTime[1] / 1e9;
    httpRequestDuration.labels(req.method, route, statusCode.toString()).observe(duration);
    httpRequestTotal.labels(req.method, route, statusCode.toString()).inc();

    // Call original send
    return originalSend.call(this, data);
  };

  next();
}

// Create base logger with JSON format
export const logger = pino({
  level: process.env.LOG_LEVEL || "info",
  transport: {
    target: "pino/file",
    options: {
      destination: 1, // stdout
    },
  },
  formatters: {
    level: (label) => {
      return { level: label.toUpperCase() };
    },
  },
  timestamp: pino.stdTimeFunctions.isoTime,
});

// Create HTTP middleware logger
export const httpLogger = pinoHttp();

// Add trace ID injection to all logs (no-op - tracing disabled)
export function injectTraceID(req: any, res: any, next: any) {
  // OpenTelemetry tracing disabled - this is a placeholder
  // When tracing is re-enabled, this will extract traceId from span context
  next();
}

// Helper to get logger with trace context (simplified version without OpenTelemetry)
export function getLogger() {
  return logger.child({ service: "profile-management" });
}
