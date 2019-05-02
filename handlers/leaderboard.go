package handlers

import (
	"net/http"
	"regexp"

	"github.com/cbtexan04/wt-test-project/data"
)

var LeaderboardRE = regexp.MustCompile("/api/1.0/leaderboard$")

func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	games, err := data.LoadGamesFromFile(data.GamePath)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	m := make(map[string]int)
	for _, g := range games {
		if g.Solved {
			m[g.Solver] = m[g.Solver] + 1
		}
	}

	Write(w, http.StatusOK, m)
}
