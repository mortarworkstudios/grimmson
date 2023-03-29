package scraper

import (
	"bufio"
	"log"
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	archivePath   = "./archive"
	archiveSuffix = "_archive_"
)

type ServerScraper struct {
	botConf *Config
	sesh    *discordgo.Session
}

func NewServerScraper(config *Config) *ServerScraper {
	return &ServerScraper{
		botConf: config,
	}
}

func (sc *ServerScraper) InitScraper() error {
	log.Println("Initializing Discord Server Scraper")

	var err error
	sc.sesh, err = discordgo.New("Bot " + sc.botConf.DiscordToken)

	if err != nil {
		return err
	}

	err = sc.sesh.Open()
	if err != nil {
		return err
	}

	// Get an array of text channels
	var textChannels []*discordgo.Channel
	for _, guild := range sc.sesh.State.Guilds {
		channels, _ := sc.sesh.GuildChannels(guild.ID)
		for _, c := range channels {
			if c.Type != discordgo.ChannelTypeGuildText {
				continue
			} else {
				textChannels = append(textChannels, c)
			}

		}
	}

	os.Mkdir(archivePath, os.ModePerm)

	var wg sync.WaitGroup
	for _, channel := range textChannels {
		log.Printf("Starting archive for %s\n", channel.Name)
		wg.Add(1)
		go sc.BulkDownloadMessages(&wg, channel, archivePath)
	}
	wg.Wait()

	sc.sesh.Close()
	return nil
}

func (sc *ServerScraper) BulkDownloadMessages(wg *sync.WaitGroup, channel *discordgo.Channel, archivePath string) {
	defer wg.Done()
	var messages []*discordgo.Message
	var err error
	dateStamp := time.Now().Format(time.RFC3339)
	// Create an archive file to write to
	var archiveFile *os.File
	archiveFile, err = os.Create(archivePath + "/" + channel.Name + archiveSuffix + dateStamp + ".txt")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer archiveFile.Close()
	archiveWriter := bufio.NewWriter(archiveFile)

	doneArchiving := false
	beforeID := ""
	for !doneArchiving {
		// Get all the messages we can (max is a limit per API call)
		messages, err = sc.sesh.ChannelMessages(channel.ID, 100, beforeID, "", "")
		if err != nil {
			log.Fatal(err.Error())
		}

		if len(messages) == 0 {
			doneArchiving = true
			break
		}

		// Loop through all the messages we fetched
		for _, msg := range messages {
			// Grab the last ID to get more messages from before
			beforeID = msg.ID
			archiveWriter.WriteString(msg.Content + "\n")
		}
	}

	log.Printf("Archiving complete for %s\n", channel.Name)
	archiveWriter.Flush()
}
