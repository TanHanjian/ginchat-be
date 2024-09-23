package middlewares

import (
	"context"
	"encoding/json"
	user_models "ginchat/models/user_basic"
	"ginchat/myredis"
	"ginchat/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从请求头获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "未提供token"})
			c.Abort()
			return
		}

		// 2. 验证token是否合法
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "无效的token"})
			c.Abort()
			return
		}

		// 3. 从Redis中获取用户信息
		ctx := context.Background()
		userJSON, err := myredis.Client.Get(ctx, "user:"+strconv.FormatUint(uint64(claims.ID), 10)).Result()
		if err == nil {
			// 如果Redis中存在用户信息，直接使用
			var user user_models.UserBasic
			err = json.Unmarshal([]byte(userJSON), &user)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "解析用户信息失败"})
				c.Abort()
				return
			}
			c.Set("user", user)
			c.Next()
			return
		}

		// 4. 如果Redis中没有，从数据库获取用户信息
		result, user := user_models.FindByID(claims.ID)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "用户不存在"})
			c.Abort()
			return
		}

		// 5. 将用户信息存储到Redis
		newJSON, err := json.Marshal(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "序列化用户信息失败"})
			c.Abort()
			return
		}
		myredis.Client.Set(ctx, "user:"+strconv.FormatUint(uint64(claims.ID), 10), string(newJSON), 30*time.Minute)

		// 6. 将用户信息存储到上下文
		c.Set("user", user)
		c.Next()
	}
}
