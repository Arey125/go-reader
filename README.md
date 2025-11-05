# Go Reader

## Installation
1. Get [tailwind-cli](https://tailwindcss.com/blog/standalone-cli)
```bash
curl -LO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
```
2. Install [migrate](https://github.com/golang-migrate/migrate)
```bash
go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
3. Copy .env.example to .env
4. Run migrations with `make migrate`
5. Run project with `make run`
