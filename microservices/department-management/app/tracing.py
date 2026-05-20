import os
import logging
from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.exporter.zipkin.json import ZipkinExporter
from opentelemetry.sdk.resources import SERVICE_NAME, Resource
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.instrumentation.requests import RequestsInstrumentor


def initialize_tracer(service_name: str) -> TracerProvider:
    """Initialize OpenTelemetry tracing for FastAPI service"""
    
    # Get Zipkin endpoint from environment or use default
    zipkin_endpoint = os.getenv(
        "OTEL_EXPORTER_ZIPKIN_ENDPOINT",
        "http://zipkin:9411/api/v2/spans"
    )
    
    # Create Zipkin exporter
    zipkin_exporter = ZipkinExporter(
        endpoint=zipkin_endpoint,
    )
    
    # Create resource
    resource = Resource(attributes={
        SERVICE_NAME: service_name,
        "service.version": "1.0.0",
    })
    
    # Create TracerProvider
    trace_provider = TracerProvider(resource=resource)
    trace_provider.add_span_processor(BatchSpanProcessor(zipkin_exporter))
    
    # Set global TracerProvider
    trace.set_tracer_provider(trace_provider)
    
    logging.info(f"Tracing initialized for service: {service_name}")
    return trace_provider


def instrument_app(app):
    """Instrument FastAPI app and requests library"""
    FastAPIInstrumentor.instrument_app(app)
    RequestsInstrumentor().instrument()


def shutdown_tracer(trace_provider: TracerProvider):
    """Gracefully shutdown tracer"""
    trace_provider.force_flush()
