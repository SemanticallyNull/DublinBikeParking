package slack

import (
	"fmt"

	"code.katiechapman.ie/dublinbikeparking/stand"

	"github.com/slack-go/slack"
)

type SlackIntegration struct {
	webhookURL string
}

func NewSlackIntegration(webhookURL string) *SlackIntegration {
	return &SlackIntegration{
		webhookURL: webhookURL,
	}
}

func (s *SlackIntegration) PostNotification(stand stand.Stand) error {
	approveButton := slack.NewButtonBlockElement("approve", fmt.Sprintf("id=%s&token=%s", stand.StandID, stand.Token), slack.NewTextBlockObject("plain_text", "Approve", false, false))
	denyButton := slack.NewButtonBlockElement("deny", fmt.Sprintf("id=%s&token=%s", stand.StandID, stand.Token), slack.NewTextBlockObject("plain_text", "Deny", false, false))
	approveButton.Style = "primary"
	denyButton.Style = "danger"

	bms := slack.NewSectionBlock(
		slack.NewTextBlockObject(
			"mrkdwn",
			fmt.Sprintf("A new bike stand has been submitted:\n*<https://dublinbikeparking.com/update.html#18/%f/%f|%s>*", stand.Lat, stand.Lng, stand.Name),
			false,
			false,
		),
		[]*slack.TextBlockObject{
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*ID:*\n%s", stand.StandID), false, false),
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Number Of Stands:*\n%d", stand.NumberOfStands), false, false),
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Location:*\n%s", stand.Name), false, false),
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Coordinates:*\n%f/%f", stand.Lat, stand.Lng), false, false),
		},
		nil,
	)
	bma := slack.NewActionBlock(
		"approval_actions",
		approveButton,
		denyButton,
	)

	if err := slack.PostWebhook(s.webhookURL, &slack.WebhookMessage{
		Blocks: &slack.Blocks{
			BlockSet: []slack.Block{bms, bma},
		},
	}); err != nil {
		return err
	}

	return nil
}

type SlackInteraction struct {
	User struct {
		ID string `json:"id"`
	} `json:"user"`
	ResponseURL string `json:"response_url"`
	Actions     []struct {
		ActionID string `json:"action_id"`
		Value    string `json:"value"`
	} `json:"actions"`
}
