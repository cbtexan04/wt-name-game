package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/cbtexan04/wt-test-project/data"
	"github.com/cbtexan04/wt-test-project/handlers"
)

const (
	EmployeeDataUrl  = "https://www.willowtreeapps.com/api/v1.0/profiles"
	RequiredUser     = "willowtree"
	RequiredPassword = "password"
)

var (
	ErrFetchEmployees = errors.New("Cannot fetch employee data from remote host")
	ErrUnauthorized   = errors.New("User is not authorized")
)

var employees []data.Employee

// We will need to grab the data remotely from another host; for now, we can
// just grab it in the init, but it would be nice if we could poll for the
// employee data every now and then in case of updates to the employee list
func init() {
	e, err := data.FetchEmployees(EmployeeDataUrl)
	if err != nil {
		panic(fmt.Sprintf("%s - %s", ErrFetchEmployees, err))
	}

	employees = e
}

func main() {
	router := &handlers.RegexpHandler{
		Routes: make([]handlers.Route, 0, 10),
	}

	router.HandleFunc(handlers.GameRE, handlers.GetGames, "GET")
	router.HandleFunc(handlers.GameDetailsRE, handlers.DeleteGame, "DELETE")
	router.HandleFunc(handlers.GameDetailsRE, handlers.GetGameDetails, "GET")
	router.HandleFunc(handlers.GameDetailsRE, handlers.CheckSolution(employees), "PUT")
	router.HandleFunc(handlers.GameRE, handlers.NewGame(employees), "POST")
	router.HandleFunc(handlers.ListEmployeesRE, handlers.GetEmployees(employees), "GET")
	router.HandleFunc(handlers.EmployeeDetailRE, handlers.GetEmployeeDetails(employees), "GET")
	router.HandleFunc(handlers.LeaderboardRE, handlers.GetLeaderboard, "GET")
	log.Fatal(http.ListenAndServe(":1234", handleAuth(router)))
}

// We can use composition to handle our auth pretty easily. While normally we
// would want to make sure we're creating new sessions for our users, we can
// just check http basic auth and move along. For demonstration purposes we'll
// enforce a default username and password, but realistically we would have
// multiple users which could authenticate.
func handleAuth(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, pass, ok := r.BasicAuth()
		if !ok || pass != RequiredPassword {
			handlers.Error(w, http.StatusInternalServerError, ErrUnauthorized.Error())
			return
		}

		h.ServeHTTP(w, r)
	}
}
