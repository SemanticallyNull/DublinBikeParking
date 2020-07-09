package slack

import (
	"fmt"

	"code.katiechapman.ie/dublinbikeparking/stand"

	"github.com/slack-go/slack"
)

type SlackIntegration struct {
	webhookURL string
}

type approvalMessage struct {
	Coordinates   string
	Location      string
	ID            string
	NoStands      string
	ApprovalValue string
}

func NewSlackIntegration(webhookURL string) *SlackIntegration {
	return &SlackIntegration{
		webhookURL: webhookURL,
	}
}

func (s *SlackIntegration) PostNotification(stand stand.Stand) error {
	approveButton := slack.NewButtonBlockElement("approve", fmt.Sprintf("id=%s&token=%s", stand.StandID, "token"), slack.NewTextBlockObject("plain_text", "Approve", false, false))
	denyButton := slack.NewButtonBlockElement("deny", fmt.Sprintf("id=%s&token=%s", stand.StandID, "token"), slack.NewTextBlockObject("plain_text", "Deny", false, false))
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

//func main() {
//	si := SlackIntegration{
//		webhookURL: "https://hooks.slack.com/services/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",
//	}
//	err := si.PostNotification(apiv0.Stand{
//		StandID:        "90d1f8",
//		Lat:            53.34953,
//		Lng:            -6.24849,
//		Name:           "Georges dock outside Ely chq",
//		Type:           "Sheffield Stand",
//		NumberOfStands: 4,
//	})
//	if err != nil {
//		panic(err)
//	}
//
//	r := gin.Default()
//	r.POST("/handle-slack", func(c *gin.Context) {
//		payload := c.PostForm("payload")
//		interaction := &SlackInteraction{}
//		err := json.Unmarshal([]byte(payload), interaction)
//		if err != nil {
//			fmt.Printf("Could not parse payload %s: %s", payload, err)
//			c.JSON(http.StatusBadRequest, err)
//			return
//		}
//
//		standID := (&url.URL{
//			RawQuery: interaction.Actions[0].Value,
//		}).Query().Get("id")
//
//		br := &bytes.Buffer{}
//		err = json.NewEncoder(br).Encode(struct {
//			ReplaceOriginal string `json:"replace_original"`
//			Text            string `json:"text"`
//		}{
//			ReplaceOriginal: "true",
//			Text:            fmt.Sprintf("<@%s> has %s stand ID %s", interaction.User.ID, interaction.Actions[0].ActionID, standID),
//		})
//		if err != nil {
//			panic(err)
//		}
//		_, err = http.Post(interaction.ResponseURL, "application/json", br)
//		if err != nil {
//			panic(err)
//		}
//
//		c.JSON(200, "OK")
//	})
//	err = http.ListenAndServe(":3488", r)
//	if err != nil {
//		panic(err)
//	}
//}
