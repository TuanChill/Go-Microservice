package pkg

import (
	"bytes"
	"context"
	"html/template"
	"log"

	"go_template/global"
	"go_template/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sestypes "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

func SendGoEmail(email string, data models.EmailData) {
	tmpl, err := template.New("email").Parse(data.Template)
	if err != nil {
		log.Printf("SES email template parse error: %v", err)
		return
	}

	var body bytes.Buffer
	if err = tmpl.Execute(&body, data); err != nil {
		log.Printf("SES email template execute error: %v", err)
		return
	}

	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(global.Cfg.SES.Sender),
		Destination: &sestypes.Destination{
			ToAddresses: []string{email},
		},
		Content: &sestypes.EmailContent{
			Simple: &sestypes.Message{
				Subject: &sestypes.Content{
					Data: aws.String(data.Title),
				},
				Body: &sestypes.Body{
					Html: &sestypes.Content{
						Data: aws.String(body.String()),
					},
				},
			},
		},
	}

	if _, err = global.SES.SendEmail(context.Background(), input); err != nil {
		log.Printf("SES send email error to %s: %v", email, err)
	}
}
