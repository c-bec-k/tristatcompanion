package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/c-bec-k/tristatcompanion/internal/data"
)

func (bot *application) ReplyDamage(w http.ResponseWriter, opts map[string]interface{}, interaction data.Interaction) {

	totalDamage := (getLevel(opts) * getDM(opts)) + getAcv(opts)

	embed := data.MessageEmbed{
		Color: 15217272,
	}

	embed.Title = fmt.Sprintf("You deal %d damage!", totalDamage)

	if opts["description"] != nil {
		embed.Author.Name = strconv.Quote(opts["description"].(string))
	}

	var fields []data.EmbedField

	calcField := data.EmbedField{
		Name:   "Damage Calculation",
		Value:  fmt.Sprintf("(%d × %d) %+d \n([level × DM] + ACV)", getLevel(opts), getDM(opts), getAcv(opts)),
		Inline: false,
	}

	fields = append(fields, calcField)

	embed.Fields = fields

	reply := data.InteractionResponse{
		Type: data.ChannelMessageWithSourceCallback,
		Data: data.InteractionCallbackData{
			Embeds: []data.MessageEmbed{embed},
		},
	}

	//fmt.Printf("Reply sent: %+v\n", reply)
	js, err := json.Marshal(reply)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getLevel(opts map[string]interface{}) int {
	return int(opts["level"].(float64))
}

func getDM(opts map[string]interface{}) int {
	if _, ok := opts["dm"]; ok {
		return int(opts["dm"].(float64))
	}
	return 5
}

func getAcv(opts map[string]interface{}) int {
	if _, ok := opts["acv"]; ok {
		return int(opts["acv"].(float64))
	}
	return 0
}
