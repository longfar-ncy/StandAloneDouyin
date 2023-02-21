// Code generated by hertz generator.

package api

import (
	"context"
	"douyin/biz/model/api"
	"douyin/biz/service"
	"douyin/pkg/constant"
	"douyin/pkg/errno"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// GetFeed .
// @router /douyin/feed/ [GET]
func GetFeed(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DouyinFeedRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	hlog.Infof("handler.feed_service.GetFeed Request: %#v", req)
	userID := c.GetUint64(constant.IdentityKey)
	hlog.Info("handler.feed_service.GetFeed GetUserID:", userID)
	resp, err := service.GetFeed(req.LatestTime, userID)
	if err != nil {
		errNo := errno.ConvertErr(err)
		c.JSON(consts.StatusOK, &api.DouyinFeedResponse{
			StatusCode: errNo.ErrCode,
			StatusMsg:  &errNo.ErrMsg,
		})
		return
	}

	c.JSON(consts.StatusOK, resp)
}
