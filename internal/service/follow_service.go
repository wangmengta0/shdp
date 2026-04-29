package service

import (
	"context"
	"errors"
	"fmt"
	"shdp/internal/model"
	"shdp/internal/repository"
	"shdp/pkg/utils"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type FollowService struct {
	repo     *repository.FollowRepo
	userRepo *repository.UserRepo
	rdb      *redis.Client
}

func NewFollowService(repo *repository.FollowRepo, userRepo *repository.UserRepo, rdb *redis.Client) *FollowService {
	return &FollowService{repo: repo, userRepo: userRepo, rdb: rdb}
}

func (s *FollowService) IsFollow(userId, followUserId int64) (bool, error) {
	return s.repo.CheckFollow(userId, followUserId)
}

func (s *FollowService) Follow(ctx context.Context, userId, followUserId int64, isFollow bool) error {
	key := fmt.Sprintf("%s%d", utils.FollowsKey, userId)
	member := strconv.FormatInt(followUserId, 10)

	if isFollow {
		if err := s.repo.CreateFollow(userId, followUserId); err != nil {
			return errors.New("关注失败")
		}
		s.rdb.SAdd(ctx, key, member)
	} else {
		if err := s.repo.DeleteFollow(userId, followUserId); err != nil {
			return errors.New("取消关注失败")
		}
		s.rdb.SRem(ctx, key, member)
	}
	return nil
}

func (s *FollowService) CommonFollows(ctx context.Context, userId, targetUserId int64) ([]model.UserDTO, error) {
	key1 := fmt.Sprintf("%s%d", utils.FollowsKey, userId)
	key2 := fmt.Sprintf("%s%d", utils.FollowsKey, targetUserId)

	intersectIds, err := s.rdb.SInter(ctx, key1, key2).Result()
	if err != nil || len(intersectIds) == 0 {
		return []model.UserDTO{}, nil
	}

	var ids []int64
	for _, idStr := range intersectIds {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		ids = append(ids, id)
	}
	
	users, err := s.userRepo.QueryUsersByIdsOrdered(ids)
	if err != nil {
		return nil, errors.New("查询用户信息失败")
	}

	var dtoList []model.UserDTO
	for _, u := range users {
		dtoList = append(dtoList, model.UserDTO{
			ID:       u.ID,
			NickName: u.NickName,
			Icon:     u.Icon,
		})
	}

	return dtoList, nil
}
