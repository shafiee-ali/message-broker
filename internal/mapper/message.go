package mapper

import (
	"therealbroker/internal/types"
	"therealbroker/pkg/broker"
	"therealbroker/pkg/models"
)

func InMemoryDBMessageToCreatedMessage(msgModel models.Message) types.CreatedMessage {
	return types.CreatedMessage{
		Id:             msgModel.Id,
		Subject:        msgModel.Subject,
		Body:           msgModel.Body,
		ExpirationTime: msgModel.ExpirationTime,
	}
}

func CreateMessageDTOToDBMessage(msg broker.CreateMessageDTO) models.PostgresMessage {
	return models.PostgresMessage{
		Subject:        msg.Subject,
		Body:           msg.Body,
		ExpirationTime: msg.ExpirationTime,
	}
}

func DBMessageToCreatedMessage(msg models.PostgresMessage) types.CreatedMessage {
	return types.CreatedMessage{
		Id:             msg.ID,
		Subject:        msg.Subject,
		Body:           msg.Body,
		ExpirationTime: msg.ExpirationTime,
	}
}
