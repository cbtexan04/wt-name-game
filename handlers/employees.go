package handlers

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/cbtexan04/wt-name-game/data"
)

var (
	ErrMarshalling = errors.New("Unable to marshal employee json")
	ErrInvalidURL  = errors.New("Invalid url")
)

var (
	ListEmployeesRE  = regexp.MustCompile("/api/1.0/employees$")
	EmployeeDetailRE = regexp.MustCompile("/api/1.0/employees/([^/]+)$")
)

func GetEmployees(employees data.Employees) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := data.Filter{
			ID:        r.URL.Query().Get("id"),
			FirstName: r.URL.Query().Get("firstName"),
			LastName:  r.URL.Query().Get("lastName"),
		}

		filteredEmployees := employees.FilterEmployees(f)

		Write(w, http.StatusOK, filteredEmployees)
	}
}

func GetEmployeeDetails(employees data.Employees) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getIDFromRequest(EmployeeDetailRE, r.URL.Path)
		if err != nil {
			Error(w, http.StatusInternalServerError, ErrInvalidURL.Error())
			return
		}

		f := data.Filter{ID: id}
		filteredEmployees := employees.FilterEmployees(f)

		Write(w, http.StatusOK, filteredEmployees)
	}
}

func getIDFromRequest(re *regexp.Regexp, url string) (id string, err error) {
	matches := re.FindStringSubmatch(url)
	if len(matches) != 2 {
		return "", ErrInvalidURL
	}

	return matches[1], nil
}
