package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/c-bec-k/tristatcompanion/internal/data"
)

func (bot *application) ReplyDamage(w http.ResponseWriter, opts map[string]interface{}, interaction data.Interaction) {

	var dm, acv int
	level := int(opts["level"].(float64))
	if _, ok := opts["dm"]; ok {
		dm = int(opts["dm"].(float64))
	} else {
		dm = 5
	}
	if _, ok := opts["acv"]; ok {
		acv = int(opts["acv"].(float64))
	}

	totalDamage := (level * dm) + acv

	fmt.Printf("Damage of: %d\nDM of: %d\nACV of: %d", totalDamage, dm, acv)

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
		Value:  fmt.Sprintf("(%d × %d) %+d \n([level × DM] + ACV)", level, dm, acv),
		Inline: false,
	}

	// modsField := data.EmbedField{"Modifiers", fmt.Sprintf("%+d%+d%+d (Stat+Attribute+Misc)", mods[0], mods[1], mods[2]), true}

	fields = append(fields, calcField)

	embed.Fields = fields

	// if edOb != "" {
	// 	moreDiceField := data.EmbedField{"Edges/Obstacles", edOb, true}
	// 	fields = append(fields, moreDiceField)
	// }

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
