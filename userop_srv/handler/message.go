package handler

import (
	"context"

	"mall_srvs/userop_srv/global"
	"mall_srvs/userop_srv/model"
	"mall_srvs/userop_srv/proto"
)

func (userOpServer *UserOpServer) MessageList(ctx context.Context, req *proto.MessageRequest) (*proto.MessageListResponse, error) {
	var rsp proto.MessageListResponse
	var messages []model.LeavingMessages
	var messageList []*proto.MessageResponse

	result := global.DB.Where(&model.LeavingMessages{
		User: req.GetUserId(),
	}).Find(&messages)
	rsp.Total = int32(result.RowsAffected)

	for _, message := range messages {
		messageList = append(messageList, &proto.MessageResponse{
			Id:          message.ID,
			UserId:      message.User,
			MessageType: message.MessageType,
			Subject:     message.Subject,
			Message:     message.Message,
			File:        message.File,
		})
	}

	rsp.Data = messageList
	return &rsp, nil
}

func (userOpServer *UserOpServer) CreateMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {
	var message model.LeavingMessages

	message.User = req.GetUserId()
	message.MessageType = req.GetMessageType()
	message.Subject = req.GetSubject()
	message.Message = req.GetMessage()
	message.File = req.GetFile()

	global.DB.Save(&message)

	return &proto.MessageResponse{
		Id: message.ID,
	}, nil
}
