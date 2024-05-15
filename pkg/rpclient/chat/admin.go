

package chat

import (
	"context"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/eerrs"
	"github.com/openimsdk/chat/pkg/protocol/admin"
)

func NewAdminClient(client admin.AdminClient) *AdminClient {
	return &AdminClient{
		client: client,
	}
}

type AdminClient struct {
	client admin.AdminClient
}

func (o *AdminClient) GetConfig(ctx context.Context) (map[string]string, error) {
	conf, err := o.client.GetClientConfig(ctx, &admin.GetClientConfigReq{})
	if err != nil {
		return nil, err
	}
	if conf.Config == nil {
		return map[string]string{}, nil
	}
	return conf.Config, nil
}

func (o *AdminClient) CheckInvitationCode(ctx context.Context, invitationCode string) error {
	resp, err := o.client.FindInvitationCode(ctx, &admin.FindInvitationCodeReq{Codes: []string{invitationCode}})
	if err != nil {
		return err
	}
	if len(resp.Codes) == 0 {
		return eerrs.ErrInvitationNotFound.Wrap()
	}
	if resp.Codes[0].UsedUserID != "" {
		return eerrs.ErrInvitationCodeUsed.Wrap()
	}
	return nil
}

func (o *AdminClient) CheckRegister(ctx context.Context, ip string) error {
	_, err := o.client.CheckRegisterForbidden(ctx, &admin.CheckRegisterForbiddenReq{Ip: ip})
	return err
}

func (o *AdminClient) CheckLogin(ctx context.Context, userID string, ip string) error {
	_, err := o.client.CheckLoginForbidden(ctx, &admin.CheckLoginForbiddenReq{Ip: ip, UserID: userID})
	return err
}

func (o *AdminClient) UseInvitationCode(ctx context.Context, userID string, invitationCode string) error {
	_, err := o.client.UseInvitationCode(ctx, &admin.UseInvitationCodeReq{UserID: userID, Code: invitationCode})
	return err
}

func (o *AdminClient) CheckNilOrAdmin(ctx context.Context) (bool, error) {
	if !mctx.HaveOpUser(ctx) {
		return false, nil
	}
	_, err := mctx.CheckAdmin(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (o *AdminClient) CreateToken(ctx context.Context, userID string, userType int32) (*admin.CreateTokenResp, error) {
	return o.client.CreateToken(ctx, &admin.CreateTokenReq{UserID: userID, UserType: userType})
}

func (o *AdminClient) GetDefaultFriendUserID(ctx context.Context) ([]string, error) {
	resp, err := o.client.FindDefaultFriend(ctx, &admin.FindDefaultFriendReq{})
	if err != nil {
		return nil, err
	}
	return resp.UserIDs, nil
}

func (o *AdminClient) GetDefaultGroupID(ctx context.Context) ([]string, error) {
	resp, err := o.client.FindDefaultGroup(ctx, &admin.FindDefaultGroupReq{})
	if err != nil {
		return nil, err
	}
	return resp.GroupIDs, nil
}
