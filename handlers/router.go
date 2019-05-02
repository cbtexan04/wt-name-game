package handlers

import (
	"net/http"
	"regexp"
	"strings"
)

// Unfortunately, go does not have built in regexp matching for http
// handlers. We'll have to build own own to be able to use a regex with
// our handlers. While we're at it, we can also make our custom router
// accept which http methods apply on a given regexp

// Route associates a regular expression with an http.Handler
type Route struct {
	url     *regexp.Regexp
	h       http.Handler
	methods []string
}

// RegexpHandler handles the routing of URL patterns onto handlers by
// means of regular expressions
type RegexpHandler struct {
	Routes []Route
}

// Handle registers an http.Handler to a regular expression
func (rr *RegexpHandler) Handle(pattern string, h http.Handler, m ...string) {
	re := regexp.MustCompile(pattern)

	r := Route{
		url:     re,
		h:       h,
		methods: m,
	}

	rr.Routes = append(rr.Routes, r)
}

// HandleFunc registers an http.HandlerFunc to a regular expression
func (rr *RegexpHandler) HandleFunc(reg *regexp.Regexp, h http.HandlerFunc, m ...string) {
	re := regexp.MustCompile(reg.String())

	r := Route{
		url:     re,
		h:       http.HandlerFunc(h),
		methods: m,
	}

	rr.Routes = append(rr.Routes, r)
}

// ServeHTTP compares the URL associated with the http.Request to each
// Route the RegexpHandler knows about to find the appropriate handler for the
// request. It returns a 404 Not Found if unable to find a match.
func (rr *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range rr.Routes {
		if route.url.MatchString(r.URL.Path) {
			// If no route specified, assume open to all methods
			if len(route.methods) == 0 || ContainsMethod(route.methods, r.Method) {
				route.h.ServeHTTP(w, r)
				return
			}
		}
	}

	http.NotFound(w, r)
}

func ContainsMethod(slice []string, method string) bool {
	for _, v := range slice {
		if strings.ToUpper(v) == method {
			return true
		}
	}

	return false
}
