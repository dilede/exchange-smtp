package exchangesmtp

import "testing"

func TestMailString(t *testing.T) {
	want := "From: gdo@gazdv.ru\n" +
		"To: DemeninDL@gazdv.ru\n" +
		"Subject: Test message\n\n" +
		"With test body."

	m := Mail{
		From:    "gdo@gazdv.ru",
		To:      []string{"DemeninDL@gazdv.ru"},
		Subject: "Test message",
		Body:    "With test body.",
	}
	got := m.String()

	if got != want {
		t.Error("Expected: " + want + " got: " + got)
	}
}
