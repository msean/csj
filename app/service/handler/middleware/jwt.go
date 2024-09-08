package middleware

import (
	"app/global"
	"app/service/common"
	"errors"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type MyClaims struct {
	Phone string `json:"phone"`
	UUID  string `json:"uuid"`
	jwt.RegisteredClaims
}

func SetToken(phone, uuid string) (token string, err error) {
	SetClaims := MyClaims{
		Phone: phone,
		UUID:  uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(global.Global.TokenEx())), //有效时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                              //签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                              //生效时间
			Issuer:    os.Getenv("JWT_ISSUER"),                                     //签发人
			Subject:   "caishuji",                                                  //主题
			ID:        "1",                                                         //JWT ID用于标识该JWT
			Audience:  []string{"caishuji"},                                        //用户
		},
	}

	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaims)
	token, err = tokenStruct.SignedString([]byte(global.Global.Signkey()))
	return
}

func GetTokenClaims(token string) (claims *MyClaims, err error) {
	tokenObj, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Global.Signkey()), nil
	})

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return
	}
	err = nil

	claims = tokenObj.Claims.(*MyClaims)
	return
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			c.Abort()
			common.Response(c, common.TokenUnValidErr, nil)
			return
		}

		claims, err := GetTokenClaims(tokenHeader)
		global.Global.Logger.Debug("[AuthMiddleware]", zap.Any("expiresat", claims.ExpiresAt))
		if err != nil {
			c.Abort()
			common.Response(c, common.TokenUnValidErr, nil)
			return
		}

		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			if t := claims.ExpiresAt.Time.Add(global.Global.TokenExRefresh()); t.After(time.Now()) {
				newToken, e := SetToken(claims.Phone, claims.UUID)
				if e != nil {
					c.Abort()
					common.Response(c, common.TokenUnGenerateErr, nil)
				}
				c.Header("new-token", newToken)
			} else {
				c.Abort()
				common.Response(c, common.TokenUnValidErr, nil)
				return
			}
		}
		c.Set("userUUID", claims.UUID)
		c.Next()
	}
}
