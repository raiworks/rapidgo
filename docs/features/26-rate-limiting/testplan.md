# Feature #26 — Rate Limiting: Test Plan

## Test cases

| TC | Description | Method | Expected |
|----|-------------|--------|----------|
| TC-32 | Default rate limit allows requests within limit | GET | 200 OK |
| TC-33 | X-RateLimit-Limit header present | GET | Header = rate limit value |
| TC-34 | X-RateLimit-Remaining header decrements | GET ×2 | Second response has lower Remaining |
| TC-35 | Requests exceeding limit return 429 | GET ×(limit+1) | 429 Too Many Requests |
| TC-36 | Custom RATE_LIMIT env var is respected | SET env, GET | Limit header matches custom value |
| TC-37 | Middleware alias "ratelimit" resolves | Resolve | Non-nil handler |

## Notes

- TC-32 to TC-35 use a low rate limit (e.g., `"2-M"`) in tests for quick exhaustion.
- TC-36 sets `RATE_LIMIT` env var before middleware creation, then restores it.
- The `ulule/limiter` library handles the 429 response body; tests verify status code and headers.
- Numbering continues from TC-31 (CORS tests in Feature #25).
