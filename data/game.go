package data

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

var (
	ErrGameNotFound = errors.New("game id not found")
)

type Games []Game

type Game struct {
	Solver   *string  `json:"solver,omitempty"`
	Solution *string  `json:"solution,omitempty"`
	Solved   bool     `json:"solved"`
	GameType string   `json:"game-type"`
	GameID   string   `json:"id"`
	Message  string   `json:"message"`
	Choices  []Choice `json:"choices"`
}

// Create a new choices struct; we probably could get away with just using
// employee.Headshot, but this would mean someone could fairly easily develop a
// programatic way of creating new games and solving by looking at the alt text
type Choice struct {
	MimeType string `json:"mimeType"`
	URL      string `json:"url"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
	AnswerID string `json:"answer-id"`
}

func NewGame(e Employees, gameType string) (Game, error) {
	if gameType != "standard" && gameType != "matt" {
		gameType = "standard"
	}

	var prefix string
	if gameType == "matt" {
		// Use mat to catch Matthew/Matt/Mat
		prefix = "Mat"
	}

	g := Game{
		GameType: gameType,
		Choices:  make([]Choice, 6),
	}

	var solutionEmployee *Employee

	for i := 0; i < 6; i++ {
		employee, err := GetRandomEmployee(e, prefix)
		if err != nil {
			return g, err
		}

		if solutionEmployee == nil {
			solutionEmployee = &employee
		}

		newChoice := EmployeeToChoice(employee.Headshot)
		newChoice.AnswerID = employee.ID
		g.Choices[i] = newChoice
	}

	g.GameID = GenerateGameURI()
	g.Message = fmt.Sprintf("Which image is a picture of %s %s?", solutionEmployee.FirstName, solutionEmployee.LastName)
	g.Solution = &solutionEmployee.ID

	return g, nil
}

func UpdateSolver(gameID string, solver string) error {
	games, err := LoadGamesFromFile(GamePath)
	if err != nil {
		return err
	}

	index := FindGameIndex(games, gameID)
	if index < 0 {
		return err
	}

	games[index].Solver = &solver
	games[index].Solved = true

	err = WriteGamesToFile(games, GamePath)
	if err != nil {
		return err
	}

	return nil
}

// Helper function that lets us easily convert an employee headshot into a
// stripped down struct with only certain fields available (we don't want to
// expose all the fields in the employee headshot to the user)
func EmployeeToChoice(eh EmployeeHeadshot) Choice {
	return Choice{
		URL:      eh.URL,
		Height:   eh.Height,
		Width:    eh.Width,
		MimeType: eh.MimeType,
	}
}

func IsCorrectSolution(gameID string, gameSolution string) (bool, error) {
	games, err := LoadGamesFromFile(GamePath)
	if err != nil {
		return false, err
	}

	index := FindGameIndex(games, gameID)
	if index < 0 {
		return false, ErrGameNotFound
	}

	solution := games[index].Solution
	if solution == nil {
		return false, errors.New("no solution recorded")
	}

	return *games[index].Solution == gameSolution, nil
}

func GetGameDetails(gameID string) (Game, error) {
	games, err := LoadGamesFromFile(GamePath)
	if err != nil {
		return Game{}, err
	}

	index := FindGameIndex(games, gameID)
	if index < 0 {
		return Game{}, ErrGameNotFound
	}

	return games[index], nil
}

// Generate a random game URI by hashing the current time
func GenerateGameURI() string {
	hasher := md5.New()
	hasher.Write([]byte(time.Now().String()))
	return hex.EncodeToString(hasher.Sum(nil))
}

type GameFilter struct {
	GameType string
	Solver   string
	GameID   string
}

func (f GameFilter) IsEmpty() bool {
	return f.GameType == "" && f.Solver == "" && f.GameID == ""
}

// Allow us to filter games by various facets. Any gamee that matches
// -any- of the filter criteria will be added
func (g Games) FilterGames(f GameFilter) Games {
	// If no filter is being used, we can just return the full list
	if f.IsEmpty() {
		return g
	}

	filterList := make(Games, 0, len(g))

	for _, game := range g {
		if f.GameID != "" && game.GameID == f.GameID {
			filterList = append(filterList, game)
		} else if f.Solver != "" && game.Solver != nil && *game.Solver == f.Solver {
			filterList = append(filterList, game)
		} else if f.GameType != "" && game.GameType == f.GameType {
			filterList = append(filterList, game)
		}
	}

	return filterList
}
