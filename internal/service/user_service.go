package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"shdp/internal/model"
	"shdp/internal/repository"
	"shdp/pkg/utils"
	"strconv"
	"time"
)

type UserService struct {
	repo *repository.UserRepo
	rdb  *redis.Client
}

func NewUserService(repo *repository.UserRepo, rdb *redis.Client) *UserService {
	return &UserService{repo: repo, rdb: rdb}
}
func (s *UserService) SendCode(ctx context.Context, phone string) (string, error) {
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	key := utils.LoginCodeKey + phone
	err := s.rdb.Set(ctx, key, code, utils.LoginCodeTTL).Err()
	if err != nil {
		return "", errors.New("保存验证码失败")
	}
	fmt.Println("【模拟发送短信】向手机号 %s 发送验证码: %s\n", phone, code)
	return code, nil
}
func (s *UserService) Login(ctx context.Context, dto model.LoginFormDTO) (string, error) {
	key := utils.LoginCodeKey + dto.Phone
	cachedCode, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil || cachedCode != dto.Code {
		return "", errors.New("验证码错误或已过期")
	} else if err != nil {
		return "", errors.New("系统异常")
	}
	user, err := s.repo.QueryByPhone(dto.Phone)
	if err != nil {
		return "", errors.New("数据库查询异常")
	}
	if user == nil {
		user = &model.User{
			Phone:    dto.Phone,
			NickName: "user_" + utils.GenerateRandomString(8),
			Icon:     "https://default-icon-url.com/avatar.jpg",
		}
		if err := s.repo.CreateUser(user); err != nil {
			return "", errors.New("创建用户失败")
		}
	}
	token := uuid.New().String()
	userDTO := model.UserDTO{
		ID:       user.ID,
		NickName: user.NickName,
		Icon:     user.Icon,
	}
	userMap := map[string]interface{}{
		"id":       strconv.FormatInt(user.ID, 10),
		"nickName": userDTO.NickName,
		"icon":     userDTO.Icon,
	}
	tokenKey := utils.LoginUserKey + token
	pipe := s.rdb.Pipeline()
	pipe.HSet(ctx, tokenKey, userMap)
	pipe.Expire(ctx, tokenKey, utils.LoginUserTTL)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return "", errors.New("缓存用户信息失败")
	}
	s.rdb.Do(ctx, key)
	return token, nil
}
