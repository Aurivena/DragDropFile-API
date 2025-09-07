## Для запуска потребуется

```bash
    yay -S nodejs
    yay -S npm
```
#### Установка зависимостей
```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    go install github.com/pressly/goose/v3/cmd/goose@latest
```
---
#### Создаем configs/config.json
```json
{
  "server": {
    "server_port": "1941",
    "server_mode": "development",
    "server_domain": "http://localhost:5173"
  },
  "business-database": {
    "db_password":"answer_postgres_password",
    "db_host": "localhost",
    "db_port": "5433",
    "db_username":"answer_postgres_user",
    "db_name":"answer_postgres_database",
    "db_ssl_mode": "disable"
  },
  "minio": {
    "endpoint":"localhost:9000",
    "minio_root_user":"answer-s3_minio-user",
    "minio_root_password":"answer-s3_minio-password",
    "minio_use_ssl": false,
    "minio_bucket_name":"answer-s3_minio-bucket"
  },
  "certificates": {
    "certificatesPath":"" ,
    "keyPath":""
  }
}
```
---
### Для получения документации
```bash
    make swagger
```