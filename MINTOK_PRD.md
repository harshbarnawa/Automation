# Mintok Product Requirements Document

## Vision

Mintok is an AI Token Optimization Gateway. It gives teams one secure endpoint for their AI traffic and reduces spend without requiring application teams to manage provider-specific routing, prompt optimization, caching, or usage analytics.

## Problem

AI applications frequently overspend on tokens, send repetitive prompts to expensive models, and lack a reliable view of cost, latency, quality, and cache effectiveness across providers. Provider APIs do not offer a unified optimization and governance layer.

## Target Users

- AI product and platform engineers operating production LLM workloads.
- Engineering leaders managing AI cost, reliability, and model quality.
- Developer teams that need a simple, OpenAI-compatible gateway endpoint.

## Core Product

### AI Gateway

- A unified, authenticated API for sending chat and completion requests.
- Provider and model abstraction with request normalization, timeouts, retries, and clear error responses.
- Project-scoped API keys, request limits, and audit logging.

### Token Optimization

- Prompt compression that preserves task-critical instructions and configurable protected content.
- Token estimates, savings estimates, and optimization decisions returned with each request.
- Policy controls to enable, disable, or tune compression by project and route.

### Intelligent Routing

- Route requests according to policy, task type, budget, latency target, context size, and model availability.
- Support explicit model selection, fallback chains, and provider failover.
- Recommend models using benchmark results while keeping routing decisions explainable.

### Caching

- Semantic and exact-match caching for eligible requests.
- Configurable cache keys, TTLs, and bypass controls.
- Cache hit, avoided-token, cost, and latency measurements.

### Benchmarking and Analytics

- Benchmark models against user-defined representative prompts and quality criteria.
- Dashboard views for token volume, cost, latency, cache rate, compression savings, provider/model use, and routing outcomes.
- Request-level logs with sensitive prompt content redacted by default.

## Non-Goals

- Repository cloning, GitHub connection, pull-request review, source-code scanning, documentation generation from repositories, security scanning, and performance analysis are not Mintok product features.
- Mintok does not train foundation models or replace AI providers.

## Key User Flows

### Gateway Request Flow

1. A customer application authenticates with a Mintok API key.
2. Mintok validates policy, estimates tokens, and optionally compresses the prompt.
3. The router checks the cache and selects the best eligible provider/model.
4. Mintok forwards the request, applies fallback when appropriate, and returns a normalized response.
5. Mintok records tokens, cost, latency, cache, compression, and routing metrics.

### Benchmark Flow

1. A workspace defines a prompt suite and evaluation criteria.
2. Mintok runs selected models through the suite.
3. Mintok stores quality, cost, token, and latency results.
4. The workspace uses the results to create or refine routing policies.

## Non-Functional Requirements

- Secure credentials, API keys, and provider secrets; never log raw secrets.
- Tenant isolation for all workspace, project, key, and analytics data.
- Configurable data retention and prompt redaction.
- P95 gateway overhead target below 100 ms excluding provider latency for cache misses; cache hits target below 50 ms.
- Idempotent request handling where supported and observable failures with correlation IDs.

## Pricing Direction

- Free: development gateway usage and basic analytics.
- Pro: optimization policies, caching, benchmarks, and extended retention.
- Enterprise: SSO, advanced governance, private deployment, and negotiated volume.
