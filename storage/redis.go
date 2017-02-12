//
//  redis.go
//  mercuryx
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package storage

import (
	"gopkg.in/redis.v3"
	"github.com/ortuman/mercury/config"
)

type Redis struct {
	client *redis.Client
}

func NewRedis() *Redis {
	r := new(Redis)
	r.client = redis.NewClient(&redis.Options{Addr: config.Redis.Host})
	return r
}

func (r *Redis) IncreaseBadge(senderID, token string) error {
	if err := r.client.Incr("badges:" + senderID + ":" + token).Err(); err != nil && err != redis.Nil {
		return err
	}
	return nil
}

func (r *Redis) GetBadge(senderID, token string) (uint64, error) {
	badge, err := r.client.Get("badges:" + senderID + ":" + token).Uint64()
	if err != nil && err != redis.Nil {
		return 0, err
	}
	return badge, nil
}

func (r *Redis) ClearBadge(senderID, token string) error {
	if err := r.client.Del("badges:" + senderID + ":" + token).Err(); err != nil && err != redis.Nil {
		return err
	}
	return nil
}
