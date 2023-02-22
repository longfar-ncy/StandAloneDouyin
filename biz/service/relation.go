package service

import (
	"douyin/biz/model/api"
	"douyin/dal/db"
	"douyin/dal/pack"
	"douyin/pkg/errno"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Follow(userID, toUserID uint64) (*api.DouyinRelationActionResponse, error) {
	if userID == toUserID {
		errNo := errno.UserRequestParameterError
		errNo.ErrMsg = "不能自己关注自己哦"
		hlog.Error("service.relation.Follow err:", errNo.Error())
		return nil, errNo
	}
	isFollow := db.IsFollow(userID, toUserID)
	if isFollow {
		hlog.Error("service.relation.Follow err:", errno.RepeatOperationError)
		return nil, errno.RepeatOperationError
	}

	//关注操作
	err := db.Follow(userID, toUserID)
	if err != nil {
		hlog.Error("service.relation.Follow err:", err.Error())
		return nil, err
	}
	return &api.DouyinRelationActionResponse{
		StatusCode: int64(api.ErrCode_Success),
	}, nil
}

func CancelFollow(userID, toUserID uint64) (*api.DouyinRelationActionResponse, error) {
	if userID == toUserID {
		errNo := errno.UserRequestParameterError
		errNo.ErrMsg = "不能自己取关自己哦"
		hlog.Error("service.relation.CancelFollow err:", errNo.Error())
		return nil, errNo
	}
	//取消关注
	err := db.CancelFollow(userID, toUserID)
	if err != nil {
		hlog.Error("service.relation.CancelFollow err:", err.Error())
		return nil, err
	}
	return &api.DouyinRelationActionResponse{
		StatusCode: int64(api.ErrCode_Success),
	}, nil
}

func GetFollowList(userID uint64) (*api.DouyinRelationFollowListResponse, error) {
	dbUserList, err := db.GetFollowList(userID)
	if err != nil {
		hlog.Error("service.relation.GetFollowList err:", err.Error())
		return nil, err
	}

	// 提前申请好数组大小来避免后续扩容
	userList := make([]*api.User, 0, len(dbUserList))
	// TODO 存在循环查询DB
	for _, v := range dbUserList {
		// 这里要查的是，关注列表的人是否关注了自己
		userList = append(userList, pack.User(v, db.IsFollow(v.ID, userID)))
	}

	return &api.DouyinRelationFollowListResponse{
		StatusCode: errno.Success.ErrCode,
		UserList:   userList,
	}, nil

}

func GetFollowerList(userID uint64) (*api.DouyinRelationFollowerListResponse, error) {
	dbUserList, err := db.GetFollowerList(userID)
	if err != nil {
		hlog.Error("service.relation.GetFollowerList err:", err.Error())
		return nil, err
	}

	// 提前申请好数组大小来避免后续扩容
	// TODO 存在循环查询DB
	userList := make([]*api.User, 0, len(dbUserList))
	for _, v := range dbUserList {
		// 这里要查的是，自己是否关注了粉丝列表的人
		userList = append(userList, pack.User(v, db.IsFollow(userID, v.ID)))
	}
	return &api.DouyinRelationFollowerListResponse{
		StatusCode: errno.Success.ErrCode,
		UserList:   userList,
	}, nil
}

func GetFriendList(userID uint64) (*api.DouyinRelationFriendListResponse, error) {
	resultList, err := db.GetFriendList(userID)
	if err != nil {
		hlog.Error("service.relation.GetFollowerList err:", err.Error())
		return nil, err
	}

	return &api.DouyinRelationFriendListResponse{
		StatusCode: errno.Success.ErrCode,
		UserList:   resultList,
	}, nil
}
