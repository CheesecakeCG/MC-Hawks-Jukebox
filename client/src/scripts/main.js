var host = "http://" + window.location.hostname;
if (window.location.port != null) {
  host += ":" + window.location.port;
}

if (localStorage.getItem("playlist") == undefined) {
  localStorage.setItem("playlist", "[]");
}

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
        playlist: JSON.parse(localStorage.getItem("playlist")),
        currentSong: new Song({Title:"Connecting to Jukebox", Youtubeid:""}),
        songLength: 0
      },
      showFabMenu: false,
      skipEnabled: false
    },
    methods: {
      voteForSkip: function (event) {
        app.player.skipVotes++;
      },
      refreshCurrentSong: function () {
        app.player.currentSong = getCurrentSong();
        setInterval(function() {
          app.player.currentSong = getCurrentSong();
        }, 4000);
      },
      refreshSongList: function () {
        getSongList();
        setInterval(function() {
          getSongList();
        }, 15000);
      },
      getSeekAsDec: function () {
        var end = (parseInt(this.player.currentSong.SongLength) +  parseInt(this.player.currentSong.StartedAt) * 1000);
        var time = Date.now();
        return (1 - ((end - time) / this.player.currentSong.SongLength)) * 100 + "%" ;
      }
    }
  })
  app.refreshCurrentSong()
  app.refreshSongList()
});

function getCurrentSong(){
  return getFromHost("/currentSong");
}

function getSongList(){
  var sl = getFromHost("/songList");

  for (var i = 0; i < sl.length; i++) {
    var si = app.player.playlist.findIndex(function(e) {
      return e.Title == sl[i].Title;
    });
    sl[i] = new Song(sl[i]);
    if (si >= 0 && app.player.currentSong.Title != sl[i].Title) {
      sl[i].votesStatus = app.player.playlist[si].votesStatus;
    }
  }
  app.player.playlist = sl.reverse();
}

function getFromHost(path) {
  var xhr = new XMLHttpRequest();
  xhr.open("GET", host + path, false);
  xhr.send(null);
  return JSON.parse(xhr.responseText);
}

class Song {
    // constructor(title, youtubeid) {
    //     this.Title = title;
    //     this.Youtubeid = youtubeid;
    //     this.votes = 0;
    //     this.votesStatus = 0;
    // }

    constructor(obj) {
        this.Title = obj.Title;
        this.Youtubeid = obj.Youtubeid;
        this.votes = obj.Votes;
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
      this.updateVote();
    }

    removeVote() {
      this.votesStatus = 0;
      this.updateVote();
    }

    updateVote() {
      if (app.player.currentSong.Title != this.Title){
        sendVote(this.Title, this.votesStatus);
      } else {
        this.removeVote();
        return;
      }
      setTimeout(getSongList(), 400);
      localStorage.setItem("playlist", JSON.stringify(app.player.playlist));
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
