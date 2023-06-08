package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/after23/telat-dong/models"
	"github.com/bwmarrin/discordgo"
)

func Request(s *discordgo.Session, ch chan<- models.Result, url string) {
	var res models.Result
	client := http.Client{
		Timeout: 6 * time.Minute,
	}

	resp, err := client.Get(url)
	if err != nil {
		res.Status = models.Failed
		res.Message = fmt.Sprintf("Error doing HTTP GET\n%v", err)
		ch <- res
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			res.Status = models.Failed
			res.Message = fmt.Sprintf("Error reading body\n%v", err)
			ch <- res
			return
		}
		res.Status = models.Failed
		res.Message = fmt.Sprintf("HTTP %d\n%v", resp.StatusCode, string(body))
		ch <- res
		return
	}
	fmt.Println("Image request resolved.")

	if !(resp.Header.Get("Content-Type") == "image/png"){
		body,err := ioutil.ReadAll(resp.Body)
		if err != nil {
			res.Status = models.Failed
			res.Message = fmt.Sprintf("Error reading body\n%v", err)
			ch <- res
			return
		}
		res.Status = models.Success
		res.Message = string(body)
		ch <- res
		return
	}

	file := &discordgo.File{
		Name:        "image.png",
		ContentType: "image/png",
		Reader:      resp.Body,
	}
	files := []*discordgo.File{file}
	test := &discordgo.MessageSend{Files: files}
	msg, err := s.ChannelMessageSendComplex(Conf().ImageDumpID, test)
	if err != nil {
		res.Status = models.Failed
		res.Message = fmt.Sprintf("Error Uploading image\n%v", err)
		ch <- res
		return
	}
	imgUrl := msg.Attachments[0].URL

	err = s.ChannelMessageDelete(Conf().ImageDumpID, msg.ID)
	if err != nil {
		log.Println("models.Failed to delete message: ", err)
	}

	res.Status = models.Success
	res.Message = "Status: Success"
	res.URL = imgUrl
	ch <- res
	return
}
