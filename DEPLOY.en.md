# DS2API Deployment Guide (Go)

Language: [中文](DEPLOY.md) | [English](DEPLOY.en.md)

## Contents

- Vercel deployment
- Docker deployment
- Local run
- systemd deployment

## Vercel Deployment

1. Import the repository into Vercel
2. Set required environment variables:
- `DS2API_ADMIN_KEY`
- `DS2API_CONFIG_JSON` (JSON or Base64)
3. Deploy and open `/admin`

The project uses `api/index.go` as the serverless entrypoint. See `vercel.json`.

## Docker Deployment

```bash
cp .env.example .env
# edit .env

docker-compose up -d

docker-compose logs -f
```

Rebuild after updates:

```bash
docker-compose up -d --build
```

## Local Run

```bash
cp config.example.json config.json
# edit config

go run ./cmd/ds2api
```

Default port is `5001` (override with `PORT`).

## systemd Deployment (Linux)

Example unit file:

```ini
[Unit]
Description=DS2API (Go)
After=network.target

[Service]
Type=simple
WorkingDirectory=/opt/ds2api
Environment=PORT=5001
Environment=DS2API_CONFIG_PATH=/opt/ds2api/config.json
Environment=DS2API_ADMIN_KEY=your-admin-secret-key
ExecStart=/opt/ds2api/ds2api
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Useful commands:

```bash
sudo systemctl daemon-reload
sudo systemctl enable ds2api
sudo systemctl start ds2api
sudo systemctl status ds2api
```
