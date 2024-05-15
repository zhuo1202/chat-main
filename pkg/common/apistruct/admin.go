

package apistruct

import "github.com/openimsdk/protocol/sdkws"

type AdminLoginResp struct {
	AdminAccount string `json:"adminAccount"`
	AdminToken   string `json:"adminToken"`
	Nickname     string `json:"nickname"`
	FaceURL      string `json:"faceURL"`
	Level        int32  `json:"level"`
	AdminUserID  string `json:"adminUserID"`
	ImUserID     string `json:"imUserID"`
	ImToken      string `json:"imToken"`
}

type SearchDefaultGroupResp struct {
	Total  uint32             `json:"total"`
	Groups []*sdkws.GroupInfo `json:"groups"`
}

type NewUserCountResp struct {
	Total     int64            `json:"total"`
	DateCount map[string]int64 `json:"date_count"`
}
