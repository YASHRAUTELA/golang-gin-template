# Go Server

Go Version: 1.22.5 or later

## Available Scripts

### Initialize Go modules

```bash
  go mod tidy
```

### Run localhost server (using startup script)
Make the script executable:
```bash
  chmod +x start.sh
```

Running the script:
```bash
  ./start.sh
```

### Run localhost server

In the `server/` directory, you can run:

```bash
  go run .
```

### Update swag documentation

```bash
  swag init
```

### Set up environment variables:

create a .env file in the project root (if it doesn't already exist) and add the necessary variables.

```bash
  DB_HOST=localhost
  DB_USER=yourusername
  DB_PASSWORD=yourpassword
  DB_NAME=yourdbname
  DB_PORT=5432
```

### Integrations:
- Postgres
- Gorm
- Bcrypt
- JWT