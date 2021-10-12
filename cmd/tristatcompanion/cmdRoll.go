package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/c-bec-k/tristatcompanion/internal/data"
)

func (bot *application) ReplyRoll(w http.ResponseWriter, opts map[string]interface{}, interaction data.Interaction) {

	var mod float64
	mods := [3]int{0, 0, 0}
	v, ok := opts["stat"]
	if ok {
		mod += v.(float64)
		mods[0] = int(v.(float64))
	}
	v, ok = opts["attribute"]
	if ok {
		mod += v.(float64)
		mods[1] = int(v.(float64))
	}
	v, ok = opts["misc"]
	if ok {
		mod += v.(float64)
		mods[2] = int(v.(float64))
	}

	numRoll := 2
	// fmt.Printf("Edge/Obstacle: %+v\n\n", opts["edge-obstacle"])
	if opts["edge-obstacle"] != nil {
		if opts["edge-obstacle"].(string) == "minorEdge" || opts["edge-obstacle"].(string) == "minorObstacle" {
			numRoll = 3
		}

		if opts["edge-obstacle"].(string) == "majorEdge" || opts["edge-obstacle"].(string) == "majorObstacle" {
			numRoll = 4
		}
	}

	var diceResult []die

	for i := 0; i < numRoll; i++ {
		diceResult = append(diceResult, getDieResult())
	}

	var emoji []string
	for _, v := range diceResult {
		emoji = append(emoji, v.emoji)
	}

	var sortedTotal []int
	for _, v := range diceResult {
		sortedTotal = append(sortedTotal, v.result)
	}
	sort.Ints(sortedTotal)

	var total int

	// fmt.Printf("%+v", sortedTotal)
	var edOb string

	if opts["edge-obstacle"] != nil {
		if opts["edge-obstacle"].(string) == "minorEdge" {
			sortedTotal = sortedTotal[1:]
			edOb = "Minor Edge"
		}

		if opts["edge-obstacle"].(string) == "minorObstacle" || opts["edge-obstacle"].(string) == "majorObstacle" {
			sortedTotal = sortedTotal[:2]
			edOb = "Major Obstacle"
		}

		if opts["edge-obstacle"].(string) == "minorObstacle" {
			edOb = "Minor Obstacle"
		}

		if opts["edge-obstacle"].(string) == "majorEdge" {
			sortedTotal = sortedTotal[2:]
			edOb = "Major Edge"
		}
	}

	for _, v := range sortedTotal {
		total += v
	}
	total += int(mod)

	embed := data.MessageEmbed{
		Color:       15217272,
		Description: fmt.Sprintf("%v %+d", strings.Join(emoji, " "), int(mod)),
	}

	if opts["description"] != nil {
		embed.Author.Name = strconv.Quote(opts["description"].(string))
	}

	embed.Title = fmt.Sprintf("You got a %+v total!", total)

	var fields []data.EmbedField

	modsField := data.EmbedField{"Modifiers", fmt.Sprintf("%+d%+d%+d (Stat+Attribute+Misc)", mods[0], mods[1], mods[2]), true}

	fields = append(fields, modsField)

	if edOb != "" {
		moreDiceField := data.EmbedField{"Edges/Obstacles", edOb, true}
		fields = append(fields, moreDiceField)
	}

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
