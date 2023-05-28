package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/after23/telat-dong/handlers"
	"github.com/after23/telat-dong/util"
	"github.com/bwmarrin/discordgo"
)

const prefix string = "!telat"

func errHandler(message string, err error) {
	if err != nil {
		log.Panic(message, err)
	}
}

var sess *discordgo.Session
var config util.Config

var (
	integerOptionMinValue = 1.0
	dmPermission = false
	defaultMemberPermission int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name: "absen",
			Description: "test",
		},
		{
			Name: "hello",
			Description: "world!",
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"absen": handlers.Absen,
		"hello": handlers.Hello,
	}
)


func init() {
	
	config, err := util.LoadConfig(".")
	errHandler("Failed to read config file : ", err)
	
	sess, err = discordgo.New("Bot "+config.Token)
	errHandler("Failed to connect to the discord bot : ",err)
}

func init() {
	sess.AddHandler(func(sess *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(sess, i)
		}
	})
}

func main() {
	sess.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
		
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

		if args[1] == "image"{
			embed := &discordgo.MessageEmbed{
				Title: "Example Embed",
				Description: "This is an example of sending a MessageEmbed with an image in the response data.",
				Image: &discordgo.MessageEmbedImage{
					URL: "https://media.discordapp.net/attachments/1112298002938347550/1112298674870042634/image.png", // URL of the image
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
				if err != nil {
					fmt.Println("Error sending message:", err)
			}
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err := sess.Open()
	errHandler("Failed to open session : ",err)
	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := sess.ApplicationCommandCreate(sess.State.User.ID, config.PlaygroundID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	defer sess.Close()

	fmt.Println("yeow")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	for _, v := range registeredCommands {
			err := sess.ApplicationCommandDelete(sess.State.User.ID, config.PlaygroundID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
}