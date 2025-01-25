package main

import (
	"bytes"
	"io"
	"log"
	"html"
	"net/http"
	"net/url"
	"net/http/httptest"
	"testing"
	"net/http/cookiejar"

	"time"
	"regexp" 
	"snippetbox.alexedwards.net/internal/models/mocks"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"

	// "fmt"
)

// // Define a regular expression which captures the CSRF token value from the HTML for our user signup page.
var csrfTokenRX = regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+?)"`)


func extractCSRFToken(t *testing.T, body string) string {
	// FindStringSubmatch returns a slice of strings 
	// containing the text of the leftmost match and the matches of the subexpressions
	matches := csrfTokenRX.FindStringSubmatch(body)

	// fmt.Println("matches: ", matches)
    
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}
    
	// return the second element of the matches slice, which is the token value.
	return html.UnescapeString(matches[1])
}


// Create a newTestApplication helper which returns an instance of our application struct containing mocked dependencies.
func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
	  t.Fatal(err)
	}

	formDecoder := form.NewDecoder()
    
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
    
	return &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog: log.New(io.Discard, "", 0),
		snippets: &mocks.SnippetModel{},
		users: &mocks.UserModel{},
		templateCache: templateCache,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	}
 
}

// Define a custom testServer type which embeds a httptest.Server instance.
type testServer struct {
   *httptest.Server
}

// Create a newTestServer helper which initalizes and returns a new instance of our custom testServer type
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// Create a new cookie jar and set the client's Jar field to it.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the client's Jar field to the cookie jar we created above. This will store cookies between requests.
	ts.Client().Jar = jar

	// Disable redirect-following for the client. 
	// This means that if the server sends a 3xx response status code, 
	// the client will return the response to the caller rather than following the redirect.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}


// Implement a get() method on our custom testServer type. This makes a GET
// request to a given url path using the test server client, and returns the response status code, headers and body.

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
    rs, err := ts.Client().Get(ts.URL + urlPath)
    if err != nil {
        t.Fatal(err)
    }

    defer rs.Body.Close()
    body, err := io.ReadAll(rs.Body)
    if err != nil {
        t.Fatal(err)
    }

    // Log the response body for debugging
    // t.Logf("Response Body: %s", body)

    bytes.TrimSpace(body)
    return rs.StatusCode, rs.Header, string(body)
}


func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	rs, err := ts.Client().PostForm(ts.URL + urlPath, form)
	if err != nil {
		t.Fatal(err)
	}
	
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	
	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}
