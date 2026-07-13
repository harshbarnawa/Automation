# Mintok Architecture

## Overview

Mintok is a multi-tenant AI gateway. Client applications use a Mintok API key (and dashboard users use JWT-based authentication) to submit AI requests through one API. Mintok optimizes and routes requests, checks caches, invokes the selected provider, and records usage data.

```text
Dashboard / Client Application
            |
            v
       Mintok API Gateway
            |
   Authentication & API Keys
            |
   Policy / Token Optimization
      |       |          |
      v       v          v
  Cache    Router    Analytics Pipeline
              |
              v
       Provider Adapters
     OpenAI / Anthropic / others
```

## Backend Layers

- **HTTP handlers:** validate and normalize external requests, authenticate callers, and return stable gateway responses.
- **Services:** apply compression, routing, cache, benchmarking, and analytics policies.
- **Repositories:** isolate PostgreSQL persistence for tenants, keys, policies, benchmark runs, and aggregated metrics.
- **Infrastructure:** Redis backs short-lived request state and caching; PostgreSQL remains the source of truth; provider adapters isolate vendor SDKs and protocols.

## Core Request Path

1. Authenticate a dashboard user with JWT or a client request with a project API key.
2. Resolve the tenant, project, gateway policy, and model options.
3. Estimate tokens and optionally compress the prompt.
4. Look up an eligible cache entry.
5. Select a provider/model and invoke it through an adapter; use policy-approved fallback on failure.
6. Normalize the response, persist usage metrics, and return metadata about routing, cache, and optimization.

## Data Ownership

- PostgreSQL stores users, organizations, projects, API-key metadata, policies, provider configurations, request records, benchmark suites/runs, and aggregates.
- Redis stores cache entries, rate-limit counters, and short-lived idempotency or routing state.
- Provider credentials are encrypted at rest when persistence is required and are never included in logs or analytics payloads.

## Security Boundaries

- Dashboard authentication is separate from project API-key authentication.
- Every gateway request is scoped to a project and organization before a policy or provider configuration is read.
- Prompts and completions are redacted or retained only according to project policy.
- Provider failures and retries must not expose credentials or prompt content in errors.

## Existing Foundation

The Gin API, PostgreSQL, Redis, Docker stack, structured logging, user repository, and authentication module are retained. GitHub/repository-analysis integrations are intentionally not part of the target architecture.
