

package mw

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/chat/pkg/common/constant"
	"github.com/openimsdk/chat/pkg/protocol/admin"
	constant2 "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
)

func New(client admin.AdminClient) *MW {
	return &MW{client: client}
}

type MW struct {
	client admin.AdminClient
}

func (o *MW) parseToken(c *gin.Context) (string, int32, string, error) {
	token := c.GetHeader("token")
	if token == "" {
		return "", 0, "", errs.ErrArgs.WrapMsg("token is empty")
	}
	resp, err := o.client.ParseToken(c, &admin.ParseTokenReq{Token: token})
	if err != nil {
		return "", 0, "", err
	}
	return resp.UserID, resp.UserType, token, nil
}

func (o *MW) parseTokenType(c *gin.Context, userType int32) (string, string, error) {
	userID, t, token, err := o.parseToken(c)
	if err != nil {
		return "", "", err
	}
	if t != userType {
		return "", "", errs.ErrArgs.WrapMsg("token type error")
	}
	return userID, token, nil
}

func (o *MW) isValidToken(c *gin.Context, userID string, token string) error {
	resp, err := o.client.GetUserToken(c, &admin.GetUserTokenReq{UserID: userID})
	if err != nil {
		return err
	}
	if len(resp.TokensMap) == 0 {
		return errs.ErrTokenExpired.Wrap()
	}
	if v, ok := resp.TokensMap[token]; ok {
		switch v {
		case constant2.NormalToken:
		case constant2.KickedToken:
			return errs.ErrTokenExpired.Wrap()
		default:
			return errs.ErrTokenUnknown.Wrap()
		}
	} else {
		return errs.ErrTokenExpired.Wrap()
	}
	return nil
}

func (o *MW) setToken(c *gin.Context, userID string, userType int32) {
	SetToken(c, userID, userType)
}

func (o *MW) CheckToken(c *gin.Context) {
	userID, userType, token, err := o.parseToken(c)
	if err != nil {
		c.Abort()
		apiresp.GinError(c, err)
		return
	}
	if err := o.isValidToken(c, userID, token); err != nil {
		c.Abort()
		apiresp.GinError(c, err)
		return
	}
	o.setToken(c, userID, userType)
}

func (o *MW) CheckAdmin(c *gin.Context) {
	userID, token, err := o.parseTokenType(c, constant.AdminUser)
	if err != nil {
		c.Abort()
		apiresp.GinError(c, err)
		return
	}
	if err := o.isValidToken(c, userID, token); err != nil {
		c.Abort()
		apiresp.GinError(c, err)
		return
	}
	o.setToken(c, userID, constant.AdminUser)
}

func (o *MW) CheckUser(c *gin.Context) {
	userID, token, err := o.parseTokenType(c, constant.NormalUser)
	if err != nil {
		c.Abort()
		apiresp.GinError(c, err)
		return
	}
	if err := o.isValidToken(c, userID, token); err != nil {
		c.Abort()
		apiresp.GinError(c, err)
		return
	}
	o.setToken(c, userID, constant.NormalUser)
}

func (o *MW) CheckAdminOrNil(c *gin.Context) {
	defer c.Next()
	userID, userType, _, err := o.parseToken(c)
	if err != nil {
		return
	}
	if userType == constant.AdminUser {
		o.setToken(c, userID, constant.AdminUser)
	}
}

func SetToken(c *gin.Context, userID string, userType int32) {
	c.Set(constant.RpcOpUserID, userID)
	c.Set(constant.RpcOpUserType, []string{strconv.Itoa(int(userType))})
	c.Set(constant.RpcCustomHeader, []string{constant.RpcOpUserType})
}
