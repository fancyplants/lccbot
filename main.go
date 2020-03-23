package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	// ! don't tell anybody!
	TOKEN = ""
)

const (
	PREFIX = "/roll"
)

func init() {
	data, err := ioutil.ReadFile("token.txt")
	if err != nil {
		panic("No token! Put it in token.txt")
	}

	TOKEN = string(data)
}

func check(err error) {
	// TODO: don't use this
	if err != nil {
		panic(err)
	}
}

func main() {
	discord, err := discordgo.New("Bot " + TOKEN)
	check(err)

	discord.AddHandler(OnMessage)

	err = discord.Open()
	check(err)

	fmt.Println("Bot is now running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

func handleVideo(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return
	}
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		return
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("https://www.discordapp.com/channels/%s/%s", vs.GuildID, vs.ChannelID))
		}
	}
}

func handleRoll(s *discordgo.Session, m *discordgo.MessageCreate) {
	rest := strings.TrimSpace(m.Content[len(PREFIX):])

	rolls := strings.Split(rest, "+")

	sum := 0
	for _, roll := range rolls {
		text := strings.TrimSpace(roll)

		// if contains d, then it's a roll
		if strings.Contains(text, "d") {
			r, err := NewRoll(text)
			if err != nil {
				_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Incorrectly formatted command: %s", m.Content))
				if err != nil {
					log.Println(err)
				}
				return
			}

			sum += r.Calc()
		} else {
			num, err := strconv.Atoi(text)
			if err != nil {
				_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Incorrectly formatted command: %s", m.Content))
				if err != nil {
					log.Println(err)
				}
				return
			}

			sum += num
		}
	}

	// if successful, go ahead and send message!
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s's Roll: %v", m.Author.Username, sum))
	if err != nil {
		log.Println(err)
	}
}

func OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	// ignore any messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// syntax is now "/roll <roll> + <roll> + "

	if strings.HasPrefix(m.Content, "/roll") {
		handleRoll(s, m)
	}

	// generate video link for voice channel that user is in
	if strings.HasPrefix(m.Content, "/video") {
		handleVideo(s, m)
	}
}
