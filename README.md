# crawley-contact-api

Contact form backend for [crawley.systems](https://crawley.systems).

- `POST /api/contact` — accepts JSON `{name, email, company, message}`, sends notifications via Pushover + Discord.
- `GET /health` — liveness/readiness probe.

## Configuration

Environment variables:

- `PORT` — listen port (default `8090`)
- `PUSHOVER_TOKEN`, `PUSHOVER_USER` — optional Pushover credentials
- `DISCORD_WEBHOOK_URL` — optional Discord webhook URL
- `LLM_API_KEY` — optional. If set, contact submissions are screened by an LLM. Spam-flagged messages still go to Discord (tagged `[SPAM]`) for audit but skip Pushover.
- `LLM_BASE_URL` — OpenAI-compatible chat-completions base (default `https://api.openai.com/v1`). Override for OpenRouter, DeepSeek, vLLM, etc.
- `LLM_MODEL` — model name (default `gpt-4.1-nano`).

## Build

```
docker build -t crawley-contact-api .
```
