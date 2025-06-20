package models

type ConfigService struct {
	Server       ServerConfig       `json:"server" binding:"required"`
	BusinessDB   BusinessDBConfig   `json:"business-database" binding:"required"`
	Minio        MinioConfig        `json:"minio" binding:"required"`
	Certificates CertificatesConfig `json:"certificates" binding:"required"`
}

type CertificatesConfig struct {
	CertificatesPath string `json:"certificatesPath"`
	KeyPath          string `json:"keyPath"`
}

type ServerConfig struct {
	Port       string `json:"server_port" binding:"required"`
	ServerMode string `json:"server_mode" binding:"required"`
	Domain     string `json:"server_domain" binding:"required"`
}

type BusinessDBConfig struct {
	Password string `json:"db_password" binding:"required"`
	Host     string `json:"db_host" binding:"required"`
	Port     string `json:"db_port" binding:"required"`
	Username string `json:"db_username" binding:"required"`
	DBName   string `json:"db_name" binding:"required"`
	SSLMode  string `json:"db_ssl_mode" binding:"required"`
}

type MinioConfig struct {
	Endpoint        string `json:"endpoint" binding:"required"`
	User            string `json:"minio_root_user" binding:"required"`
	Password        string `json:"minio_root_password" binding:"required"`
	SSL             bool   `json:"minio_use_ssl"  binding:"required"`
	MinioBucketName string `json:"minio_bucket_name" binding:"required"`
}
