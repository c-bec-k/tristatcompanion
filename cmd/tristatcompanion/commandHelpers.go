package main

import (
	"math/rand"
	"time"
)

type die struct {
	result int
	emoji  string
}

var dice = []die{
	{result: 1, emoji: "<:Result1:768110374742130698>"},
	{result: 2, emoji: "<:Result2:768110374570164285>"},
	{result: 3, emoji: "<:Result3:768110375131938816>"},
	{result: 4, emoji: "<:Result4:768110375006502914>"},
	{result: 5, emoji: "<:Result5:768110375275069460>"},
	{result: 6, emoji: "<:Result6:768110375416889344>"},
}

func getDieResult() die {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(6)
	return dice[num]
}
