

package imapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/openimsdk/chat/pkg/common/constant"
	pconstant "github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"gorm.io/gorm/utils"
)

type baseApiResponse[T any] struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
	Data    *T     `json:"data"`
}

var client = &http.Client{
	Timeout: time.Second * 10,
}

type ApiCaller[Req, Resp any] interface {
	Call(ctx context.Context, apiPrefix string, req *Req) (*Resp, error)
}

func NewApiCaller[Req, Resp any](api string) ApiCaller[Req, Resp] {
	return &caller[Req, Resp]{
		api: api,
	}
}

type caller[Req, Resp any] struct {
	api string
}

func (a caller[Req, Resp]) Call(ctx context.Context, apiPrefix string, req *Req) (*Resp, error) {
	start := time.Now()
	resp, err := a.call(ctx, apiPrefix, req)
	if err != nil {
		log.ZError(ctx, "api caller failed", err, "api", a.api, "duration", time.Since(start), "req", req)
		return nil, err
	}
	log.ZInfo(ctx, "api caller success resp", "api", a.api, "duration", time.Since(start), "req", req, "resp", resp)
	return resp, nil
}

func (a caller[Req, Resp]) call(ctx context.Context, apiPrefix string, req *Req) (*Resp, error) {
	url := apiPrefix + a.api
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	operationID := utils.ToString(ctx.Value(pconstant.OperationID))
	request.Header.Set(pconstant.OperationID, operationID)
	if token, _ := ctx.Value(constant.CtxApiToken).(string); token != "" {
		request.Header.Set(pconstant.Token, token)
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errs.WrapMsg(err, "read http response body", "url", url, "code", response.StatusCode)
	}
	var resp baseApiResponse[Resp]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, errs.WrapMsg(err, string(data))
	}
	if resp.ErrCode != 0 {
		return nil, errs.NewCodeError(resp.ErrCode, resp.ErrMsg).WithDetail(resp.ErrDlt).Wrap()
	}
	return resp.Data, nil
}
