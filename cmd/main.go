package main

//Система заявок на производстве.
//Услуги - используемые программы станков с ЧПУ,
//заявки - заказ на изготовление детали

import (
	"web/internal/api"
	mClient "web/internal/minio"
	"web/internal/redis"
	"web/internal/repo"

	_ "web/docs"
)

// @title           Система заявок на производстве
// @version         1.0

// @contact.name   Корецкий К.В.
// @contact.url    https://github.com/korzck
// @contact.email  konstantin.koretskiy@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  Курс РИП ИУ5
// @externalDocs.url          https://github.com/iu5git/Web/
func main() {

	db, _ := repo.NewRepo()
	redisRepo, _ := redis.New()
	minioClient := mClient.NewMinioClient()
	srv := api.NewService(db, minioClient, redisRepo)

	srv.Run()
}
