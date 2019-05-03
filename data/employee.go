package data

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strings"
)

var ErrNoEmployeeMatch = errors.New("Unable to find employee with specified prefix")

type Employees []Employee

type EmployeeHeadshot struct {
	Type     string `json:"type"`
	MimeType string `json:"mimeType"`
	ID       string `json:"id"`
	URL      string `json:"url"`
	Alt      string `json:"alt"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
}

type Employee struct {
	ID          string           `json:"id"`
	Type        string           `json:"type"`
	Slug        string           `json:"slug"`
	JobTitle    string           `json:"jobTitle,omitempty"`
	FirstName   string           `json:"firstName"`
	LastName    string           `json:"lastName"`
	Headshot    EmployeeHeadshot `json:"headshot"`
	SocialLinks []interface{}    `json:"socialLinks"`
	Bio         string           `json:"bio,omitempty"`
}

func FetchEmployees(remoteUrl string) ([]Employee, error) {
	var employees Employees

	rs, err := http.Get(remoteUrl)
	if err != nil {
		return employees, err
	}

	err = json.NewDecoder(rs.Body).Decode(&employees)
	return employees, err
}

func GetRandomEmployee(employees Employees, prefix string) (Employee, error) {
	for i := 0; i < len(employees); i++ {
		e := employees[rand.Intn(len(employees))]
		if strings.HasPrefix(e.FirstName, prefix) {
			return e, nil
		}
	}

	// If we ran through everyone, and no one has the requested name, error
	return Employee{}, ErrNoEmployeeMatch
}

type Filter struct {
	FirstName string
	LastName  string
	ID        string
}

func (f Filter) IsEmpty() bool {
	return f.FirstName == "" && f.LastName == "" && f.ID == ""
}

// Allow us to filter employees by various facets. Any employee that matches
// -any- of the filter criteria will be added
func (e Employees) FilterEmployees(f Filter) Employees {
	// If no filter is being used, we can just return the full list
	if f.IsEmpty() {
		return e
	}

	filterList := make(Employees, 0, len(e))

	for _, employee := range e {
		if f.ID != "" && employee.ID == f.ID {
			filterList = append(filterList, employee)
		} else if f.FirstName != "" && employee.FirstName == f.FirstName {
			filterList = append(filterList, employee)
		} else if f.LastName != "" && employee.LastName == f.LastName {
			filterList = append(filterList, employee)
		}
	}

	return filterList
}
