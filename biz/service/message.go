package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/errno"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func SendMessage(fromUserID, toUserID uint64, content string) (*api.DouyinMessageActionResponse, error) {
	isFriend := db.IsFriend(fromUserID, toUserID)
	if !isFriend {
		errNo := errno.UserRequestParameterError
		errNo.ErrMsg = "不能给非好友发消息"
		hlog.Error("service.message.SendMessage err:", errNo.Error())
		return nil, errNo
	}
	err := db.CreateMessage(fromUserID, toUserID, content)
	if err != nil {
		hlog.Error("service.message.SendMessage err:", err.Error())
		return nil, err
	}
	return &api.DouyinMessageActionResponse{
		StatusCode: errno.Success.ErrCode,
	}, nil
}

func GetMessageChat(userID, oppositeID uint64, preMsgTime int64) (*api.DouyinMessageChatResponse, error) {
	if userID == oppositeID {
		return nil, errno.UserRequestParameterError
	}
	messages, err := db.GetMessagesByUserIDAndPreMsgTime(userID, oppositeID, preMsgTime)
	if err != nil {
		hlog.Error("service.message.GetMessageChat err:", err.Error())
		return nil, err
	}
	return &api.DouyinMessageChatResponse{
		StatusCode:  errno.Success.ErrCode,
		MessageList: pack.Messages(messages),
	}, nil
}
