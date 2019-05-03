package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/cbtexan04/wt-name-game/data"
)

var ErrInvalidLimit = errors.New("Invalid limit; must be >= 1")

var LeaderboardRE = regexp.MustCompile("/api/1.0/leaderboard$")

func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	// Limit our leaderboard to 10 unless the user tells us otherwise
	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		var err error
		limit, err = strconv.Atoi(l)
		if err != nil || limit < 1 {
			Error(w, http.StatusBadRequest, ErrInvalidLimit.Error())
			return
		}
	}

	lb, err := data.GetLeaderboard(data.GamePath, limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Write(w, http.StatusOK, lb)
}
