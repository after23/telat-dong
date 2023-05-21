package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/after23/telat-dong/util"
	"github.com/bwmarrin/discordgo"
)

const prefix string = "!telat"

func errHandler(message string, err error) {
	if err != nil {
		log.Panic(message, err)
	}
}

func main() {
	config, err := util.LoadConfig(".")
	errHandler("Failed to read config file : ", err)

	sess, err := discordgo.New("Bot "+config.Token)
	errHandler("Failed to connect to the discord bot : ",err)

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate){
		if m.Author.ID == s.State.User.ID {
			return
		}
		args := strings.Split(m.Content, " ")
		if args[0] != prefix{
			return
		}

		if args[1] == "hello" {
			s.ChannelMessageSend(m.ChannelID, "world!")
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	errHandler("Failed to open session : ",err)
	defer sess.Close()

	fmt.Println("yeow")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}