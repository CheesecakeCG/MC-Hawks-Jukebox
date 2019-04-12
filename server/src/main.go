package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"html"
	"log"
	"io/ioutil"
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

func countVotes(s Song)  int {
	var count int
	for _, v := range s.Votes {
		count += v.VotesStatus
	}
	return count
}

var SongList []Song

func getSong(t string) *Song  {
	for _, s := range SongList {
		if s.Title == t {
			return &s;
		}
	}
	return &Song{"Not Found", "", []Vote{}}
}

func castVote(v Vote) {

	s := getSong(v.SongTitle)
	s.Votes = append(s.Votes, v)
	println("Vote cast by ", v.Sender, "on", v.SongTitle, "VotesStatus", v.VotesStatus, "Vote List:", s.Votes)
}

func pushSongList() {
	fmt.Println(SongList)
}

func main()  {
  var testJson = "{\"Sender\":\"127.0.0.1\", \"SongTitle\":\"Habu - Exit\", \"VotesStatus\":-1}"
  testVote := parseVote(testJson)
  if testVote.Sender != "127.0.0.1" {
    fmt.Println("This Vote is from the host!")
  }
  fmt.Println(testVote)
	fmt.Println("Songs", SongList)

	// Host Website
	fmt.Println("Starting Website!")
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// RESTful API

  http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {

		println("   --- New Vote Update! ---")
  	fmt.Fprintf(w, "Recieved vote", html.EscapeString(r.URL.Path))
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		} else {
			castVote(parseVote(string(reqBody)))
		}
		println("   ---\n")
  })

	InitAdmin()
	InitPlayer()
	
  log.Fatal(http.ListenAndServe(":8000", nil))
}


func parseVote(data string) (Vote) {
  var out Vote
  err := json.Unmarshal([]byte(data), &out)
  if err != nil {
    fmt.Println("ERROR! - ", err)
  }
  return out
}
