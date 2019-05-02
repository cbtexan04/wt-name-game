package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/cbtexan04/wt-test-project/data"
)

var (
	GameRE        = regexp.MustCompile("/api/1.0/games$")
	GameDetailsRE = regexp.MustCompile("/api/1.0/games/([^/]+)$")
)

func GetGameDetails(w http.ResponseWriter, r *http.Request) {
	gid, err := getIDFromRequest(GameDetailsRE, r.URL.Path)
	if err != nil {
		http.Error(w, ErrInvalidURL.Error(), http.StatusInternalServerError)
		return
	}

	game, err := data.GetGameDetails(gid)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}

	Write(w, http.StatusOK, stripSolutions(game))
}

func DeleteGame(w http.ResponseWriter, r *http.Request) {
	gid, err := getIDFromRequest(GameDetailsRE, r.URL.Path)
	if err != nil {
		http.Error(w, ErrInvalidURL.Error(), http.StatusInternalServerError)
		return
	}

	games, err := data.LoadGamesFromFile(data.GamePath)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	games = data.DeleteGame(games, gid)

	err = data.WriteGamesToFile(games, data.GamePath)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetGames(w http.ResponseWriter, r *http.Request) {
	games, err := data.LoadGamesFromFile(data.GamePath)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	Write(w, http.StatusOK, stripSolutions(games...))
}

func NewGame(employees data.Employees) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		g, err := data.NewGame(employees, r.URL.Query().Get("game-type"))
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Use local storage to save the game details.. In reality we would
		// want to use a database, but for time considerations just write it
		// out to a file
		games, err := data.LoadGamesFromFile(data.GamePath)
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		games = append(games, g)

		err = data.WriteGamesToFile(games, data.GamePath)
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		new := stripSolutions(g)
		fmt.Println(new[0].Solution)

		Write(w, http.StatusOK, stripSolutions(g))
	}
}

type Solution struct {
	Solution string `json:"solution"`
}

func CheckSolution(employees data.Employees) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gid, err := getIDFromRequest(GameDetailsRE, r.URL.Path)
		if err != nil {
			http.Error(w, ErrInvalidURL.Error(), http.StatusInternalServerError)
			return
		}

		var solution Solution
		err = json.NewDecoder(r.Body).Decode(&solution)
		if err != nil {
			Error(w, http.StatusBadRequest, fmt.Sprintf("Unexpected json response: %s", err))
			return
		}

		solved, err := data.IsCorrectSolution(gid, solution.Solution)
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		} else if !solved {
			Error(w, http.StatusBadRequest, "Invalid game solution")
			return
		}

		user, _, _ := r.BasicAuth()
		err = data.UpdateSolver(gid, user)
		if err != nil {
			Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		response := struct {
			Message string
		}{
			"You got the right answer!",
		}

		Write(w, http.StatusOK, response)
	}
}

// Strip the solutions from unsolved games
func stripSolutions(games ...data.Game) data.Games {
	newGames := make([]data.Game, len(games))
	for i, game := range games {
		newGames[i] = game
		if !game.Solved {
			newGames[i].Solution = nil
		}
	}

	return newGames
}

func getGameID(re *regexp.Regexp, url string) (id string, err error) {
	matches := re.FindStringSubmatch(url)
	if len(matches) != 2 {
		return "", ErrInvalidURL
	}

	return matches[1], nil
}