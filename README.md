### Для запуска потребуется

#### config.json
```json
{
  "server": {
    "server_port": "1941",
    "server_mode": "development",
    "server_domain": "http://localhost:3000"
  },
  "business-database": {
    "db_password":пароль,
    "db_host": "localhost",
    "db_port": "5433",
    "db_username":имя,
    "db_name":название_бд,
    "db_ssl_mode": "disable"
  },
  "minio": {
    "endpoint":"localhost:9000" ,
    "minio_root_user":имя ,
    "minio_root_password":пароль,
    "minio_use_ssl": false,
    "minio_bucket_name":название_бакета
  }
}

```
### Для получения документации
```
./swagger_docs.sh
```
