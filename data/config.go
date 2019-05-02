package data

import (
	"encoding/json"
	"os"
)

var GamePath = "/tmp/gamesConfig.json"

func LoadGamesFromFile(path string) (Games, error) {
	if _, err := os.Stat(path); err != nil {
		return Games{}, nil
	}

	fd, err := os.Open(path)
	if err != nil {
		return Games{}, err
	}
	defer fd.Close()

	var games Games
	err = json.NewDecoder(fd).Decode(&games)
	return games, err
}

func WriteGamesToFile(g Games, path string) error {
	var fd *os.File
	var err error
	if fd, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0666); err != nil {
		fd, err = os.Create(path)
		if err != nil {
			return err
		}
	}
	defer fd.Close()

	enc := json.NewEncoder(fd)
	enc.SetIndent("", "    ")
	return enc.Encode(&g)
}

func FindGameIndex(games Games, gameID string) int {
	for i, g := range games {
		if g.GameID == gameID {
			return i
		}
	}

	return -1
}

func DeleteGame(games Games, gameID string) Games {
	for i, game := range games {
		if game.GameID == gameID {
			return append(games[:i], games[i+1:]...)
		}
	}

	return games
}
