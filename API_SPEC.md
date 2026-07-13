POST /auth/register

POST /auth/login

POST /auth/refresh

GET /auth/me

## Project API keys

All API-key management endpoints require a dashboard JWT bearer token.

POST /projects/:project_id/api-keys

Creates a project-scoped gateway credential. The plaintext `key` is returned only in this response; Mintok stores only its hash. Request body: `{"name": "Production"}`.

GET /projects/:project_id/api-keys

Returns API-key metadata, including a non-sensitive prefix and revocation state.

DELETE /projects/:project_id/api-keys/:key_id

Revokes an active API key. Revoked keys cannot be used by gateway clients.

Planned gateway API

POST /v1/chat/completions

GET /v1/analytics/usage

POST /v1/benchmarks
