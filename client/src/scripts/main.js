var host = "localhost";

function uuidv1() {
  if (localStorage.getItem("voterId") == null) {
    var xhr = new XMLHttpRequest();
    xhr.open("GET", "https://www.uuidgenerator.net/api/version1", false);
    xhr.send(null);
    localStorage.setItem("voterId", xhr.responseText);
  }
  return localStorage.getItem("voterId");
}



var app;
console.log("Hey!");
window.addEventListener("load", function(event) {
  console.log("Loaded!");
  app = new Vue({
    el: '#app',
    data: {
      appMeta: {'name':'Hawks Jukebox',
                'version':'v0-prealpha'},
      player: {
        isPlaying: true,
        skipVotes: 0,
        numberOfListeners: 5,
        playlist: [new Song("Habu - Exit", "ctwNCX3uG_4"), new Song("Travis Wes - Sicko Bomba [CLEAN]", ""), new Song("Sheck Scott - Mo Mode (Big Room Biology Remix) [CLEAN]", "")],
        currentSong: new Song("Habu - Exit", "ctwNCX3uG_4"),
      },
      showFabMenu: false,
      skipEnabled: false
    },
    methods: {
      voteForSkip: function (event) {
        app.player.skipVotes++;
      }
    }
  })
});

class Song {
    constructor(title, youtubeid) {
        this.title = title;
        this.youtubeid = youtubeid;
        this.votes = 0;
        this.votesStatus = 0;
    }

    upvote() {
      if (this.votesStatus == 1) {
        this.removeVote()
        return
      }
      this.votesStatus = 1;
      this.updateVote()
    }

    downvote() {
      if (this.votesStatus == -1) {
        this.removeVote()
        return
      }
      this.votesStatus = -1;
      this.updateVote()
    }

    removeVote() {
      this.votesStatus = 0;
      this.updateVote()
    }

    updateVote() {
      sendVote(this.title, this.votesStatus);
    }

}

function sendVote(title, upvote){
  var xhr = new XMLHttpRequest();
  xhr.open("POST", host + "/vote", false);
  xhr.setRequestHeader('Content-Type', 'application/json');
  var j = JSON.stringify({
    Sender: uuidv1(),
    SongTitle: title,
    VotesStatus: upvote
  })
  console.log(j);
  xhr.send(j);
}
