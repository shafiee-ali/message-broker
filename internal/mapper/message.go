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

func CreateMessageDTOToPostgresMessage(msg broker.CreateMessageDTO) models.PostgresMessage {
	return models.PostgresMessage{
		Subject:        msg.Subject,
		Body:           msg.Body,
		ExpirationTime: msg.ExpirationTime,
	}
}

func CreateMessageDTOToCassandraMessage(msg broker.CreateMessageDTO) models.CassandraMessage {
	return models.CassandraMessage{
		Subject:        msg.Subject,
		Body:           msg.Body,
		ExpirationTime: msg.ExpirationTime,
	}
}

func PostgresMessageToCreatedMessage(msg models.PostgresMessage) types.CreatedMessage {
	return types.CreatedMessage{
		Id:             msg.ID,
		Subject:        msg.Subject,
		Body:           msg.Body,
		ExpirationTime: msg.ExpirationTime,
	}
}

func CassandraMessageToCreatedMessage(msg models.CassandraMessage) types.CreatedMessage {
	return types.CreatedMessage{
		Id:             msg.ID,
		Subject:        msg.Subject,
		Body:           msg.Body,
		ExpirationTime: msg.ExpirationTime,
	}
}
