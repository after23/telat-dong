package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Hello(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "world",
				},
			})
		}

func Absen(s *discordgo.Session, i *discordgo.InteractionCreate) {

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
	var wg sync.WaitGroup
	wg.Add(1)
	go absen(&wg, s, i.Interaction, "https://telat-api.onrender.com/talenta/absen/?api-key=hQcx29p8gWXyq6wdQykFAxcpb8bqnwsx")
	wg.Wait()
	embeds[0].Description = "Status: Finished"
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &embeds,
	})
}

func absen(wg *sync.WaitGroup, s *discordgo.Session , i *discordgo.Interaction, url string ) {
		// Make an HTTP request to retrieve the image.
		resp, err := http.Get(url)
		if err != nil {
			// fmt.Println("Error retrieving image:", err)
			wg.Done()
			log.Panic("Error doing http request ", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			wg.Done()
			log.Panic("somehing went oopsies ", err)
		}

		// The image request has resolved successfully.
		fmt.Println("Image request resolved.")

		// Upload the file to Discord and get the attachment URL.
		msg, err := s.ChannelFileSend("1112298002938347550" , "image.png", resp.Body)
		if err != nil {
			wg.Done()
			log.Panic("error sending image :", err)
		}

		embed := &discordgo.MessageEmbed{
				Title: "Absen Result",
				Description: "Success",
				Image: &discordgo.MessageEmbedImage{
					URL: msg.Attachments[0].URL, // URL of the image
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}
		_, err = s.ChannelMessageSendEmbed("1109748811405983786", embed)
			if err != nil {
				fmt.Println("Error sending message:", err)
		}
		wg.Done()

		// embed = &discordgo.MessageEmbed{
		// Title: "Absen",
		// Description: "Status : Processing",
		// Timestamp: time.Now().Format(time.RFC3339),
		// }

		// embeds := []*discordgo.MessageEmbed{embed}

		// updateResponse := &discordgo.WebhookEdit{
		// 	Embeds: &embeds,
		// }
		
		// s.InteractionResponseEdit(i, updateResponse)		
		// return attachment.URL
}