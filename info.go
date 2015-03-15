package main

import (
	"fmt"
	"math/rand"
	"time"
)

var lines = []string{
	"Talking's for functioning people",
	"Follow me on Twitter @THUNDERGROOVE",
	"FA for lyfe",
	"The exoplanet has failed",
	`"Oh, hey Jon." "Oh, yeah." "How's it going, man?" "Oh, I'm alright, whate- whatever." "Did you-did you hear about that party?" "No, no, no, where's that?" "Dude. I-I... is it 26 and P or-or what?" "Dude, it's-it's L." "What?!" "Yeah." "Dude, I've been driving around for like, fuckin' half an hour." "I thought you were picking me up right now." "Well, I was going to, but I had to stop at the store" "and get some fucking shit." "You had to stop at the store?" "Well, what are we going to fucking drink?" "Do you-do you still have the money that I gave you earlier?" "Well, not really, because I fuckin' had to buy beer!" "That's fucked up, man. Every time, I pick you up, and I..." "Didn't you..." "Spend my money on you." "You know what? Whatever, all right. What the..." "Whatever. I don't give a shit." "Yeah, well, I'm still going. Are you going to come?" "I'll pick you up still, but... I mean," "the party's gonna be over by the time I get over there." "Whatever. You know, it doesn't really matter right now because" "you are rock solid." "Rock solid?" "Rock solid." "We're both rock solid." "That's right."`,
}

var info = fmt.Sprintln("SDETool2 twice as good as SDETool",
	"Written by Nick Powell 'THUNDERGROOVE' in Golang")

func PrintInfo() {
	fmt.Println(info)
	fmt.Println("	\u266B " + randomLine() + " \u266B")
}
func randomLine() string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(lines))
	return lines[i]
}
