# 🚦 High-Throughput Rate Limiter (Gin + Redis + Lua + Prometheus)

A blazing-fast, production-ready rate limiter built in Go (Gin), powered by Redis Lua scripts for atomic token bucket operations, and instrumented with Prometheus & Grafana for real-time monitoring and observability.

✨ Features

⚡ High Throughput – Uses Redis Lua scripts for atomic token bucket updates with minimal latency.

🛡 Distributed & Consistent – Works across multiple Gin instances using Redis as a shared store.

🎛 Configurable – Easily set bucket size, refill rate, and key granularity (per user/IP/endpoint).

📊 Metrics & Observability – Exposes Prometheus metrics at /metrics endpoint for scraping.

🧩 Gin Middleware – Drop-in middleware for any Gin project.

🧪 Tested – Unit tests for Lua script, middleware, and metrics.

# ⚙️ How It Works

1. Token Bucket Algorithm

  - Each request tries to consume a token from its bucket.

  - If available → token is decremented and request proceeds.

  - If empty → request is rejected with 429 Too Many Requests.

2. Redis Lua Script

  - Atomic refill & consume in a single round trip.

  - Prevents race conditions under high concurrency.

3. Prometheus Integration

  - Exposes metrics: allowed requests, rejections, bucket usage, request latency.

  - Compatible with Grafana dashboards for monitoring.
