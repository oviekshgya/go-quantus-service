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

func (s *ContentServiceImpl) GetAllContents(c *gin.Context, id int64) ([]entities.Content, error) {
	result, err := pkg.WithTransaction(s.DB, func(tx *gorm.DB) (interface{}, error) {
		data, err := s.Repo.ListContentByUserID(tx, id)
		if err != nil {
			return nil, err
		}
		return data, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]entities.Content), nil
}

func (s *ContentServiceImpl) GetContentByID(c *gin.Context, id int64) (*entities.Content, error) {
	result, err := pkg.WithTransaction(s.DB, func(tx *gorm.DB) (interface{}, error) {
		data, err := s.Repo.GetContentByID(tx, id)
		if err != nil {
			return nil, err
		}
		return data, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*entities.Content), nil
}

func (s *ContentServiceImpl) UpdateContent(c *gin.Context, req *entities.Content) (*entities.Content, error) {
	result, err := pkg.WithTransaction(s.DB, func(tx *gorm.DB) (interface{}, error) {
		err := s.Repo.UpdateContent(tx, req)
		if err != nil {
			return nil, err
		}
		return req, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*entities.Content), nil
}

func (s *ContentServiceImpl) DeleteContent(c *gin.Context, id int64) (*int64, error) {
	result, err := pkg.WithTransaction(s.DB, func(tx *gorm.DB) (interface{}, error) {
		err := s.Repo.DeleteContent(tx, id)
		if err != nil {
			return nil, err
		}
		return &id, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*int64), nil
}
