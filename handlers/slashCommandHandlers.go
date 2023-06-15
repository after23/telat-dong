package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	// go absen(ch, s, i.Interaction, config.ImageDumpID)

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

func absen(ch chan<- models.Result, s *discordgo.Session , i *discordgo.Interaction, imageDumpChannel string ) {
	// Make an HTTP request to retrieve the image.
	client := http.Client{
    Timeout: 6 * time.Minute,
	}
	resp, err := client.Get("https://telat-api.onrender.com/talenta/absen/?api-key=hQcx29p8gWXyq6wdQykFAxcpb8bqnwsx")
	var res models.Result
	if err != nil {
		res.Status = models.Failed
		res.Message = fmt.Sprintf("Error doing http request: %v", err)
		ch <- res
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body,_ := ioutil.ReadAll(resp.Body)
		res = models.Result{
			Status: models.Failed,
			Message: fmt.Sprintf("Something went oopsies: %d - %s", resp.StatusCode, string(body)),
		}
		ch <- res
		return
	}

	// The image request has resolved successfully.
	fmt.Println("Image request resolved.")

	// Upload the file to Discord and get the attachment URL.
	msg, err := s.ChannelFileSend(imageDumpChannel , "image.png", resp.Body)
	if err != nil {
		res = models.Result{
			Status: models.Failed,
			Message: fmt.Sprintf("error sending image: %v", err),
		}
		ch <- res
		return
	}

	imageURL := msg.Attachments[0].URL

	err = s.ChannelMessageDelete(imageDumpChannel, msg.ID)
	
	if err != nil {
		res = models.Result{
			Status: models.Failed,
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
			res = models.Result{
				Status: models.Failed,
				Message: fmt.Sprintf("error sending Message Embed: %v", err),
			}
			ch <- res
			return
	}
	res = models.Result{
		Status: models.Success,
		Message: "",
	}
	ch <- res

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
	go ping(ch)
	defer close(ch)

	res := <- ch
	embeds[0].Description = res.Message
	embeds[0].Timestamp = time.Now().Format(time.RFC3339)
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &embeds,
	})
	return
}

func ping(ch chan<- models.Result){
	client := http.Client{
		Timeout: 5 * time.Minute,
	}

	var res models.Result

	resp, err := client.Get("https://telat-api.onrender.com/ping")
	if err != nil {
		res.Status = models.Failed
		res.Message = fmt.Sprintf("Status: Failed\n%v", err)
		ch <- res
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK{
		res.Status = models.Failed
		res.Message = fmt.Sprintf("Status: Failed\nHTTP %d", resp.StatusCode)
		ch <- res
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		res.Status = models.Failed
		res.Message = fmt.Sprintf("Status: Failed\n%v", err)
		ch <- res
		return
	}

	res.Status = models.Success
	res.Message = fmt.Sprintf("Status: Success\n%s", string(body))
	ch <- res
	return
}

