package data

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
)

var (
	ErrGameNotFound = errors.New("game id not found")
)

type Games []Game

type Game struct {
	Solver   string             `json:"solver,omitempty"`
	Solved   bool               `json:"solved"`
	GameType string             `json:"game-type"`
	GameID   string             `json:"id"`
	Message  string             `json:"message"`
	Choices  []EmployeeHeadshot `json:"employees"`
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
		Choices:  make([]EmployeeHeadshot, 6),
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

		employee.Headshot.EmployeeID = &employee.ID
		g.Choices[i] = employee.Headshot
	}

	g.GameID = GameSolutionID(solutionEmployee.ID)
	g.Message = fmt.Sprintf("Which image is a picture of %s %s?", solutionEmployee.FirstName, solutionEmployee.LastName)

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

	games[index].Solver = solver
	games[index].Solved = true

	err = WriteGamesToFile(games, GamePath)
	if err != nil {
		return err
	}

	return nil
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

// We can cheat and use the game ID as a hashed version of the solution. This
// makes it nice because we can check if a user guessed a solution by simply
// hashing their solution and seeing if it's equivalent to our game URI. This
// *does* come with the drawback that eventually someone is going to end up
// with the same game ID URI, but since this is a proof of concept, it should
// suffice
func GameSolutionID(solutionID string) string {
	hasher := md5.New()
	hasher.Write([]byte(solutionID))
	return hex.EncodeToString(hasher.Sum(nil))
}
