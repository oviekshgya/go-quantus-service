package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-quantus-service/src/config"
	"go-quantus-service/src/entities"
	"go-quantus-service/src/pkg"
	"go-quantus-service/src/rabbitMQ"
	rds "go-quantus-service/src/redis"
	"go-quantus-service/src/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	DB     *gorm.DB
	Repo   repository.RepositoryUser
	Rabbit *rabbitMQ.RabbitMQImpl
	Redis  *rds.RedisClient
}

func (s *UserServiceImpl) RegisterUser(c *gin.Context, req *entities.User) (*entities.User, error) {
	result, err := pkg.WithTransaction(s.DB, func(tx *gorm.DB) (interface{}, error) {
		if err := s.Repo.CreateUser(tx, req); err != nil {
			return nil, err
		}
		return req, nil
	})
	if err != nil {
		config.Logger.Println("[err service] register query", err.Error())
		return nil, err
	}
	return result.(*entities.User), nil
}

func (s *UserServiceImpl) LoginUser(c *gin.Context, req *entities.User) (*entities.User, error) {
	result, err := pkg.WithTransaction(s.DB, func(tx *gorm.DB) (interface{}, error) {
		getUser, err := s.Repo.GetUserByEmail(tx, req.Email)
		if err != nil {
			config.Logger.Println("[err service] not found email", err.Error())
			return nil, fmt.Errorf("invalid email/passowrd")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(getUser.Password), []byte(req.Password)); err != nil {
			config.Logger.Println("[err service] invalid password", err.Error())
			return nil, fmt.Errorf("invalid email/passowrd")
		}
		return getUser, nil
	})
	if err != nil {
		config.Logger.Println("err service", err.Error())
		return nil, err
	}
	return result.(*entities.User), nil
}

func (s *UserServiceImpl) UserDetail(c *gin.Context, id int64) (*entities.User, error) {
	result, err := pkg.WithTransaction(s.DB, func(tx *gorm.DB) (interface{}, error) {
		getUser, err := s.Repo.GetUserByID(tx, id)
		if err != nil {
			config.Logger.Println("[err service] not found user", err.Error())
			return nil, fmt.Errorf("not found user")
		}
		return getUser, nil
	})
	if err != nil {
		config.Logger.Println("[err service] not found user", err.Error())
		return nil, fmt.Errorf("invalid user")
	}
	return result.(*entities.User), nil
}

func (s *UserServiceImpl) UpdateUSser(c *gin.Context, req *entities.User) (*entities.User, error) {
	result, err := pkg.WithTransaction(s.DB, func(tx *gorm.DB) (interface{}, error) {
		err := s.Repo.UpdateUser(tx, req)
		if err != nil {
			config.Logger.Println("[err service] update user", err.Error())
			return nil, err
		}
		return req, nil
	})
	if err != nil {
		config.Logger.Println("[err service] update user", err.Error())
		return nil, err
	}
	return result.(*entities.User), nil
}

func (s *UserServiceImpl) DeleteUser(c *gin.Context, id int64) (*int64, error) {
	result, err := pkg.WithTransaction(s.DB, func(tx *gorm.DB) (interface{}, error) {
		err := s.Repo.DeleteUser(tx, id)
		if err != nil {
			config.Logger.Println("[err service] delete user", err.Error())
			return nil, err
		}
		return &id, nil
	})
	if err != nil {
		config.Logger.Println("[err service] delete user", err.Error())
		return nil, err
	}
	return result.(*int64), nil
}
