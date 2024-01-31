package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Poll structure
type Poll struct {
	Question string
	Options  []string
	Votes    map[string]int
}

// Map to keep track of active polls
var activePolls = make(map[string]*Poll)
var token = os.Getenv("DISCORD_BOT_TOKEN")

func main() {
	log.Println(token)
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// Check if the message is a command to create a poll
	if strings.HasPrefix(m.Content, "!createpoll") {
		createPoll(s, m)
	}
}

// Function to create a poll
func createPoll(s *discordgo.Session, m *discordgo.MessageCreate) {
	content := strings.TrimPrefix(m.Content, "!createpoll")
	args := strings.Split(content, "\"")

	// Check for the correct format
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Incorrect poll format. Usage: !createpoll \"Question?\" \"Option1\" \"Option2\" ...")
		return
	}

	// Extract question and options
	question := args[1]
	var options []string
	for i := 3; i < len(args); i += 2 {
		options = append(options, args[i])
	}

	// Create a new poll
	poll := &Poll{
		Question: question,
		Options:  options,
		Votes:    make(map[string]int),
	}
	pollID := generatePollID() // Implement a function to generate a unique poll ID
	activePolls[pollID] = poll

	// Display the poll
	displayPoll(s, m, pollID, poll)
}

func generatePollID() string {
	// rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("poll_%d_%d", time.Now().Unix(), rand.Intn(1000))
}

// Function to display the poll
func displayPoll(s *discordgo.Session, m *discordgo.MessageCreate, pollID string, poll *Poll) {
	// Create an embed or a simple message to display the poll
	// For each option, add a reaction for voting
	// ...

	// Example: Send a simple message (implement better formatting and reaction handling)
	msg := fmt.Sprintf("Poll: %s\n", poll.Question)
	for _, option := range poll.Options {
		msg += fmt.Sprintf("- %s\n", option)
	}
	s.ChannelMessageSend(m.ChannelID, msg)
}
