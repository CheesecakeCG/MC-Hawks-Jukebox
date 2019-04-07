package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"html"
	"log"
	"io/ioutil"
)

type vote struct {
  Sender string
  SongTitle string
  VotesStatus int
}

type song struct {
  Title string
  Youtubeid string
  Votes []vote
}

func countVotes(s song)  int {
	var count int
	for _, v := range s.Votes {
		count += v.VotesStatus
	}
	return count
}

var songList []song

func getSong(t string) *song  {
	for _, s := range songList {
		if s.Title == t {
			return &s;
		}
	}
	return &song{"Not Found", "", []vote{}}
}

func castVote(v vote) {
	println("Vote cast by ", v.Sender, "on", v.SongTitle, "VotesStatus", v.VotesStatus)
	s := getSong(v.SongTitle)
	s.Votes = append(s.Votes, v)
}

func removeVote (toDelete vote) {
	s := getSong(toDelete.SongTitle)
	for i, v := range s.Votes {
		if v.Sender == toDelete.Sender {
			s.Votes = append(s.Votes[:i], s.Votes[i+1:]...)
		}
	}

}

func main()  {
  var testJson = "{\"Sender\":\"127.0.0.1\", \"SongTitle\":\"Habu - Exit\", \"VotesStatus\":-1}"
  testVote := parseVote(testJson)
  if testVote.Sender != "127.0.0.1" {
    fmt.Println("This vote is from the host!")
  }
  fmt.Println(testVote)
	fmt.Println("Songs", songList)

  http.HandleFunc("/vote", func(w http.ResponseWriter, r *http.Request) {
  	fmt.Fprintf(w, "Recieved vote", html.EscapeString(r.URL.Path))
		println("Voted Recieved!")
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		} else {
			println("No Don't Crash!", r.Body)
			castVote(parseVote(string(reqBody)))
		}

  })

	http.HandleFunc("/removeVote", func(w http.ResponseWriter, r *http.Request) {
  	fmt.Fprintf(w, "Recieved vote", html.EscapeString(r.URL.Path))
		println("RemoveVote Recieved!")
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		} else {
			println("No Don't Crash!", r.Body)
			removeVote(parseVote(string(reqBody)))
		}

  })

  log.Fatal(http.ListenAndServe(":8001", nil))
}



func parseVote(data string) (vote) {
  var out vote
  err := json.Unmarshal([]byte(data), &out)
  if err != nil {
    fmt.Println("ERROR! - ", err)
  }
  return out
}
