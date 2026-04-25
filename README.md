# crawley-contact-api

Contact form backend for [crawley.systems](https://crawley.systems).

- `POST /api/contact` — accepts JSON `{name, email, company, message}`, sends notifications via Pushover + Discord.
- `GET /health` — liveness/readiness probe.

## Configuration

Environment variables:

- `PORT` — listen port (default `8090`)
- `PUSHOVER_TOKEN`, `PUSHOVER_USER` — optional Pushover credentials
- `DISCORD_WEBHOOK_URL` — optional Discord webhook URL

## Build

```
docker build -t crawley-contact-api .
```
