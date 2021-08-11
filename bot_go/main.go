package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack/socketmode"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func main() {
	webApi := slack.New(
		"xoxb-2363568815571-2368123169266-ta18hwKdJxfxaXOlcY2zA0fr",
		slack.OptionAppLevelToken("xapp-1-A02AR2NG99B-2366685455557-99fa21aad96f4e2206edd3e1fbdadbb6a0f81b3b8bcb3cf07bf7b5392e3c1ec9"),
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "api: ", log.Lshortfile|log.LstdFlags)),
	)
	socketMode := socketmode.New(
		webApi,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "sm: ", log.Lshortfile|log.LstdFlags)),
	)
	authTest, authTestErr := webApi.AuthTest()
	if authTestErr != nil {
		fmt.Fprintf(os.Stderr, "SLACK_BOT_TOKEN is invalid: %v\n", authTestErr)
		os.Exit(1)
	}
	selfUserId := authTest.UserID

	go func() {
		for envelope := range socketMode.Events {
			switch envelope.Type {
			case socketmode.EventTypeEventsAPI:
				// イベント API のハンドリング

				// 3 秒以内にとりあえず ack
				socketMode.Ack(*envelope.Request)

				eventPayload, _ := envelope.Data.(slackevents.EventsAPIEvent)
				switch eventPayload.Type {
				case slackevents.CallbackEvent:
					switch event := eventPayload.InnerEvent.Data.(type) {
					case *slackevents.MessageEvent:
						if event.User != selfUserId && strings.Contains(event.Text, "さよなら") {
							_, _, err := webApi.PostMessage(
								event.Channel,
								slack.MsgOptionText(
									fmt.Sprintf(":wave: さよなら <@%v> さん！", event.User),
									false,
								),
							)
							if err != nil {
								log.Printf("Failed to reply: %v", err)
							}
						}
					default:
						socketMode.Debugf("Skipped: %v", event)
					}
				default:
					socketMode.Debugf("unsupported Events API eventPayload received")
				}
			case socketmode.EventTypeInteractive:
				// ショートカットのハンドリングとモーダル起動
				payload, _ := envelope.Data.(slack.InteractionCallback)
				switch payload.Type {
				case slack.InteractionTypeShortcut:
					if payload.CallbackID == "socket-mode-shortcut" {
						socketMode.Ack(*envelope.Request)
						modalView := slack.ModalViewRequest{
							Type:       "modal",
							CallbackID: "modal-id",
							Title: slack.NewTextBlockObject(
								"plain_text",
								"タスク登録",
								false,
								false,
							),
							Submit: slack.NewTextBlockObject(
								"plain_text",
								"送信",
								false,
								false,
							),
							Close: slack.NewTextBlockObject(
								"plain_text",
								"キャンセル",
								false,
								false,
							),
							Blocks: slack.Blocks{
								BlockSet: []slack.Block{
									slack.NewInputBlock(
										"input-task",
										slack.NewTextBlockObject(
											"plain_text",
											"タスク",
											false,
											false,
										),
										// multiline is not yet supported
										slack.NewPlainTextInputBlockElement(
											slack.NewTextBlockObject(
												"plain_text",
												"タスクの詳細・期限などを書いてください",
												false,
												false,
											),
											"input",
										),
									),
								},
							},
						}
						resp, err := webApi.OpenView(payload.TriggerID, modalView)
						if err != nil {
							log.Printf("Failed to opemn a modal: %v", err)
						}
						socketMode.Debugf("views.open response: %v", resp)
					}
				case slack.InteractionTypeViewSubmission:
					// モーダルからの送信をハンドリング
					if payload.CallbackID == "modal-id" {
						socketMode.Debugf("Submitted Data: %v", payload.View.State.Values)
						socketMode.Ack(*envelope.Request)
					}
				default:
					socketMode.Debugf("Skipped: %v", payload)
				}

			default:
				socketMode.Debugf("Skipped: %v", envelope.Type)
			}
		}
	}()

	socketMode.Run()
}
