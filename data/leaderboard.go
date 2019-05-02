package data

import (
	"sort"
)

type Leaderboard []Entry

type Entry struct {
	User  string `json"user"`
	Score int    `json:"score"`
}

func GetLeaderboard(path string, limit int) (Leaderboard, error) {
	games, err := LoadGamesFromFile(path)
	if err != nil {
		return []Entry{}, err
	}

	// Add all the entries to a map (easier for adding score totals)
	m := make(map[string]int)
	for _, g := range games {
		if g.Solved && g.Solver != nil {
			m[*g.Solver] = m[*g.Solver] + 1
		}
	}

	// Pull our user and their score totals into a formal struct we can return
	lb := make([]Entry, len(m))
	var i int
	for k, v := range m {
		lb[i] = Entry{
			k,
			v,
		}
		i++
	}

	// Sort our slice by the highest scoring users
	sort.Slice(lb, func(i, j int) bool {
		return lb[i].Score > lb[j].Score
	})

	if len(lb) > limit {
		lb = lb[:limit]
	}

	return lb, nil
}
