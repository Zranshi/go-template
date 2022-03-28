package redis

import (
	"context"
	"errors"
	"go-template/pkg/jwt"
	"strconv"
	"time"
)

func AddJwt(email string, id uint) (string, error) {
	ctx := context.Background()
	key := strconv.FormatInt(time.Now().Unix(), 10)

	// 查找当前email是否存在token，如果存在则删除
	existToken, _ := r.HGet(ctx, "email_set", email).Result()
	if existToken != "" {
		r.HDel(ctx, "tokens", existToken)
	}

	// 生成新的token
	token, err := jwt.GenerateToken(email, id, key)
	if err != nil {
		return "", err
	}

	// 分别插入两个HASH中
	err = r.HSet(ctx, "tokens", token, key).Err()
	if err != nil {
		return "", err
	}
	err = r.HSet(ctx, "email_set", email, token).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

func CheckJWT(token string) (string, error) {
	ctx := context.Background()

	// 查询token是否存在
	key, err := isExistJwt(ctx, token)
	if err != nil {
		return "", err
	}

	// 如果存在再判断时候过期, 如果过期了则删除
	claims, err := jwt.ParseToken(token, key)
	if err != nil {
		return "", err
	}
	if claims.ExpiresAt <= time.Now().Unix() {
		// 过期了，删除两个HASH中的键值对
		if err := deleteJwt(ctx, token, claims.Email); err != nil {
			return "", errors.New("delete token fail")
		}
	}

	// 如果均成功则返回email
	return claims.Email, nil
}

func DeleteJwt(token string) (string, error) {
	ctx := context.Background()

	// 查询token是否存在
	key, err := isExistJwt(ctx, token)
	if err != nil {
		return "", err
	}

	claims, err := jwt.ParseToken(token, key)
	if err != nil {
		return "", err
	}
	// 删除两个HASH中的键值对
	if err := deleteJwt(ctx, token, claims.Email); err != nil {
		return "", errors.New("delete token fail")
	}

	// 如果均成功则返回email
	return claims.Email, nil
}

// 检查token是否存在于redis中
func isExistJwt(ctx context.Context, token string) (string, error) {
	// 查询token是否存在
	key, err := r.HGet(ctx, "tokens", token).Result()
	if err != nil {
		return "", err
	}
	if key == "" {
		return "", errors.New("token is not exist")
	}
	return key, nil
}

// 删除token
func deleteJwt(ctx context.Context, token string, email string) error {
	err := r.HDel(ctx, "tokens", token).Err()
	if err != nil {
		return err
	}

	err = r.HDel(ctx, "email_set", email).Err()
	if err != nil {
		return err
	}
	return nil
}
