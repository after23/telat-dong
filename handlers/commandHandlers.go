package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	Success = iota
	Failed
)

type Result struct {
	Status int
	Message string
}

func Hello(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "world",
		},
	})
}

func Absen(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	ch := make(chan Result)
	go absen(ch, s, i.Interaction, "https://telat-api.onrender.com/talenta/absen/?api-key=hQcx29p8gWXyq6wdQykFAxcpb8bqnwsx")

	res := <- ch
	embeds[0].Description = "Status: Finished"
	if res.Status == Failed {
		embeds[0].Description = fmt.Sprintf("Status: Failed\n %s", res.Message)
	}
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &embeds,
	})
}

func absen(ch chan<- Result, s *discordgo.Session , i *discordgo.Interaction, url string ) {
	// Make an HTTP request to retrieve the image.
	resp, err := http.Get(url)
	var res Result
	if err != nil {
		res.Status = Failed
		res.Message = fmt.Sprintf("Error doing http request: %v", err)
		ch <- res
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body,_ := ioutil.ReadAll(resp.Body)
		res = Result{
			Status: Failed,
			Message: fmt.Sprintf("Something went oopsies: %d - %s", resp.StatusCode, string(body)),
		}
		ch <- res
		return
	}

	// The image request has resolved successfully.
	fmt.Println("Image request resolved.")

	// Upload the file to Discord and get the attachment URL.
	msg, err := s.ChannelFileSend("1112298002938347550" , "image.png", resp.Body)
	if err != nil {
		res = Result{
			Status: Failed,
			Message: fmt.Sprintf("error sending image: %v", err),
		}
		ch <- res
		return
	}

	imageURL := msg.Attachments[0].URL

	err = s.ChannelMessageDelete("1112298002938347550", msg.ID)
	
	if err != nil {
		res = Result{
			Status: Failed,
			Message: fmt.Sprintf("Failed Deleting Message %v", err),
		}
		return
	}

	embed := &discordgo.MessageEmbed{
			Title: "Absen Result",
			Description: "Success",
			Image: &discordgo.MessageEmbedImage{
				URL: imageURL, // URL of the image
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
	_, err = s.ChannelMessageSendEmbed(i.ChannelID, embed)
		if err != nil {
			res = Result{
				Status: Failed,
				Message: fmt.Sprintf("error sending Message Embed: %v", err),
			}
			ch <- res
			return
	}
	res = Result{
		Status: Success,
		Message: "",
	}
	ch <- res

}