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

var sess *discordgo.Session
var config util.Config

var (
	integerOptionMinValue = 1.0
	dmPermission = false
	defaultMemberPermission int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
		{
			Name: "absen",
			Description: "Ga lagi telat-telat absen",
		},
		{
			Name: "hello",
			Description: "testing",
		}, 
		{
			Name: "ping",
			Description: "Bangunin service",
		},
		// {
		// 	Name: "responses",
		// 	Description: "responses",
		// 	Options: []*discordgo.ApplicationCommandOption{
		// 		{
		// 			Name:        "resp-type",
		// 			Description: "Response type",
		// 			Type:        discordgo.ApplicationCommandOptionInteger,
		// 			Choices: []*discordgo.ApplicationCommandOptionChoice{
		// 				{
		// 					Name:  "Channel message with source",
		// 					Value: 4,
		// 				},
		// 				{
		// 					Name:  "Deferred response With Source",
		// 					Value: 5,
		// 				},
		// 			},
		// 			Required: true,
		// 		},
		// 	},
		// },
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, config *util.Config){
		"absen": handlers.Absen,
		"hello": handlers.Hello,
		"ping": handlers.SlashPing,
		// "responses": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// 	// Responses to a command are very important.
		// 	// First of all, because you need to react to the interaction
		// 	// by sending the response in 3 seconds after receiving, otherwise
		// 	// interaction will be considered invalid and you can no longer
		// 	// use the interaction token and ID for responding to the user's request

		// 	content := ""
		// 	// As you can see, the response type names used here are pretty self-explanatory,
		// 	// but for those who want more information see the official documentation
		// 	switch i.ApplicationCommandData().Options[0].IntValue() {
		// 	case int64(discordgo.InteractionResponseChannelMessageWithSource):
		// 		content =
		// 			"You just responded to an interaction, sent a message and showed the original one. " +
		// 				"Congratulations!"
		// 		content +=
		// 			"\nAlso... you can edit your response, wait 5 seconds and this message will be changed"
		// 	default:
		// 		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// 			Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
		// 		})
		// 		if err != nil {
		// 			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		// 				Content: "Something went wrong",
		// 			})
		// 		}
		// 		return
		// 	}

		// 	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// 		Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
		// 		Data: &discordgo.InteractionResponseData{
		// 			Content: content,
		// 		},
		// 	})
		// 	if err != nil {
		// 		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		// 			Content: "Something went wrong",
		// 		})
		// 		return
		// 	}
		// 	time.AfterFunc(time.Second*5, func() {
		// 		content := content + "\n\nWell, now you know how to create and edit responses. " +
		// 			"But you still don't know how to delete them... so... wait 10 seconds and this " +
		// 			"message will be deleted."
		// 		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		// 			Content: &content,
		// 		})
		// 		if err != nil {
		// 			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		// 				Content: "Something went wrong",
		// 			})
		// 			return
		// 		}
		// 		time.Sleep(time.Second * 10)
		// 		s.InteractionResponseDelete(i.Interaction)
		// 	})
		// },
	}
)


func init() {
	
	config, err := util.LoadConfig(".")
	util.ErrHandler("Failed to read config file : ", err)
	
	sess, err = discordgo.New("Bot "+config.Token)
	util.ErrHandler("Failed to connect to the discord bot : ",err)
}

func init() {
	sess.AddHandler(func(sess *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(sess, i, &config)
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
				Title: "coba image",
				Description: "image embed",
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

		if args[1] == "absen"{
			handlers.TempAbsen(s, m.ChannelID, m.Author.ID, &config)
		}

		if args[1] == "ping"{
			handlers.Ping(s, m.ChannelID)
		}
		if args[1] == "status"{
			handlers.Status(s, m.ChannelID, m.Author.ID)
		}
		if args[1] == "skip"{
			handlers.Skip(s, m.ChannelID)
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err := sess.Open()
	util.ErrHandler("Failed to open session : ",err)
	// log.Println("Adding commands...")
	// registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	// for i, v := range commands {
	// 	cmd, err := sess.ApplicationCommandCreate(sess.State.User.ID, config.PlaygroundID, v)
	// 	if err != nil {
	// 		log.Printf("Cannot create '%v' command: %v.\nSkipping slash command.", v.Name, err)
	// 		break
	// 	}
	// 	registeredCommands[i] = cmd
	// }
	defer sess.Close()

	fmt.Println("yeow")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc

	// for _, v := range registeredCommands {
	// 		err := sess.ApplicationCommandDelete(sess.State.User.ID, config.NFGuildID, v.ID)
	// 		if err != nil {
	// 			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
	// 		}
	// 	}
}