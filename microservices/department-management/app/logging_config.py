import logging
import json
import sys
from pythonjsonlogger import jsonlogger
from opentelemetry import trace

def initialize(service_name: str) -> logging.Logger:
    """Initialize structured JSON logging for the service."""
    
    # Create logger
    logger = logging.getLogger(service_name)
    logger.setLevel(logging.DEBUG)
    
    # Create JSON formatter
    logHandler = logging.StreamHandler(sys.stdout)
    formatter = jsonlogger.JsonFormatter(
        fmt='%(timestamp)s %(level)s %(service)s %(traceId)s %(message)s',
        timestamp=True
    )
    
    # Custom formatter to extract trace ID
    class TraceIDFilter(logging.Filter):
        def filter(self, record):
            # Get trace ID from OpenTelemetry context
            span = trace.get_current_span()
            if span and span.is_recording():
                record.traceId = span.get_span_context().trace_id
            else:
                record.traceId = ""
            record.service = service_name
            return True
    
    # Add trace ID filter
    logHandler.addFilter(TraceIDFilter())
    logHandler.setFormatter(formatter)
    
    # Remove default handlers and add JSON handler
    logger.handlers.clear()
    logger.addHandler(logHandler)
    
    return logger
