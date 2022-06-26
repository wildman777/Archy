package channelMessage

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// Data struct from Pubsub
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// Payload struct that is expected
type Payload struct {
	ChannelId string `json:"channel_id"`
	Message   string `json:"message"`
	Image     string `json:"image"`
}

// Unmarshal received context and call proper function that send message
func ChannelMessage(ctx context.Context, m PubSubMessage) error {

	log.Printf("Starting!")

	// Unmarshal data to a valid payload
	var payload Payload
	json.Unmarshal(m.Data, &payload)

	return SendChannelMessage(&payload)
}

// Send a message to a specific Discord channel
func SendChannelMessage(payload *Payload) error {

	// Instanciate discord bot
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return err
	}

	// Instanciate a direct message channel
	channel, err := dg.Channel(payload.ChannelId)

	if err != nil {
		error_message := []byte(err.Error())
		error_400_regex, _ := regexp.Compile("400")
		if len(error_400_regex.Find(error_message)) > 0 {
			panic("Can't create Channel - Bad ChannelId")
		}
		error_401_regex, _ := regexp.Compile("401")
		if len(error_401_regex.Find(error_message)) > 0 {
			panic("Unauthorized to create the connection. Verify Discord Token")
		}
		return err
	}

	// Create Message object
	var messageData discordgo.MessageSend
	messageData.Content = payload.Message
	// Create file reader for image
	imageAsBytes, _ := base64.StdEncoding.DecodeString(payload.Image)

	// Add the image only if it is provided
	if len(imageAsBytes) > 0 {
		reader := bytes.NewReader(imageAsBytes)
		var file discordgo.File
		var files []*discordgo.File = []*discordgo.File{&file}

		files[0].Name = "image.png"
		files[0].ContentType = "png"
		files[0].Reader = reader

		messageData.Files = files
	}

	// Send message to user
	log.Printf("Sending to channel: " + payload.ChannelId)
	log.Printf("Message is: " + payload.Message)
	log.Printf("Image was send with it: " + strconv.FormatBool(len(imageAsBytes) > 0))
	_, err = dg.ChannelMessageSendComplex(channel.ID, &messageData)

	if err != nil {
		panic("Message didn't make it" + err.Error())
	}

	log.Printf("Done!")
	return nil
}