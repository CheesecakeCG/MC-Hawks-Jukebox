package main

import (
	"log"
	"os"
  "github.com/BrianAllred/goydl"
	"fmt"
)

var cachePath string = "/home/chrisd/.cache/mch-jukebox/"
var youtubeDl = goydl.NewYoutubeDl()

func AddSongToDB(s Song)  {
  if _, err := os.Stat(GetSongPath(s)); err == nil {
    fmt.Println("Song,", s.Title, ", already cached!")
  } else if os.IsNotExist(err) {
		CacheSong(s, 0)
  } else {
    log.Fatal(err)
    return
  }
	SongList = append(SongList, s)
}

func GetSongPath(s Song) (string) {
	return cachePath + s.Youtubeid + ".mp3"
}

func CacheSong(s Song, attempts int) {
	fmt.Println("Downloading song,", s.Title, "to", GetSongPath(s))
	youtubeDl.Options.ExtractAudio.Value = true
	youtubeDl.Options.AudioFormat.Value = "mp3"
	youtubeDl.Options.Output.Value = GetSongPath(s)

	cmd, ytErr := youtubeDl.Download("https://www.youtube.com/watch?v=" + s.Youtubeid)

	if ytErr != nil {
		if attempts >= 2 {
			fmt.Println("Download Failed (3 tries) on ", s.Title)
			// log.Fatal(ytErr)
		} else {
			CacheSong(s, attempts + 1)
		}

		return
	}
	cmd.Wait()
}

func IsSongCached(s Song) (bool){
	if _, err := os.Stat(GetSongPath(s)); err == nil {
    return true
  } else if os.IsNotExist(err) {
		return false
	} else {
    log.Panic(err)
    return false
  }
}


func InitAdmin()  {

  // go io.Copy(os.Stdout, youtubeDl.Stdout)
  // go io.Copy(os.Stderr, youtubeDl.Stderr)

  AddSongToDB(Song{"Habu - Exit", "ctwNCX3uG_4", []Vote{}})
  AddSongToDB(Song{"Days Of Thunder", "UiSB2Fbw9gs", []Vote{}})
	//AddSongToDB(Song{"Covox - Switchblade Squadron", "Covox - Switchblade Squadron", []Vote{}})
	PlaySong(SongList[1])
	fmt.Println("Admin Ready!")
}
