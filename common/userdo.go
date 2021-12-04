package common

import "github.com/garyburd/redigo/redis"

type UserDo struct {
	pool *redis.Pool
}

func NewUserDo(pool *redis.Pool) (userDO *UserDo) {
	userDO = &UserDo{
		pool: pool,
	}
	return
}
