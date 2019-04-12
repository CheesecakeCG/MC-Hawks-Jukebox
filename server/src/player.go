package main

import (
	"log"
	"os"
  "time"
  "fmt"
	"github.com/faiface/beep"
  // "github.com/faiface/beep/mp3"
	"github.com/faiface/beep/flac"
  "github.com/faiface/beep/speaker"
)


func PlaySong(s Song) {

  fmt.Println("Loading:", s.Title)

	if !IsSongCached(s) {
		CacheSong(s, 0)
		return
	}
  // f, err := os.Open("/home/chrisd/.cache/mch-jukebox/Covox - Switchblade Squadron.mp3") // Works
  f, err := os.Open(GetSongPath(s))
	if err != nil {
		log.Fatal(err)
	}

  fmt.Println("Streeammer Init:", s.Title)
	streamer, format, err := flac.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
  defer streamer.Close()

  fmt.Println("Speaker Init:", s.Title)
  // ReIniting stops all other sounds
  speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

  fmt.Println("Now Playing:", s.Title)
  done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
    fmt.Println("Finished Playing:", s.Title)
    done <- true
	})))

	<-done
}

func InitPlayer(){
  fmt.Println("Player Ready!")
}
