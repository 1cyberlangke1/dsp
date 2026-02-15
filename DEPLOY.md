# DS2API 部署指南（Go 版本）

语言 / Language: [中文](DEPLOY.md) | [English](DEPLOY.en.md)

## 目录

- Vercel 部署
- Docker 部署
- 本地运行
- systemd 部署

## Vercel 部署

1. 导入仓库到 Vercel
2. 设置环境变量（至少）：
- `DS2API_ADMIN_KEY`
- `DS2API_CONFIG_JSON`（JSON 或 Base64）
3. 部署后访问 `/admin` 管理界面

说明：项目使用 `api/index.go` 作为 Serverless 入口，配置见 `vercel.json`。

## Docker 部署

```bash
cp .env.example .env
# 编辑 .env

docker-compose up -d

docker-compose logs -f
```

更新后重建：

```bash
docker-compose up -d --build
```

## 本地运行

```bash
cp config.example.json config.json
# 编辑配置

go run ./cmd/ds2api
```

默认端口 `5001`，可通过 `PORT` 环境变量覆盖。

## systemd 部署（Linux）

示例服务文件：

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

常用命令：

```bash
sudo systemctl daemon-reload
sudo systemctl enable ds2api
sudo systemctl start ds2api
sudo systemctl status ds2api
```
