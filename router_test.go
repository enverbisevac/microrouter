package microrouter

import (
	"net/http"
	"testing"
)

func TestNewRegexResolver(t *testing.T) {
	object := newRegexResolver()
	if object.handlers == nil {
		t.Error("handlers not created")
	}
	if object.cache == nil {
		t.Error("cache not created")
	}
	if object.notFound == nil {
		t.Error("not found handlers not created")
	}
	if object.methodNotFound == nil {
		t.Error("method handlers not created")
	}
}

func TestCheckMethod(t *testing.T) {
	cases := map[string]struct {
		inputMethod string
		inputPath   string
		pattern     string
		exp         uint8
	}{
		"Test empty method":                    {"", "/hello", "GET /hello", methodNotFound},
		"Test path":                            {"GET", "/hello", "GET /$", pathNotFound},
		"Testing POST method with GET request": {"GET", "/hello", "POST /hello", methodNotFound},
		"Test method and path":                 {"GET", "/hello", "GET /hello", methodAndPathFound},
		"Test GET method":                      {"GET", "/hello", "(GET|POST) /hello", methodAndPathFound},
		"Test POST method":                     {"POST", "/hello", "(GET|POST) /hello", methodAndPathFound},
		"Test OPTIONS method":                  {"OPTIONS", "/hello", "(GET|POST) /hello", methodAndPathFound},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := checkMethod(tc.inputMethod, tc.inputPath, tc.pattern)
			if got != tc.exp {
				t.Errorf("Expected result %d, but we got %d", tc.exp, got)
			}
		})
	}
}

func TestGetContentType(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	req.Header.Set("Content-Type", "")
	got := getContentType(req)
	exp := defaultContentType
	if got != exp {
		t.Errorf("Expected content type %s, but we got %s", exp, got)
	}

	req.Header.Set("Content-Type", "application/json")
	got = getContentType(req)
	exp = "application/json"
	if got != exp {
		t.Errorf("Expected content type %s, but we got %s", exp, got)
	}
}

func TestGenerateFullPattern(t *testing.T) {
	cases := map[string]struct {
		pattern string
		methods []string
		exp     string
	}{
		"empty pattern and method": {"", []string{}, "GET /$"},
		"empty pattern":            {"", []string{"GET"}, "GET /$"},
		"pattern 1":                {"/$", []string{"GET"}, "GET /$"},
		"pattern 2":                {"/hello", []string{"GET"}, "GET /hello"},
		"pattern 3":                {"/hello", []string{"GET", "POST"}, "(GET|POST) /hello"},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := generateFullPattern(tc.pattern, tc.methods...)
			if got != tc.exp {
				t.Errorf("Expected content type %s, but we got %s", tc.exp, got)
			}
		})
	}
}
