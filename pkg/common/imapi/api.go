

package imapi

import (
	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/protocol/friend"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/user"
)

// im caller.
var (
	importFriend      = NewApiCaller[friend.ImportFriendReq, friend.ImportFriendResp]("/friend/import_friend")
	userToken         = NewApiCaller[auth.UserTokenReq, auth.UserTokenResp]("/auth/user_token")
	inviteToGroup     = NewApiCaller[group.InviteUserToGroupReq, group.InviteUserToGroupResp]("/group/invite_user_to_group")
	updateUserInfo    = NewApiCaller[user.UpdateUserInfoReq, user.UpdateUserInfoResp]("/user/update_user_info")
	registerUser      = NewApiCaller[user.UserRegisterReq, user.UserRegisterResp]("/user/user_register")
	forceOffLine      = NewApiCaller[auth.ForceLogoutReq, auth.ForceLogoutResp]("/auth/force_logout")
	getGroupsInfo     = NewApiCaller[group.GetGroupsInfoReq, group.GetGroupsInfoResp]("/group/get_groups_info")
	registerUserCount = NewApiCaller[user.UserRegisterCountReq, user.UserRegisterCountResp]("/statistics/user/register")
	friendUserIDs     = NewApiCaller[friend.GetFriendIDsReq, friend.GetFriendIDsResp]("/friend/get_friend_id")
)
