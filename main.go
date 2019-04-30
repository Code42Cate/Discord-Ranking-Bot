package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	cache "github.com/patrickmn/go-cache"
)

// Discord constants
const (
	token     = "YOUR_DISCORD_BOT_TOKEN"
	channelID = "CHANNEL_WHERE_YOU_WANT_TO_TRACK"
)

var c = cache.New(24*30*time.Hour, cache.DefaultExpiration) // Clear cache every 30 days

func main() {

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Open a websocket connection to Discord and begin listening.
	dg.AddHandler(messageHandler)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	// Cleanly close down the Discord session.
	dg.Close()

}

// isValidCommand returns true if the given string has the format: !ranking n, where n is an optional number of any length
func isValidCommand(command string) bool {
	result, _ := regexp.MatchString("!ranking( \\d*)?", command)
	return result
}

// getArgument returns the n of !ranking n if it exists, if not the max parameter gets returned
// If the given argument is bigger than max or the String to Integer parsing failed, max gets returned
func getArgument(command string, max int) int {
	s := strings.Split(command, " ")
	if len(s) > 1 {
		index, err := strconv.Atoi(s[1])
		if err != nil || index > max {
			return max
		}
		return index
	}
	return max
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Case 1. Discord user uses the !ranking n command
	if m.Author.ID != s.State.User.ID && isValidCommand(m.Content) {

		ranking := rankByPostCount(c.Items())          // Get the cache as map and then create a sorted PairList
		firstN := getArgument(m.Content, len(ranking)) // Get command argument
		ranking = ranking[0:firstN]                    // Splice the ranking results
		output := "Ranking:\n"                         // Output string

		for i, pair := range ranking {
			output += fmt.Sprintf("%d. <@%s>: %d Posts\n", i+1, pair.Key, pair.Value)
		}
		if len(ranking) == 0 {
			output = "Nobody posted anything yet!"
		}
		_, sendError := s.ChannelMessageSend(m.ChannelID, output)
		if sendError != nil {
			// Handle the error however you think is appropriate
		}
	}
	// If the message is from the bot or n the wrong channel, return
	if m.Author.ID == s.State.User.ID || m.ChannelID != channelID {
		return
	}
	// Message has an attachement!
	// Note: Users only get +1 for any amount of attachments per message.
	// If you want to count the exact amount, just add len(m.Attachments) instead of 1 to the score
	if len(m.Attachments) > 0 {
		// Search for the userID in our cache storage
		user, found := c.Get(m.Author.ID)
		if found { // Found the user in cache, add 1 to his score
			c.Set(m.Author.ID, user.(int)+1, cache.DefaultExpiration)
		} else { // User not found in cache, create entry with 1 as score and his userID as key
			c.Set(m.Author.ID, 1, cache.DefaultExpiration)
		}
	}
}

// Modified an example which I found here:
// https://groups.google.com/forum/#!topic/golang-nuts/FT7cjmcL7gw
func rankByPostCount(postFrequencies map[string]cache.Item) PairList {
	pl := make(PairList, len(postFrequencies))
	i := 0
	for k, v := range postFrequencies {
		pl[i] = Pair{k, v.Object.(int)}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

// Pair is an element in our PairList
type Pair struct {
	Key   string
	Value int
}

// PairList contains Pair elements
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
