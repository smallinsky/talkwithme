// use idea from https://github.com/moby/moby/blob/master/pkg/namesgenerator/names-generator.go
// to generates random user name when client connects to server via telnet and username is not yet privided.

package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	left = [...]string{
		"adoring",
		"awesome",
		"beautiful",
		"blissful",
		"bold",
		"clever",
		"cool",
		"compassionate",
		"condescending",
		"dazzling",
		"dreamy",
		"ecstatic",
		"mystifying",
		"naughty",
		"nervous",
		"nifty",
		"nostalgic",
		"pedantic",
		"pensive",
		"quirky",
		"quizzical",
		"recursing",
		"romantic",
		"sad",
		"serene",
		"sharp",
		"silly",
		"sleepy",
		"strange",
		"stupefied",
		"sweet",
	}

	right = [...]string{
		"adam",
		"angelina",
		"ariel",
		"ela",
		"elon",
		"eva",
		"jeo",
		"john",
		"julia",
		"lily",
		"magda",
		"mark",
		"richard",
		"robert",
		"dubinsky",
		"gould",
		"ptolemy",
		"wu",
		"nobel",
		"newton",
	}
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func getRandomName() string {
	return fmt.Sprintf("%s_%s", left[rand.Intn(len(left))], right[rand.Intn(len(right))])
}
