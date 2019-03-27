var app;
console.log("Hey!");
window.addEventListener("load", function(event) {
  console.log("Loaded!");
  app = new Vue({
    el: '#app',
    data: {
      appMeta: {'name':'MC Hawk Juke Box',
                'version':'v0-prealpha'},
      player: {
        isPlaying: false,
        skipVotes: 0,
        numberOfListeners: 5,
        playlist: [new Song("Habu - Exit", "ctwNCX3uG_4"), new Song("Travis Wes - Sicko Bomba [CLEAN]", ""), new Song("Sheck Scott - Mo Mode (Big Room Biology Remix) [CLEAN]", "")],
        currentSong: new Song("Habu - Exit", "ctwNCX3uG_4"),
      }
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
    }

    upvote() {
      this.votes++;
      app.player.playlist.sort();
      // TODO: Actually contact server
    }

    downvote() {
      this.votes--;
      app.player.playlist.sort();
      // TODO: Actually contact server
    }
}
