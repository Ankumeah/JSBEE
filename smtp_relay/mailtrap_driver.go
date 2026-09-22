package main

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-go"
)

type MailtrapClient struct{ client *mailtrap.Client }

func GetMailtrapClient(apiKey string) (MailtrapClient, error) {
	client, err := mailtrap.NewClient(
		apiKey,
		mailtrap.WithSandbox(true),
		mailtrap.WithSandboxID(4919444),
	)

	return MailtrapClient{client}, err
}

func (c *MailtrapClient) BatchSend(
	ctx context.Context,
	request BatchSendRequest,
) error {
	var requests []mailtrap.SendRequest
	for _, r := range request.Requests {
		requests = append(requests, processMailtrapSend(r))
	}
	base := processMailtrapSend(request.Base)

	if resp, _, err := c.client.SendBatch(ctx,
		&mailtrap.BatchSendRequest{
			Base:     &base,
			Requests: requests,
		},
	); err != nil {
		return err
	} else if !resp.Success {
		return fmt.Errorf("Errors while sending btach emails: %v", err)
	}

	return nil
}

func (c *MailtrapClient) Send(
	ctx context.Context,
	request SendRequest,
) error {
	r := processMailtrapSend(request)
	_, _, err := c.client.Send(ctx, &r)

	return err
}

func processMailtrapAttachments(
	attachments []Attachment,
) []mailtrap.Attachment {
	var att []mailtrap.Attachment
	for _, attachment := range attachments {
		disposition := mailtrap.DispositionAttachment
		if attachment.Inline {
			disposition = mailtrap.DispositionInline
		}

		att = append(att, mailtrap.Attachment{
			Filename:    attachment.Filename,
			Type:        attachment.Type,
			Disposition: disposition,
			ContentID:   attachment.ContentID,
			Content:     attachment.Data,
		})
	}

	return att
}

func processMailtrapSend(
	request SendRequest,
) mailtrap.SendRequest {
	return mailtrap.SendRequest{
		From: mailtrap.Address{
			Name:  envVars["FROM_NAME"],
			Email: envVars["FROM_EMAIL"],
		},
		To: []mailtrap.Address{
			mailtrap.Address{
				Email: request.To,
			},
		},
		Subject: request.Subject,
		Text:    request.Text,
		HTML:    request.HTML,
		Attachments: processMailtrapAttachments(
			request.Attachments,
		),
	}
}
