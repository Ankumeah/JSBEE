package emails_test

import (
	"github.com/Ankumeah/JSBEE/smtp_relay/emails"

	"testing"
)

func TestInitTemplates(t *testing.T) {
	if err := emails.InitTemplates(); err != nil {
		t.Fatal(err.Error())
	}
}
