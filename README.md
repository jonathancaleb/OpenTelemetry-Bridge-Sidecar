# OpenTelemetry Bridge Sidecar  
*A transparent, drop-in reverse-proxy that adds W3C tracing to any HTTP service—no code changes.*

## Why I’m building it  
I’m tackling a **personal challenge**: give any legacy service modern observability **without touching its code**.  
Right now I’m wiring it into a farmer photo-upload API I run in rural Uganda; on 2G links the uploads vanish into silent 502/504 holes while the NestJS monolith has **zero tracing**.  
The sidecar sits unseen between phone and backend, measuring upload latency, body size, retry storms and error rates then ships everything over OTLP to any collector. No redeploy, no restarts.

## Inspiration  
The side-car pattern distilled by Mrinal in [“All about Sidecar”](https://medium.com/@mrinaldoesanything/all-about-sidecar-de79f93565d1) plus OpenTelemetry research papers and a lot of AI-assisted experimentation.

## Quick start  
```bash
# 1. clone & build
git clone https://github.com/you/otel-bridge.git
cd otel-bridge
go run ./cmd/bridge -upstream http://localhost:3000 -otlp localhost:4317 -listen :8080

# 2. send traffic
curl -X POST localhost:8080/upload -F "photo=@maize.jpg"

# 3. open Jaeger
open http://localhost:16686
```

## bones 
```
cmd/bridge          # main()
internal/proxy      # reverse-proxy core
internal/telemetry  # OTel tracer + meter providers
internal/config     # file/env hot-reloader
examples/           # NestJS mock + 2G HTML client
deployments/        # Dockerfile + K8s sidecar patch
```

## Current focus  
- gRPC tracing  
- Prometheus `/metrics` endpoint  
- Back-pressure & circuit-breaker  
- TCP/UDP stream sampling  