package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/after23/telat-dong/util"
	"github.com/bwmarrin/discordgo"
)

func initMessageEmbedEdit(s *discordgo.Session, channelID string, messageID string, embed *discordgo.MessageEmbed, res Result) {
	newDesc := fmt.Sprintf("Status: %s", statusMap[res.Status])
	if res.Status == Failed{
		newDesc += fmt.Sprintf("\n%v", res.Message)
	}
	embed.Description = newDesc
	embed.Timestamp = time.Now().Format(time.RFC3339)
	_,err := s.ChannelMessageEditEmbed(channelID, messageID, embed)
	if (err != nil){
		log.Println(err)
		return
	}
}

func TempAbsen(s *discordgo.Session, channelID string, author_id string, config *util.Config) {
	//author validation
	if author_id != "188656104673247232" {
		embed := &discordgo.MessageEmbed{
			Title: "gk dlu",
			Image: &discordgo.MessageEmbedImage{
				URL: "https://cdn.discordapp.com/attachments/1112298002938347550/1113268135169118208/gdlu.jpg", // URL of the image
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		_, err := s.ChannelMessageSendEmbed(channelID, embed)
		if err != nil {
			log.Println("Error sending MessageEmbed: ", err)
			return
		}
		return
	}
	
	embed := &discordgo.MessageEmbed{
		Title:       "Absen",
		Description: "Status: Processing...",
	}

	//initial message
	msg, err := s.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		log.Println("Error sending MessageEmbed: ", err)
		return
	}

	client := http.Client{
		Timeout: 6 * time.Minute,
	}
	var res Result
	//get request to the service
	resp, err := client.Get("https://telat-api.onrender.com/talenta/absen/?api-key=hQcx29p8gWXyq6wdQykFAxcpb8bqnwsx")
	if err != nil {
		res.Status=Failed
		res.Message=(fmt.Sprintf("Error doing Get Request: %v", err))
		initMessageEmbedEdit(s, channelID, msg.ID, embed, res)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body,_ := ioutil.ReadAll(resp.Body)
		res.Status=Failed
		res.Message=fmt.Sprintf("Error - %d:\n%v", resp.StatusCode, string(body))
		initMessageEmbedEdit(s, channelID, msg.ID, embed, res)
		return
	}
	// The image request has resolved successfully.
	fmt.Println("Image request resolved.")

	// Upload the file to Discord and get the attachment URL.
	image, err := s.ChannelFileSend(config.ImageDumpID, "image.png", resp.Body)
	if err != nil {
		return
	}
	imageURL := image.Attachments[0].URL

	err = s.ChannelMessageDelete(config.ImageDumpID, msg.ID)
	messageEmbed := &discordgo.MessageEmbed{
		Title:       "Absen Result",
		Description: "Success",
		Image: &discordgo.MessageEmbedImage{
			URL: imageURL, // URL of the image
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	_, err = s.ChannelMessageSendEmbed(channelID, messageEmbed)
	if err != nil {
		log.Println("Error sending embed message(result): ", err)
		return
	}
	res.Status=Success
	initMessageEmbedEdit(s, channelID, msg.ID, embed, res)

}

func Ping(s *discordgo.Session, channelID string) {
	embed := &discordgo.MessageEmbed{}
	embed.Title = "Pinging the Service"
	embed.Description = "Status: Processing.."
	embed.Timestamp = time.Now().Format(time.RFC3339)

	msg, err := s.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		log.Println("Error sending embed message: ", err)
	}

	client := http.Client{
		Timeout: 5 * time.Minute,
	}

	resp, err := client.Get("https://telat-api.onrender.com/ping")
	if err != nil {
		embed.Description = fmt.Sprintf("Status: Failed\n%v", err)
		embed.Timestamp = time.Now().Format(time.RFC3339)
		s.ChannelMessageEditEmbed(channelID, msg.ID, embed)
		return
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		embed.Description = fmt.Sprintf("Status: Failed\n HTTP %d\n%v", resp.StatusCode, string(body))
		embed.Timestamp = time.Now().Format(time.RFC3339)
		s.ChannelMessageEditEmbed(channelID, msg.ID, embed)
	}

	defer resp.Body.Close()
	body,_ := ioutil.ReadAll(resp.Body)
	embed.Description = fmt.Sprintf("Status: Success\n%s", string(body))
	embed.Timestamp = time.Now().Format(time.RFC3339)
	s.ChannelMessageEditEmbed(channelID, msg.ID, embed)
}