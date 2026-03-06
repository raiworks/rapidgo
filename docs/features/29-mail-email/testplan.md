# Feature #29 — Mail / Email: Test Plan

## Test cases

| TC | Description | Method | Expected |
|----|-------------|--------|----------|
| TC-01 | NewMailer reads env vars | Set env, call NewMailer() | All fields populated |
| TC-02 | NewMailer defaults port to 587 when unset | Unset MAIL_PORT | Port == 587 |
| TC-03 | NewMailer parses port from string | Set MAIL_PORT=2525 | Port == 2525 |
| TC-04 | NewMailer handles invalid port gracefully | Set MAIL_PORT=abc | Port == 587 (fallback) |
| TC-05 | Send returns error with invalid host | Call Send() with bad SMTP config | Non-nil error |

## Notes

- TC-01 to TC-04 test configuration wiring without sending real emails.
- TC-05 verifies that `Send()` fails gracefully when SMTP connection cannot be established (uses unreachable host).
- No live SMTP server required — all tests use invalid/unreachable hosts.
