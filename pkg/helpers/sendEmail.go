package helpers

import (
	"context"
	"log"

	"go_template/global"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sestypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

// SendEmail sends a plain-text email via AWS SES.
func SendEmail(email string, body string) {
	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(global.Cfg.SES.Sender),
		Destination: &sestypes.Destination{
			ToAddresses: []string{email},
		},
		Content: &sestypes.EmailContent{
			Simple: &sestypes.Message{
				Subject: &sestypes.Content{
					Data: aws.String("Notification"),
				},
				Body: &sestypes.Body{
					Text: &sestypes.Content{
						Data: aws.String(body),
					},
				},
			},
		},
	}

	if _, err := global.SES.SendEmail(context.Background(), input); err != nil {
		log.Printf("SES send email error to %s: %v", email, err)
	}
}
