

package database

import (
	"context"
	"github.com/openimsdk/chat/pkg/common/db/cache"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/db/tx"
	"github.com/redis/go-redis/v9"

	"github.com/openimsdk/chat/pkg/common/db/model/admin"
	table "github.com/openimsdk/chat/pkg/common/db/table/admin"
)

type AdminDatabaseInterface interface {
	GetAdmin(ctx context.Context, account string) (*table.Admin, error)
	GetAdminUserID(ctx context.Context, userID string) (*table.Admin, error)
	UpdateAdmin(ctx context.Context, userID string, update map[string]any) error
	ChangePassword(ctx context.Context, userID string, newPassword string) error
	AddAdminAccount(ctx context.Context, admin []*table.Admin) error
	DelAdminAccount(ctx context.Context, userIDs []string) error
	SearchAdminAccount(ctx context.Context, pagination pagination.Pagination) (int64, []*table.Admin, error)
	CreateApplet(ctx context.Context, applets []*table.Applet) error
	DelApplet(ctx context.Context, appletIDs []string) error
	GetApplet(ctx context.Context, appletID string) (*table.Applet, error)
	FindApplet(ctx context.Context, appletIDs []string) ([]*table.Applet, error)
	SearchApplet(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*table.Applet, error)
	FindOnShelf(ctx context.Context) ([]*table.Applet, error)
	UpdateApplet(ctx context.Context, appletID string, update map[string]any) error
	GetConfig(ctx context.Context) (map[string]string, error)
	SetConfig(ctx context.Context, cs map[string]string) error
	DelConfig(ctx context.Context, keys []string) error
	FindInvitationRegister(ctx context.Context, codes []string) ([]*table.InvitationRegister, error)
	DelInvitationRegister(ctx context.Context, codes []string) error
	UpdateInvitationRegister(ctx context.Context, code string, fields map[string]any) error
	CreatInvitationRegister(ctx context.Context, invitationRegisters []*table.InvitationRegister) error
	SearchInvitationRegister(ctx context.Context, keyword string, state int32, userIDs []string, codes []string, pagination pagination.Pagination) (int64, []*table.InvitationRegister, error)
	SearchIPForbidden(ctx context.Context, keyword string, state int32, pagination pagination.Pagination) (int64, []*table.IPForbidden, error)
	AddIPForbidden(ctx context.Context, ms []*table.IPForbidden) error
	FindIPForbidden(ctx context.Context, ms []string) ([]*table.IPForbidden, error)
	DelIPForbidden(ctx context.Context, ips []string) error
	FindDefaultFriend(ctx context.Context, userIDs []string) ([]string, error)
	AddDefaultFriend(ctx context.Context, ms []*table.RegisterAddFriend) error
	DelDefaultFriend(ctx context.Context, userIDs []string) error
	SearchDefaultFriend(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*table.RegisterAddFriend, error)
	FindDefaultGroup(ctx context.Context, groupIDs []string) ([]string, error)
	AddDefaultGroup(ctx context.Context, ms []*table.RegisterAddGroup) error
	DelDefaultGroup(ctx context.Context, groupIDs []string) error
	SearchDefaultGroup(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*table.RegisterAddGroup, error)
	FindBlockInfo(ctx context.Context, userIDs []string) ([]*table.ForbiddenAccount, error)
	GetBlockInfo(ctx context.Context, userID string) (*table.ForbiddenAccount, error)
	BlockUser(ctx context.Context, f []*table.ForbiddenAccount) error
	DelBlockUser(ctx context.Context, userID []string) error
	SearchBlockUser(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*table.ForbiddenAccount, error)
	FindBlockUser(ctx context.Context, userIDs []string) ([]*table.ForbiddenAccount, error)
	SearchUserLimitLogin(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*table.LimitUserLoginIP, error)
	AddUserLimitLogin(ctx context.Context, ms []*table.LimitUserLoginIP) error
	DelUserLimitLogin(ctx context.Context, ms []*table.LimitUserLoginIP) error
	CountLimitUserLoginIP(ctx context.Context, userID string) (uint32, error)
	GetLimitUserLoginIP(ctx context.Context, userID string, ip string) (*table.LimitUserLoginIP, error)
	CacheToken(ctx context.Context, userID string, token string) error
	GetTokens(ctx context.Context, userID string) (map[string]int32, error)
}

func NewAdminDatabase(cli *mongoutil.Client, rdb redis.UniversalClient) (AdminDatabaseInterface, error) {
	a, err := admin.NewAdmin(cli.GetDB())
	if err != nil {
		return nil, err
	}
	forbidden, err := admin.NewIPForbidden(cli.GetDB())
	if err != nil {
		return nil, err
	}
	forbiddenAccount, err := admin.NewForbiddenAccount(cli.GetDB())
	if err != nil {
		return nil, err
	}
	limitUserLoginIP, err := admin.NewLimitUserLoginIP(cli.GetDB())
	if err != nil {
		return nil, err
	}
	invitationRegister, err := admin.NewInvitationRegister(cli.GetDB())
	if err != nil {
		return nil, err
	}
	registerAddFriend, err := admin.NewRegisterAddFriend(cli.GetDB())
	if err != nil {
		return nil, err
	}
	registerAddGroup, err := admin.NewRegisterAddGroup(cli.GetDB())
	if err != nil {
		return nil, err
	}
	applet, err := admin.NewApplet(cli.GetDB())
	if err != nil {
		return nil, err
	}
	clientConfig, err := admin.NewClientConfig(cli.GetDB())
	if err != nil {
		return nil, err
	}
	return &AdminDatabase{
		tx:                 cli.GetTx(),
		admin:              a,
		ipForbidden:        forbidden,
		forbiddenAccount:   forbiddenAccount,
		limitUserLoginIP:   limitUserLoginIP,
		invitationRegister: invitationRegister,
		registerAddFriend:  registerAddFriend,
		registerAddGroup:   registerAddGroup,
		applet:             applet,
		clientConfig:       clientConfig,
		cache:              cache.NewTokenInterface(rdb),
	}, nil
}

type AdminDatabase struct {
	tx                 tx.Tx
	admin              table.AdminInterface
	ipForbidden        table.IPForbiddenInterface
	forbiddenAccount   table.ForbiddenAccountInterface
	limitUserLoginIP   table.LimitUserLoginIPInterface
	invitationRegister table.InvitationRegisterInterface
	registerAddFriend  table.RegisterAddFriendInterface
	registerAddGroup   table.RegisterAddGroupInterface
	applet             table.AppletInterface
	clientConfig       table.ClientConfigInterface
	cache              cache.TokenInterface
}

func (o *AdminDatabase) GetAdmin(ctx context.Context, account string) (*table.Admin, error) {
	return o.admin.Take(ctx, account)
}

func (o *AdminDatabase) GetAdminUserID(ctx context.Context, userID string) (*table.Admin, error) {
	return o.admin.TakeUserID(ctx, userID)
}

func (o *AdminDatabase) UpdateAdmin(ctx context.Context, userID string, update map[string]any) error {
	return o.admin.Update(ctx, userID, update)
}

func (o *AdminDatabase) ChangePassword(ctx context.Context, userID string, newPassword string) error {
	return o.admin.ChangePassword(ctx, userID, newPassword)
}
func (o *AdminDatabase) AddAdminAccount(ctx context.Context, admins []*table.Admin) error {
	return o.admin.Create(ctx, admins)
}

func (o *AdminDatabase) DelAdminAccount(ctx context.Context, userIDs []string) error {
	return o.admin.Delete(ctx, userIDs)
}

func (o *AdminDatabase) SearchAdminAccount(ctx context.Context, pagination pagination.Pagination) (int64, []*table.Admin, error) {
	return o.admin.Search(ctx, pagination)
}

func (o *AdminDatabase) CreateApplet(ctx context.Context, applets []*table.Applet) error {
	return o.applet.Create(ctx, applets)
}

func (o *AdminDatabase) DelApplet(ctx context.Context, appletIDs []string) error {
	return o.applet.Del(ctx, appletIDs)
}

func (o *AdminDatabase) GetApplet(ctx context.Context, appletID string) (*table.Applet, error) {
	return o.applet.Take(ctx, appletID)
}

func (o *AdminDatabase) FindApplet(ctx context.Context, appletIDs []string) ([]*table.Applet, error) {
	return o.applet.FindID(ctx, appletIDs)
}

func (o *AdminDatabase) SearchApplet(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*table.Applet, error) {
	return o.applet.Search(ctx, keyword, pagination)
}

func (o *AdminDatabase) FindOnShelf(ctx context.Context) ([]*table.Applet, error) {
	return o.applet.FindOnShelf(ctx)
}

func (o *AdminDatabase) UpdateApplet(ctx context.Context, appletID string, update map[string]any) error {
	return o.applet.Update(ctx, appletID, update)
}

func (o *AdminDatabase) GetConfig(ctx context.Context) (map[string]string, error) {
	return o.clientConfig.Get(ctx)
}

func (o *AdminDatabase) SetConfig(ctx context.Context, cs map[string]string) error {
	return o.clientConfig.Set(ctx, cs)
}

func (o *AdminDatabase) DelConfig(ctx context.Context, keys []string) error {
	return o.clientConfig.Del(ctx, keys)
}

func (o *AdminDatabase) FindInvitationRegister(ctx context.Context, codes []string) ([]*table.InvitationRegister, error) {
	return o.invitationRegister.Find(ctx, codes)
}

func (o *AdminDatabase) DelInvitationRegister(ctx context.Context, codes []string) error {
	return o.invitationRegister.Del(ctx, codes)
}

func (o *AdminDatabase) UpdateInvitationRegister(ctx context.Context, code string, fields map[string]any) error {
	return o.invitationRegister.Update(ctx, code, fields)
}

func (o *AdminDatabase) CreatInvitationRegister(ctx context.Context, invitationRegisters []*table.InvitationRegister) error {
	return o.invitationRegister.Create(ctx, invitationRegisters)
}

func (o *AdminDatabase) SearchInvitationRegister(ctx context.Context, keyword string, state int32, userIDs []string, codes []string, pagination pagination.Pagination) (int64, []*table.InvitationRegister, error) {
	return o.invitationRegister.Search(ctx, keyword, state, userIDs, codes, pagination)
}

func (o *AdminDatabase) SearchIPForbidden(ctx context.Context, keyword string, state int32, pagination pagination.Pagination) (int64, []*table.IPForbidden, error) {
	return o.ipForbidden.Search(ctx, keyword, state, pagination)
}

func (o *AdminDatabase) AddIPForbidden(ctx context.Context, ms []*table.IPForbidden) error {
	return o.ipForbidden.Create(ctx, ms)
}

func (o *AdminDatabase) FindIPForbidden(ctx context.Context, ms []string) ([]*table.IPForbidden, error) {
	return o.ipForbidden.Find(ctx, ms)
}

func (o *AdminDatabase) DelIPForbidden(ctx context.Context, ips []string) error {
	return o.ipForbidden.Delete(ctx, ips)
}

func (o *AdminDatabase) FindDefaultFriend(ctx context.Context, userIDs []string) ([]string, error) {
	return o.registerAddFriend.FindUserID(ctx, userIDs)
}

func (o *AdminDatabase) AddDefaultFriend(ctx context.Context, ms []*table.RegisterAddFriend) error {
	return o.registerAddFriend.Add(ctx, ms)
}

func (o *AdminDatabase) DelDefaultFriend(ctx context.Context, userIDs []string) error {
	return o.registerAddFriend.Del(ctx, userIDs)
}

func (o *AdminDatabase) SearchDefaultFriend(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*table.RegisterAddFriend, error) {
	return o.registerAddFriend.Search(ctx, keyword, pagination)
}

func (o *AdminDatabase) FindDefaultGroup(ctx context.Context, groupIDs []string) ([]string, error) {
	return o.registerAddGroup.FindGroupID(ctx, groupIDs)
}

func (o *AdminDatabase) AddDefaultGroup(ctx context.Context, ms []*table.RegisterAddGroup) error {
	return o.registerAddGroup.Add(ctx, ms)
}

func (o *AdminDatabase) DelDefaultGroup(ctx context.Context, groupIDs []string) error {
	return o.registerAddGroup.Del(ctx, groupIDs)
}

func (o *AdminDatabase) SearchDefaultGroup(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*table.RegisterAddGroup, error) {
	return o.registerAddGroup.Search(ctx, keyword, pagination)
}

func (o *AdminDatabase) FindBlockInfo(ctx context.Context, userIDs []string) ([]*table.ForbiddenAccount, error) {
	return o.forbiddenAccount.Find(ctx, userIDs)
}

func (o *AdminDatabase) GetBlockInfo(ctx context.Context, userID string) (*table.ForbiddenAccount, error) {
	return o.forbiddenAccount.Take(ctx, userID)
}

func (o *AdminDatabase) BlockUser(ctx context.Context, f []*table.ForbiddenAccount) error {
	return o.forbiddenAccount.Create(ctx, f)
}

func (o *AdminDatabase) DelBlockUser(ctx context.Context, userID []string) error {
	return o.forbiddenAccount.Delete(ctx, userID)
}

func (o *AdminDatabase) SearchBlockUser(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*table.ForbiddenAccount, error) {
	return o.forbiddenAccount.Search(ctx, keyword, pagination)
}

func (o *AdminDatabase) FindBlockUser(ctx context.Context, userIDs []string) ([]*table.ForbiddenAccount, error) {
	return o.forbiddenAccount.Find(ctx, userIDs)
}

func (o *AdminDatabase) SearchUserLimitLogin(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*table.LimitUserLoginIP, error) {
	return o.limitUserLoginIP.Search(ctx, keyword, pagination)
}

func (o *AdminDatabase) AddUserLimitLogin(ctx context.Context, ms []*table.LimitUserLoginIP) error {
	return o.limitUserLoginIP.Create(ctx, ms)
}

func (o *AdminDatabase) DelUserLimitLogin(ctx context.Context, ms []*table.LimitUserLoginIP) error {
	return o.limitUserLoginIP.Delete(ctx, ms)
}

func (o *AdminDatabase) CountLimitUserLoginIP(ctx context.Context, userID string) (uint32, error) {
	return o.limitUserLoginIP.Count(ctx, userID)
}

func (o *AdminDatabase) GetLimitUserLoginIP(ctx context.Context, userID string, ip string) (*table.LimitUserLoginIP, error) {
	return o.limitUserLoginIP.Take(ctx, userID, ip)
}

func (o *AdminDatabase) CacheToken(ctx context.Context, userID string, token string) error {
	return o.cache.AddTokenFlag(ctx, userID, token, constant.NormalToken)
}

func (o *AdminDatabase) GetTokens(ctx context.Context, userID string) (map[string]int32, error) {
	return o.cache.GetTokensWithoutError(ctx, userID)
}
