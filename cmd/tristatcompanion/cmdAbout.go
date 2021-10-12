package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/c-bec-k/tristatcompanion/internal/data"
)

func (bot *application) ReplyAbout(w http.ResponseWriter, opts map[string]interface{}, interaction data.Interaction) {

	embed := data.MessageEmbed{
		Title:       "About the Tri-Stat Companion",
		Type:        data.EmbedRichType,
		Description: fmt.Sprintf("This is the About section, yo!"),
		Color:       15217272,
		Thumbnail:   data.EmbedThumbnail{URL: "https://cdn.discordapp.com/app-icons/895003843375558687/d449adb5420b7f5c0145d3415abef46b.png"},
		Fields: []data.EmbedField{
			{Name: "Version", Value: "1.0.0"},
			{Name: "Invite To Another Server", Value: "[Click here to invite](https://discord.com/api/oauth2/authorize?client_id=895003843375558687&scope=applications.commands)"},
			{Name: "Get BESM4e!", Value: "You can get the PDF for BESM4e over on [DriveThruRPG](https://www.drivethrurpg.com/product/297755/BESM-Fourth-Edition-Big-Eyes-Small-Mouth) or get [BESM Naked](https://www.drivethrurpg.com/product/297761/BESM-Naked--Fourth-Edition-Big-Eyes-Small-Mouth?affiliate_id=275246)"},
			{Name: "Support the App!", Value: "You can support the app by purchasing [BESM4e with my affiliate link](https://www.drivethrurpg.com/product/297755/BESM-Fourth-Edition-Big-Eyes-Small-Mouth?affiliate_id=275246) or [BESM Naked with my affiliate link](https://www.drivethrurpg.com/product/297761/BESM-Naked--Fourth-Edition-Big-Eyes-Small-Mouth?affiliate_id=275246)"},
		},
	}

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
