package main

import (
	"log"
  "time"
  "fmt"
	vlc "github.com/adrg/libvlc-go"
)

var currentSong Song
var SongLength int64
var StartedAt int64
var sendSong chan Song
var SendFinished chan Song

func setSongLength (s int64, t int64) {
	SongLength = t
	StartedAt = s
}

func PlaySong(s Song) {
	if !IsSongCached(s) {
		CacheSong(s, 0)
	} else {
		fmt.Println("\n\n===\nSending", s.Title, "to VLC!")
		currentSong = s
		sendSong <- s
	}
}

func songLoop(sS *chan Song, sF *chan Song) {
	fmt.Println("VLC Init")
	if err := vlc.Init("--no-video", "--quiet"); err != nil {
			log.Fatal(err)
	}
	defer vlc.Release()

  fmt.Println("VLC Player Init")

	// Create a new player.
	player, err := vlc.NewPlayer()
	if err != nil {
			log.Fatal(err)
	}
	defer func() {
			player.Stop()
			player.Release()
	}()

	fmt.Println("VLC Player Ready")

	for s := range *sS {

		fmt.Println("Recieved and Loading Song:", s.Title)
	  media, err := player.LoadMediaFromURL("file://" + GetSongPath(s))
	  if err != nil {
	      log.Fatal(err)
	  }
	  defer media.Release()

		fmt.Println("Play Song:", s.Title)
	  // Start playing the media.
		player.Stop()
	  err = player.Play()
	  if err != nil {
	      log.Fatal(err)
	  }

		go func() {
			// Weird Delay to fix media lengths being seen as 0
			time.Sleep(1 * time.Second)

		  length, err := player.MediaLength()
		  if err != nil || length == 0 {
		    length = 1000 * 60
		  }
			setSongLength(time.Now().Unix(), int64(length))
			fmt.Println("Song:", s.Title, "Length:", length)
		  time.Sleep(time.Duration(length) * time.Millisecond)
			fmt.Println("Finished Playing:", s.Title)
			*sF <- s
		}()
	}
}

func InitPlayer(){
	sS := make(chan Song)
	sendSong = sS
	sF := make(chan Song)
	SendFinished = sF
	go songLoop(&sendSong, &SendFinished)

  fmt.Println("Player Ready!")
}
