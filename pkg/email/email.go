package email

import (
	"context"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

//SendEmailRequest represents an email message message
type SendEmailRequest struct {
	Subject   string
	Body      string
	Recipient string
}

//SendEmailResponse represents a response from an email provider
type SendEmailResponse struct {
	ID       string
	Response string
}

//EmailProviderInterface is a thing that sends emails
type EmailProviderInterface interface {
	Send(context.Context, *SendEmailRequest) (*SendEmailResponse, error)
}

//MailGun represents a MailGun EmailProvider
type MailGun struct {
	Domain string
	Sender string
	APIKey string
}

//Send an email based on FIXME: we should have a retry mechanism for Send failures
func (m *MailGun) Send(ctx context.Context, request *SendEmailRequest) (*SendEmailResponse, error) {
	mg := mailgun.NewMailgun(m.Domain, m.APIKey)
	//FIXME: this should be configurable?
	mg.SetAPIBase(mailgun.APIBaseEU)

	message := mg.NewMessage(
		m.Sender,
		request.Subject,
		request.Body,
		request.Recipient,
	)

	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	_, id, err := mg.Send(ctx, message)
	return &SendEmailResponse{ID: id}, err
}
