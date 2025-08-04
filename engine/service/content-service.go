package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-quantus-service/src/entities"
	"go-quantus-service/src/pkg"
	"go-quantus-service/src/rabbitMQ"
	rds "go-quantus-service/src/redis"
	"go-quantus-service/src/repository"
	"gorm.io/gorm"
	"net/http"
)

type ContentServiceImpl struct {
	DB     *gorm.DB
	Repo   repository.RepositoryContent
	Rabbit *rabbitMQ.RabbitMQImpl
	Redis  *rds.RedisClient
}

// helper generic internal untuk mengurangi boilerplate
func (s *ContentServiceImpl) withTx(c *gin.Context, fn func(tx *gorm.DB) (interface{}, error)) (interface{}, error) {
	return pkg.WithTransaction(s.DB, fn)
}

func (s *ContentServiceImpl) RegisterContent(c *gin.Context, req *entities.Content) (*entities.Content, error) {
	result, err := s.withTx(c, func(tx *gorm.DB) (interface{}, error) {
		if createErr := s.Repo.CreateContent(tx, req); createErr != nil {
			return nil, createErr
		}
		return req, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*entities.Content), nil
}

func (s *ContentServiceImpl) GetAllContents(c *gin.Context, userID int64) ([]entities.Content, error) {
	result, err := s.withTx(c, func(tx *gorm.DB) (interface{}, error) {
		return s.Repo.ListContentByUserID(tx, userID)
	})
	if err != nil {
		return nil, err
	}
	return result.([]entities.Content), nil
}

func (s *ContentServiceImpl) GetContentByID(c *gin.Context, id int64) (*entities.Content, error) {
	result, err := s.withTx(c, func(tx *gorm.DB) (interface{}, error) {
		return s.Repo.GetContentByID(tx, id)
	})
	if err != nil {
		return nil, err
	}
	return result.(*entities.Content), nil
}

func (s *ContentServiceImpl) UpdateContent(c *gin.Context, req *entities.Content) (*entities.Content, error) {
	result, err := s.withTx(c, func(tx *gorm.DB) (interface{}, error) {
		if err := s.Repo.UpdateContent(tx, req); err != nil {
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
	result, err := s.withTx(c, func(tx *gorm.DB) (interface{}, error) {
		if err := s.Repo.DeleteContent(tx, id); err != nil {
			return nil, err
		}
		return &id, nil
	})
	if err != nil {
		return nil, err
	}
	return result.(*int64), nil
}

func (s *ContentServiceImpl) HandleContentUpdateOrDelete(c *gin.Context, content *entities.Content) (interface{}, error) {
	method := c.Request.Method

	return s.withTx(c, func(tx *gorm.DB) (interface{}, error) {
		switch method {
		case http.MethodDelete:
			if err := s.Repo.DeleteContent(tx, content.ID); err != nil {
				return nil, err
			}
			return &content.ID, nil

		case http.MethodPut, http.MethodPatch:
			if err := s.Repo.UpdateContent(tx, content); err != nil {
				return nil, err
			}
			return content, nil

		default:
			return nil, fmt.Errorf("unknown method: %s", method)
		}
	})
}
