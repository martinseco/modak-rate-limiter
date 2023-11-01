package sender

import (
	"fmt"

	"github.com/martinseco/modak-rate-limiter/src/api/models"
)

type (
	senderService struct{}
)

func NewSenderService() *senderService {
	return &senderService{}
}

// Send sends notification to the user. In this implementation, it just logs to the console
func (s *senderService) Send(request models.NotificationRequest) error {
	fmt.Println(fmt.Sprintf("Sending %s message to user %s...", request.Type, request.UserID))

	return nil
}
