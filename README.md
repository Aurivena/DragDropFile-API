### Для запуска потребуется

Соаздем config/config.json

#### config.json
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
    "minio_root_user":"answer-minio-user",
    "minio_root_password":"answer-minio-password",
    "minio_use_ssl": false,
    "minio_bucket_name":"answer-minio-bucket"
  },
  "certificates": {
    "certificatesPath":"" ,
    "keyPath":""
  }
}
```
### Для получения документации
```
./swagger_docs.sh
```
