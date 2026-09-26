package main

import (
	"context"
)

type Attachment struct {
	Filename  string
	Type      string
	Inline    bool
	ContentID string
	Data      string
}

type SendRequest struct {
	To          string
	Subject     string
	Text        string
	HTML        string
	Attachments []Attachment
}

type BatchSendRequest struct {
	Base     SendRequest
	Requests []SendRequest
}

type EmailController interface {
	Send(ctx context.Context, request SendRequest) error
	BatchSend(ctx context.Context, request BatchSendRequest) error
}
