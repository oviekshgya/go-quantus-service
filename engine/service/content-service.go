package service

import (
	"github.com/gin-gonic/gin"
	"go-quantus-service/src/entities"
	"go-quantus-service/src/pkg"
	"go-quantus-service/src/rabbitMQ"
	rds "go-quantus-service/src/redis"
	"go-quantus-service/src/repository"
	"gorm.io/gorm"
)

type ContentServiceImpl struct {
	DB     *gorm.DB
	Repo   repository.RepositoryContent
	Rabbit *rabbitMQ.RabbitMQImpl
	Redis  *rds.RedisClient
}

func (s *ContentServiceImpl) RegisterContent(c *gin.Context, req *entities.Content) (*entities.Content, error) {
	result, err := pkg.WithTransaction(s.DB, func(tx *gorm.DB) (interface{}, error) {
		cerate := s.Repo.CreateContent(tx, req)
		if cerate != nil {
			return nil, cerate
		}
		return req, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*entities.Content), nil
}
