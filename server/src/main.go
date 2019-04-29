package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"html"
	"log"
	"io/ioutil"
	"strconv"
)

type Vote struct {
  Sender string
  SongTitle string
  VotesStatus int
}

type Song struct {
  Title string
  Youtubeid string
  Votes []Vote
}

type StrippedSong struct {
  Title string
  Youtubeid string
  Votes int
}

func CountVotes(s Song)  int {
	count := 0
	for _, v := range s.Votes {
		count += v.VotesStatus
	}
	return count
}

func getSongByTitle(t string) int  {
	for i, s := range SongList {
		if s.Title == t {
			return i
		}
	}
	return -1
}

func castVote(v Vote) {
	s := getSongByTitle(v.SongTitle)
	if s < 0 {
		fmt.Println("Vote on song,", v.SongTitle, "failed! Song does not exist!")
		return
	}
	for i, b := range SongList[s].Votes {
		if b.Sender == v.Sender {
			SongList[s].Votes = append(SongList[s].Votes[:i], SongList[s].Votes[i+1:]...)
		}
	}
	SongList[s].Votes = append(SongList[s].Votes, v)
	fmt.Println("Vote cast by ", v.Sender, "on", v.SongTitle, "VotesStatus", v.VotesStatus, "Vote List:", SongList[s].Votes)
	SaveSongList()
}

func pushSongList() {
	fmt.Println(SongList)
}

func main()  {


	// RESTful API

  http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {
		println("   --- New Vote Update! ---")
  	fmt.Fprintf(w, "Recieved vote", html.EscapeString(r.URL.Path))
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		castVote(parseVote(string(reqBody)))
		println("   ---\n")
  })

	http.HandleFunc("/addsong", func(w http.ResponseWriter, r *http.Request) {
		println("   --- Song Recieved! ---")
		fmt.Fprintf(w, "Song Recieved", html.EscapeString(r.URL.Path))
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		AddSong(parseSong(string(reqBody)))
		println("   ---\n")
	})

	http.HandleFunc("/currentSong", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{\"Title\":\"" + currentSong.Title + "\", \"SongLength\":\""+ strconv.FormatInt(SongLength, 10) + "\", \"StartedAt\":\""+ strconv.FormatInt(StartedAt, 10) + "\"}")
	})

	http.HandleFunc("/songList", func(w http.ResponseWriter, r *http.Request) {
		j, err := json.Marshal(StripSongListSender())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, string(j))
	})

	InitPlaylist()
	InitPlayer()
	InitAdmin()


	fmt.Println("\n                                 __                               _     _      _       \n |_|     /\\     \\    /    |/    (_             |    | |    |/    |_    |_)    / \\    \\/\n | |    /--\\     \\/\\/     |\\    __)          \\_|    |_|    |\\    |_    |_)    \\_/    /\\")
	fmt.Println("                        === [Initialization Complete!] ===\n\n")

	// Host Website
	fmt.Println("Starting Website!")
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

 	defer log.Fatal(http.ListenAndServe(":8000", nil))


}

func StripSongListSender() ([]StrippedSong){
	var out []StrippedSong
	for _,s := range SongList {
		var ss StrippedSong
		ss.Title = s.Title
		ss.Youtubeid = s.Youtubeid
		ss.Votes = CountVotes(s)
		out = append(out, ss)
	}
	return out;
}

func parseSong(data string) (Song) {
  var out Song
  err := json.Unmarshal([]byte(data), &out)
  if err != nil {
    fmt.Println("ERROR! - ", err)
  }
  return out
}

func parseVote(data string) (Vote) {
  var out Vote
  err := json.Unmarshal([]byte(data), &out)
  if err != nil {
    fmt.Println("ERROR! - ", err)
  }
  return out
}
