package service

import (
	"errors"
	"fmt"
	"shdp/internal/model"
	"shdp/internal/repository"
	"shdp/pkg/utils"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

type BlogService struct {
	repo     *repository.BlogRepo
	rdb      *redis.Client
	userRepo *repository.UserRepo
}

func NewBlogService(repo *repository.BlogRepo, rdb *redis.Client, userRepo *repository.UserRepo) *BlogService {
	return &BlogService{repo: repo, rdb: rdb, userRepo: userRepo}
}
func (s *BlogService) LikeBlog(ctx context.Context, blogId int64, userId int64) error {
	key := fmt.Sprintf("%s%d", utils.BlogLikedKey, blogId)
	member := strconv.FormatInt(userId, 10)
	_, err := s.rdb.ZScore(ctx, key, member).Result()
	if err == redis.Nil {
		err := s.repo.UpdateBlogLiked(blogId, 1)
		if err != nil {
			return errors.New("更新数据库失败")
		}
		s.rdb.ZAdd(ctx, key, redis.Z{
			Member: member,
			Score:  float64(time.Now().UnixNano() / 1e6),
		})
		return nil
	} else if err != nil {
		return err
	}
	err = s.repo.UpdateBlogLiked(blogId, -1)
	if err != nil {
		return errors.New("更新数据库失败")
	}
	s.rdb.ZRem(ctx, key, member)
	return nil
}
func (s *BlogService) QueryBlogLikes(ctx context.Context, blogId int64) ([]model.UserDTO, error) {
	key := fmt.Sprintf("%s%d", utils.BlogLikedKey, blogId)

	// 1. 从 ZSET 中查询 top5 的点赞用户 zrange key 0 4
	top5IdsStr, err := s.rdb.ZRange(ctx, key, 0, 4).Result()
	if err != nil || len(top5IdsStr) == 0 {
		return []model.UserDTO{}, nil
	}

	// 2. 解析 ID
	var ids []int64
	for _, idStr := range top5IdsStr {
		id, _ := strconv.ParseInt(idStr, 10, 64)
		ids = append(ids, id)
	}

	// 3. 根据 ID 查数据库，并转为 DTO (这里调用了刚才加的有序查询方法)
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
