# Mintok Data Model

## Identity and Tenancy

- `users` — dashboard accounts and password hashes.
- `organizations` — tenant boundary for teams and billing.
- `team_members` — organization membership and roles.
- `projects` — project-scoped gateway configuration and API keys.
- `api_keys` — project-scoped, hashed gateway client credentials with public prefixes, revocation state, and usage metadata.
- `audit_logs` — security-relevant actions and configuration changes.

## Gateway Configuration

- `provider_configs` — encrypted provider credential references and availability settings.
- `routing_policies` — model preference, budget, latency, and fallback rules.
- `cache_policies` — cache eligibility, key strategy, and TTL settings.
- `compression_policies` — prompt-compression behavior and protected sections.

## Gateway Operations

- `ai_requests` — normalized request metadata, selected provider/model, token counts, cost, latency, routing outcome, and redacted status information.
- `cache_entries` — optional persistent cache metadata; response payloads and hot cache state live in Redis.
- `usage_aggregates` — project and organization metric rollups for analytics.

## Benchmarking

- `benchmark_suites` — representative prompt collections and evaluation rules.
- `benchmark_runs` — suite executions and lifecycle state.
- `benchmark_results` — per-model quality, cost, token, and latency outcomes.

## Existing Schema

The initial PostgreSQL migration retains legacy repository-analysis tables from the pre-pivot prototype. New gateway tables will be introduced additively as their roadmap milestones are implemented; no completed authentication or infrastructure tables are removed.
