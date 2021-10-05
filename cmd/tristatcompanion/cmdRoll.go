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
	v, ok := opts["stat"]
	if ok {
		mod += v.(float64)
	}
	v, ok = opts["attribute"]
	if ok {
		mod += v.(float64)
	}
	v, ok = opts["misc"]
	if ok {
		mod += v.(float64)
	}

	numRoll := 2
	fmt.Printf("Edge/Obstacle: %+v\n\n", opts["edge-obstacle"])
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

	fmt.Printf("%+v", sortedTotal)

	if opts["edge-obstacle"].(string) == "minorEdge" {
		sortedTotal = sortedTotal[1:]
	}

	if opts["edge-obstacle"].(string) == "minorObstacle" || opts["edge-obstacle"].(string) == "majorObstacle" {
		sortedTotal = sortedTotal[:2]
	}

	if opts["edge-obstacle"].(string) == "majorEdge" {
		sortedTotal = sortedTotal[2:]
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
