# Mintok Roadmap

## Milestone 1 — Platform Foundation

- [x] Project setup
- [x] Docker local stack
- [x] Configuration and structured logging
- [x] PostgreSQL and Redis connectivity

## Milestone 2 — Identity and Access

- [x] Account registration and password login
- [x] JWT access tokens
- [x] Refresh-token rotation
- [x] Project API keys
- [ ] OAuth and SSO

## Milestone 3 — Gateway MVP

- [ ] OpenAI-compatible chat gateway endpoint
- [ ] Provider adapter contract and request normalization
- [ ] Provider credentials and per-project model configuration
- [ ] Timeouts, retries, fallback, and request correlation IDs

## Milestone 4 — Token Optimization

- [ ] Token estimation and request metering
- [ ] Prompt compression with opt-out and protected sections
- [ ] Savings reporting per request and project

## Milestone 5 — Prompt Routing and Caching

- [ ] Policy-based model routing
- [ ] Exact-match response cache
- [ ] Semantic cache and cache controls
- [ ] Budget, latency, and availability fallback rules

## Milestone 6 — Analytics

- [ ] Gateway request logs with redaction
- [ ] Cost, token, latency, compression, and cache metrics
- [ ] Project dashboard and exportable reports

## Milestone 7 — Model Benchmarking

- [ ] Benchmark suite definition
- [ ] Multi-model benchmark execution
- [ ] Quality, cost, latency comparison and routing recommendations

## Milestone 8 — Teams and Governance

- [ ] Organizations, projects, roles, and audit logs
- [ ] Usage limits, quotas, and alerting
- [ ] Data retention and privacy controls

## Milestone 9 — Production Readiness

- [ ] Provider health monitoring and circuit breaking
- [ ] Deployment automation, observability, and backups
- [ ] Security review and load testing

## Current Task

Implement the OpenAI-compatible chat gateway endpoint.
