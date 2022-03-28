package router

import (
	"errors"
	"go-template/model"
	"go-template/pkg/hash"
	"go-template/pkg/mysql"
	"go-template/pkg/redis"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db = mysql.Db

type User = model.User
type Users = []model.User

func init() {
	userGroup := Router.Group("/user")
	{
		// @获取用户列表
		userGroup.GET("/list", func(c *gin.Context) {
			users := &Users{}
			result := db.Find(users)
			c.JSON(http.StatusOK, gin.H{
				"ok":        true,
				"user_list": *users,
				"number":    result.RowsAffected,
			})
		})

		// @创建用户
		userGroup.POST("", func(c *gin.Context) {
			user := &User{}
			if err := c.ShouldBindJSON(user); err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"ok":     false,
					"msg":    "Invalid json provided",
					"detail": err,
				})
				return
			}
			user.Password = hash.Encode(user.Password + user.Email)
			exist := &Users{}
			result := db.Where(user).Find(exist)
			if result.RowsAffected != 0 {
				c.JSON(http.StatusOK, gin.H{
					"ok":     false,
					"msg":    "email has been register.",
					"detail": errors.New("could not make a new user with exist email"),
				})
				return
			}
			err := db.Create(user).Error
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"ok":     false,
					"msg":    "datebase handle error",
					"detail": err,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"ok":   true,
				"user": user,
			})
		})

		// @获取用户信息
		userGroup.GET("", func(c *gin.Context) {
			token := c.DefaultQuery("token", "")
			email, err := redis.CheckJWT(token)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"ok":     false,
					"msg":    "token is not exist",
					"detail": err,
				})
				return
			}
			user := &User{Email: email}
			err = db.Where(user).First(user).Error
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"ok":     false,
					"msg":    "datebase handle error",
					"detail": err,
				})
			}
			c.JSON(http.StatusOK, gin.H{
				"ok":   true,
				"user": user,
			})

		})

		// @更新用户信息
		userGroup.PUT("", func(c *gin.Context) {
			token := c.DefaultQuery("token", "")
			email, err := redis.CheckJWT(token)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"ok":     false,
					"msg":    "token is not exist",
					"detail": err,
				})
				return
			}
			user := &User{Email: email}
			db.Where(user).First(user)
			err = c.ShouldBindJSON(user)
			user.Email = email
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"ok":     false,
					"msg":    "Invalid json provided",
					"detail": err,
				})
				return
			}
			if user.Password != "" {
				user.Password = hash.Encode(user.Password + user.Email)
			}
			err = db.Save(user).Error
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"ok":     false,
					"msg":    "datebase handle error",
					"detail": err,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"ok":   true,
				"user": user,
			})
		})

		// @删除用户
		userGroup.DELETE("", func(c *gin.Context) {
			token := c.DefaultQuery("token", "")
			email, err := redis.CheckJWT(token)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"ok":     false,
					"msg":    "token is not exist",
					"detail": err,
				})
				return
			}
			user := &User{Email: email}
			err = db.Find(user).First(user).Error
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"ok":     false,
					"msg":    "datebase handle error",
					"detail": err,
				})
				return
			}
			_, err = redis.DeleteJwt(token)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"ok":     false,
					"msg":    "token is not exist",
					"detail": err,
				})
				return
			}
			err = db.Delete(user).Error
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"ok":     false,
					"msg":    "datebase handle error",
					"detail": err,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"ok":   true,
				"user": user,
			})
		})
	}

	tokenGroup := Router.Group("/token")
	{
		// @验证token
		tokenGroup.GET("", func(c *gin.Context) {
			token := c.DefaultQuery("token", "")
			email, err := redis.CheckJWT(token)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"ok":     false,
					"msg":    "token is not exist",
					"detail": err,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"ok":    true,
				"email": email,
			})
		})

		// @获取token
		tokenGroup.POST("", func(c *gin.Context) {
			user := &User{}
			err := c.ShouldBindJSON(user)
			user.Password = hash.Encode(user.Password + user.Email)
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"ok":     false,
					"msg":    "Invalid json provided",
					"detail": err,
				})
				return
			}
			if err = db.Where(user).First(user).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					c.JSON(http.StatusOK, gin.H{
						"ok":     false,
						"msg":    "Could not found user",
						"detail": err,
					})
					return
				} else {
					c.JSON(http.StatusOK, gin.H{
						"ok":     false,
						"msg":    "datebase handle error",
						"detail": err,
					})
					return
				}
			}
			token, err := redis.AddJwt(user.Email, user.ID)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"ok":     false,
					"msg":    "Could not generate token",
					"detail": err,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"ok":    true,
				"token": token,
			})
		})

		// @删除token
		tokenGroup.DELETE("", func(c *gin.Context) {
			token := c.DefaultQuery("token", "")
			email, err := redis.DeleteJwt(token)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"ok":     false,
					"msg":    "token is not exist",
					"detail": err,
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"ok":    true,
				"email": email,
			})
		})
	}
}
