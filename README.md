### Для запуска потребуется

#### .env
```.env
IS_READ_CONFIG=true
```

#### config.json
```json
{
  "server": {
    "server_port": "1941",
    "server_mode": "development",
    "server_domain": "http://localhost:3000"
  },
  "business-database": {
    "db_password":ваш_пароль,
    "db_host": "localhost",
    "db_port": "5433",
    "db_username":ваше_имя,
    "db_name":ваше_название_бд,
    "db_ssl_mode": "disable"
  },
  "minio": {
    "endpoint":"localhost:9000" ,
    "minio_root_user":ваше_имя ,
    "minio_root_password":ваш_пароль,
    "minio_use_ssl": false,
    "minio_bucket_name":ваше_название_бакета
  }
}

```
### Для получения документации
```
./swagger_init.sh
```