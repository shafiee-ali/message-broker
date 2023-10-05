package repository

import (
	"therealbroker/internal/types"
	pkgBroker "therealbroker/pkg/broker"
)

type IMessageRepository interface {
	Add(message pkgBroker.CreateMessageDTO) types.CreatedMessage
	FetchUnexpiredBySubjectAndId(subject string, id int) (types.CreatedMessage, error)
}
