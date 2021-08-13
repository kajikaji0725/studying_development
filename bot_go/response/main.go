package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/slack-go/slack/socketmode"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func main() {
	resp, err := http.Get("https://confrage.jp/go-%E8%A8%80%E8%AA%9E%E3%81%AEgo-func-%E3%81%A8channel%E3%81%A8%E3%81%AF/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

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
									fmt.Sprintln(string(b)),
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
			default:
				socketMode.Debugf("Skipped: %v", envelope.Type)
			}
		}
	}()

	er := socketMode.Run()
	if er != nil {
		log.Fatal(err)
	}
}
