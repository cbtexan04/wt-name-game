package tests

import (
	"testing"

	"github.com/cbtexan04/wt-test-project/data"
	"github.com/cbtexan04/wt-test-project/handlers"
)

func TestStripSolutions(t *testing.T) {
	solution := "123456789"
	game := data.Game{
		Solution: &solution,
	}

	strippedGame := handlers.StripSolutionFromGame(game)
	if strippedGame.Solution != nil {
		t.Error("Game ID was not stripped from game struct")
	}
}

func TestGenerateURIGivesUniqueResults(t *testing.T) {
	m := make(map[string]bool)
	for i := 0; i < 50; i++ {
		uri := data.GenerateGameURI()
		if m[uri] == true {
			t.Error("URI generated was not unique")
		}

		m[uri] = true
	}
}

func TestFindGameIndex(t *testing.T) {
	// Populate the games struct
	games := make([]data.Game, 10)
	for i := 0; i < 10; i++ {
		uri := data.GenerateGameURI()

		g := data.Game{
			GameID: uri,
		}

		games[i] = g
	}

	// Check that the FindGameIndex function finds the right index
	for i := 0; i < 10; i++ {
		if data.FindGameIndex(games, games[i].GameID) != i {
			t.Fatal("Game index was incorrectly calculated")
		}
	}
}

func TestDeleteGame(t *testing.T) {
	// Populate the games struct
	games := make([]data.Game, 10)
	for i := 0; i < 10; i++ {
		uri := data.GenerateGameURI()

		g := data.Game{
			GameID: uri,
		}

		games[i] = g
	}

	// Check that the FindGameIndex function finds the right index
	for i := 0; i < len(games); i++ {
		// Delete the first one
		deleteID := games[0].GameID
		games = data.DeleteGame(games, deleteID)

		if contains(games, deleteID) {
			t.Fatal("Game index was not deleted")
		}
	}
}

func contains(games data.Games, gameID string) bool {
	for _, g := range games {
		if g.GameID == gameID {
			return true
		}
	}

	return false
}
