package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/koraygocmen/golang-boilerplate/internal/config"
	"github.com/koraygocmen/golang-boilerplate/internal/context"
	"github.com/koraygocmen/golang-boilerplate/internal/env"
	slackgo "github.com/slack-go/slack"
)

type messageErrorFn func(ctx context.Ctx, err error) error
type messageEventFn func(ctx context.Ctx, title, msg string) error

type client struct {
	MessageError messageErrorFn
	MessageEvent messageEventFn
}

var (
	httpClient = &http.Client{}
	Client     = client{
		MessageError: messageError(),
		MessageEvent: messageEvent(),
	}
)

func message(webhook string, message *slackgo.WebhookMessage) error {
	if env.IsDev() || webhook == "" {
		return nil
	}

	reqbody, err := json.Marshal(message)
	if err != nil {
		err := fmt.Errorf("slack message error: marshal message error: %w", err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, webhook, bytes.NewReader(reqbody))
	if err != nil {
		err := fmt.Errorf("slack message error: http new request error: %w", err)
		return err
	}
	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	res, err := httpClient.Do(req)
	if err != nil {
		err := fmt.Errorf("slack message error: http client do error: %w", err)
		return err
	}

	if res.StatusCode != http.StatusOK {
		resbody, err := io.ReadAll(res.Body)
		if err != nil {
			err := fmt.Errorf("slack message error: read response body error: %w", err)
			return err
		}

		err = fmt.Errorf("slack message error: status code: %v, body: %v", res.StatusCode, string(resbody))
		return err
	}

	return nil
}

func messageError() messageErrorFn {
	return func(ctx context.Ctx, err error) error {
		var contextFields []*slackgo.TextBlockObject
		for k, v := range context.Map(ctx) {
			contextFields = append(contextFields, &slackgo.TextBlockObject{
				Type: slackgo.PlainTextType,
				Text: fmt.Sprintf("%v: %v", k, v),
			})
		}

		m := &slackgo.WebhookMessage{
			Blocks: &slackgo.Blocks{
				BlockSet: []slackgo.Block{
					&slackgo.SectionBlock{
						Type: slackgo.MBTSection,
						Text: &slackgo.TextBlockObject{
							Type: slackgo.MarkdownType,
							Text: "<!channel|channel>",
						},
					},
					&slackgo.SectionBlock{
						Type: slackgo.MBTHeader,
						Text: &slackgo.TextBlockObject{
							Type:  slackgo.PlainTextType,
							Text:  "Application Error",
							Emoji: true,
						},
					},
					&slackgo.SectionBlock{
						Type: slackgo.MBTDivider,
					},
					&slackgo.SectionBlock{
						Type:   slackgo.MBTSection,
						Fields: contextFields,
					},
					&slackgo.SectionBlock{
						Type: slackgo.MBTDivider,
					},
					&slackgo.SectionBlock{
						Type: slackgo.MBTSection,
						Text: &slackgo.TextBlockObject{
							Type: slackgo.MarkdownType,
							Text: fmt.Sprintf("```%v```", err),
						},
					},
				},
			},
		}

		return message(config.Slack.Webhook.Errors, m)
	}
}

func messageEvent() messageEventFn {
	return func(ctx context.Ctx, title, msg string) error {
		m := &slackgo.WebhookMessage{
			Blocks: &slackgo.Blocks{
				BlockSet: []slackgo.Block{
					&slackgo.SectionBlock{
						Type: slackgo.MBTSection,
						Text: &slackgo.TextBlockObject{
							Type: slackgo.MarkdownType,
							Text: "<!channel|channel>",
						},
					},
					&slackgo.SectionBlock{
						Type: slackgo.MBTHeader,
						Text: &slackgo.TextBlockObject{
							Type:  slackgo.PlainTextType,
							Text:  title,
							Emoji: true,
						},
					},
					&slackgo.SectionBlock{
						Type: slackgo.MBTDivider,
					},
					&slackgo.SectionBlock{
						Type: slackgo.MBTSection,
						Text: &slackgo.TextBlockObject{
							Type: slackgo.MarkdownType,
							Text: msg,
						},
					},
				},
			},
		}

		return message(config.Slack.Webhook.Events, m)
	}
}
