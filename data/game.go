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
	Solver   *string            `json:"solver,omitempty"`
	Solution *string            `json:"solution,omitempty"`
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
