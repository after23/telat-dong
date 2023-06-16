package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/after23/telat-dong/models"
	"github.com/after23/telat-dong/util"
	"github.com/bwmarrin/discordgo"
)


func Hello(s *discordgo.Session, i *discordgo.InteractionCreate, config *util.Config) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Heya",
		},
	})
}

func Absen(s *discordgo.Session, i *discordgo.InteractionCreate, config *util.Config) {
	if i.Member.User.ID != "188656104673247232"{
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "siape lo",
			},
		})
		return
	}
	embed := &discordgo.MessageEmbed{
		Title: "Absen",
		Description: "Status: Processing",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	embeds := []*discordgo.MessageEmbed{embed}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
		},
	})
	ch := make(chan models.Result)
	defer close(ch)
	go util.Request(s, ch, util.Conf().AbsenURL)

	res := <- ch
	embeds[0].Description = fmt.Sprintf("Status: %s", models.StatusMap[res.Status])
	if res.Status == models.Failed {
		embeds[0].Description = fmt.Sprintf("Status: %s\n %s", models.StatusMap[res.Status],res.Message)
	}
	if res.URL != "" {
		embeds[0].Image = &discordgo.MessageEmbedImage{
			URL: res.URL,
		}
	}
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &embeds,
	})
}

func SlashPing(s *discordgo.Session, i *discordgo.InteractionCreate, config *util.Config){
	embed := &discordgo.MessageEmbed{
		Title: "Pinging the Service",
		Description: "Status: Processing..",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	embeds := []*discordgo.MessageEmbed{embed}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
		},
	})

	if err != nil {
		log.Println("Failed Sending interaction response: ", err)
		return
	}

	ch := make(chan models.Result)
	go util.Request(s, ch, models.PingURL)
	defer close(ch)

	res := <- ch
	embeds[0].Description = fmt.Sprintf("Status: %s\n%s", models.StatusMap[res.Status], res.Message)
	embeds[0].Timestamp = time.Now().Format(time.RFC3339)
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &embeds,
	})
	return
}

