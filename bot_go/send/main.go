package main

import "github.com/slack-go/slack"

func main() {
	api := slack.New("xoxb-1491378823348-2286618846368-uIxUrd1gGmuRC2QBPe1euDjr")
	_, _, _ = api.PostMessage(
		"bot-sandbox-nomiyama",
		slack.MsgOptionText("情報セキュリティはやりましたか？モデリングもやりましたか？　のみ", false),
	)
}

//hogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehogehoge
