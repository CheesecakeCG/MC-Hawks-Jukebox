package main

import (
  "fmt"
  "github.com/nanobox-io/golang-scribble"
  "encoding/json"
  "sort"
)

type SongArray []Song //Needed for a custom sorting the songlist
var SongList []Song
var Db scribble.Driver
var DatabaseDir string = "/home/chrisd/.config/mch-jukebox/db"

func (s SongArray) Len() int {
    return len(s)
}
func (s SongArray) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s SongArray) Less(i, j int) bool {
    return CountVotes(s[i]) < CountVotes(s[j])
}

func SaveSongList() {
  SortSongList()
  for _, s := range SongList {
    if err := Db.Write("SongList", s.Title, s); err != nil {
      fmt.Println("Error", err)
    }
  }
}

func SortSongList()  {
  // fmt.Println("Sorting SongList!")
  sort.Sort(SongArray(SongList))
  // fmt.Println("Done Sorting:", SongList)
}

func AddSongToDB(s Song) {
  SongList = append(SongList, s)
  SaveSongList()
}

func LoadAllSongsFromDB()  {
  slRead, err := Db.ReadAll("SongList")
  if err != nil {
    fmt.Println("Error", err)
  }

  for _, s := range slRead {
    songFound := Song{}
    if err := json.Unmarshal([]byte(s), &songFound); err != nil {
      fmt.Println("Error", err)
    }
    SongList = append(SongList, songFound)
  }

}

func OnSongFinished(sF *chan Song, sl *[]Song)  {
  for s := range *sF {
    // if currentSong.Title != s.Title {
    //   continue
    // }
    fmt.Println("Recieved CB for Finishing Playing:", s.Title)
    (*sl)[getSongByTitle(s.Title)].Votes = []Vote{}
    SortSongList()
    if (*sl)[len(*sl) - 1].Title != currentSong.Title {
      PlaySong((*sl)[len(*sl) - 1])
    } else {
      fmt.Println("No better song to play, playing: ", getSongByTitle(currentSong.Title) - 1)
      PlaySong((*sl)[getSongByTitle(currentSong.Title) - 1])
    }

  }
}

func InitPlaylist() {
  db, err := scribble.New(DatabaseDir, nil)
  if err != nil {
    fmt.Println("Error: ", err)
  }
  Db = *db

  LoadAllSongsFromDB()

  fmt.Println("Songs Loaded: ", SongList)
  go OnSongFinished(&SendFinished, &SongList)
  go PlaySong(SongList[0],)
  fmt.Println("Playlist Ready!")
}
