package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/after23/telat-dong/models"
	"github.com/after23/telat-dong/util"
	"github.com/bwmarrin/discordgo"
)

func initMessageEmbedEdit(s *discordgo.Session, channelID string, messageID string, embed *discordgo.MessageEmbed, res models.Result) {
	newDesc := fmt.Sprintf("Status: %s", models.StatusMap[res.Status])
	if res.Message != ""{
		newDesc += fmt.Sprintf("\n%v", res.Message)
	}
	embed.Description = newDesc
	embed.Timestamp = time.Now().Format(time.RFC3339)
	if res.URL != ""{
		embed.Image = &discordgo.MessageEmbedImage{
			URL: res.URL,
		}
	}
	_,err := s.ChannelMessageEditEmbed(channelID, messageID, embed)
	if (err != nil){
		log.Println(err)
		return
	}
}

func TempAbsen(s *discordgo.Session, channelID string, author_id string, config *util.Config) {
	//author validation
	if author_id != util.Conf().MyID {
		Skip(s, channelID)
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

	ch := make(chan models.Result)
	defer close(ch)
	go util.Request(s, ch, util.Conf().AbsenURL)
	res := <- ch
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
		return
	}
	ch := make(chan models.Result)
	go util.Request(s, ch, models.PingURL)
	res := <- ch

	initMessageEmbedEdit(s, channelID, msg.ID, embed, res)
}

func Skip(s *discordgo.Session, channelID string) {

	embed := &discordgo.MessageEmbed{
		Title: "skip",
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

func Status(s *discordgo.Session, channelID, authorID string){
	if authorID != util.Conf().MyID {
		Skip(s,channelID)
		return
	}
	embed := &discordgo.MessageEmbed{}
	embed.Title = "Checking absen"
	embed.Description = "Status: Processing.."
	embed.Timestamp = time.Now().Format(time.RFC3339)

	msg, err := s.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		log.Println("Error sending MessageEmbed: ", err)
		return
	}

	ch := make(chan models.Result)
	defer close(ch)
	go util.Request(s, ch, util.Conf().StatusURL)

	res := <- ch
	initMessageEmbedEdit(s, channelID, msg.ID, embed, res)
	return
	
}
