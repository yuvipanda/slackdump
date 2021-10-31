package slackdump

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/slack-go/slack"
)

// Messages keeps slice of messages
type Messages struct {
	Messages  []slack.Message
	ChannelID string
	SD        *SlackDumper
}

// ToText outputs Messages m to io.Writer w in Text format
func (m Messages) ToText(w io.Writer) (err error) {
	const minMsgTimeApart = 2 //minutes
	writer := bufio.NewWriter(w)
	defer writer.Flush()

	var prevMsg = slack.Message{}
	var prevTime = time.Time{}
	// var lastMsgFrom string
	for _, message := range m.Messages {
		t, err := fromSlackTime(message.Timestamp)
		if err != nil {
			return err
		}
		diff := t.Sub(prevTime)
		if prevMsg.User == message.User && diff.Minutes() < minMsgTimeApart {
			writer.WriteString(fmt.Sprintf(
				"%s\n", message.Text,
			))
		} else {
			writer.WriteString(fmt.Sprintf(
				"\n> %s @ %s:\n%s\n",
				m.GetUserForMessage(&message),
				t.Format("02/01/2006 15:04:05 Z0700"),
				message.Text,
			))
		}
		prevMsg = message
		prevTime = t

	}
	return nil
}

// GetUserForMessage returns username for the message
func (m *Messages) GetUserForMessage(msg *slack.Message) string {
	var userid string
	if msg.Comment != nil {
		userid = msg.Comment.User
	} else {
		userid = msg.User
	}

	if userid != "" {
		return m.SD.UserForID[userid].Name
	}

	return ""
}