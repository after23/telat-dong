package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/after23/telat-dong/models"
	"github.com/after23/telat-dong/util"
	"github.com/bwmarrin/discordgo"
)


func Hello(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Heya",
		},
	})
}

func absen(s *discordgo.Session, i *discordgo.InteractionCreate, url string){
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
	if url == util.Conf().StatusURL{
		embed.Title = "Checkin absen"
	}
	embeds := []*discordgo.MessageEmbed{embed}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
		},
	})
	chanOwner := func() <-chan models.Result {
		ch := make(chan models.Result)
		go func() {
			defer close(ch)
			util.Request(s, ch, url)
		}()
		return ch
	}

	
	res := <-chanOwner()
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

func Absen(s *discordgo.Session, i *discordgo.InteractionCreate) {
	absen(s, i, util.Conf().AbsenURL)
}

func SlashStatus(s *discordgo.Session, i *discordgo.InteractionCreate){
	absen(s,i,util.Conf().StatusURL)
}

func SlashPing(s *discordgo.Session, i *discordgo.InteractionCreate){
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

func SlashSkip(s *discordgo.Session, i *discordgo.InteractionCreate){
	embed := &discordgo.MessageEmbed{
		Title: "Skip",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	embed.Image = &discordgo.MessageEmbedImage{
		URL: "https://cdn.discordapp.com/attachments/1112298002938347550/1113268135169118208/gdlu.jpg",
	}
	embeds := []*discordgo.MessageEmbed{embed}

	respData := &discordgo.InteractionResponseData{
		Embeds: embeds,
	}

	resp := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: respData,
	}

	err := s.InteractionRespond(i.Interaction, resp)

	if err != nil {
		log.Println("SlashSkip Error: ", err)
		return
	}
}

