package rapid

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"testing"
)

type RapidClientTestSuite struct {
	suite.Suite
	mockController *gomock.Controller
	client         RapidClient
}

func TestRapidClientTestSuite(b *testing.T) {
	suite.Run(b, new(RapidClientTestSuite))
}

func (suite *RapidClientTestSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.client = NewRapidClient(&http.Client{})
}

func (suite *RapidClientTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

type mockRoundTripper struct {
	mockFunc func(req *http.Request) *http.Response
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.mockFunc(req), nil
}

func newTestClient(mockFunc func(req *http.Request) *http.Response) *http.Client {
	return &http.Client{
		Transport: &mockRoundTripper{mockFunc},
	}
}

func (suite *RapidClientTestSuite) TestFetchAllMovies_Success() {
	jsonResponse := `[{"id":1,"title":"Movie 1"},{"id":2,"title":"Movie 2"}]`

	mockClient := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(jsonResponse)),
			Header:     make(http.Header),
		}
	})

	rc := NewRapidClient(mockClient)

	movies, err := rc.FetchAllMovies()

	suite.Nil(err)
	suite.Len(movies, 2)
}

func (suite *RapidClientTestSuite) TestFetchAllMovies_HTTPError() {
	mockClient := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(bytes.NewBufferString("Internal Server Error")),
			Header:     make(http.Header),
		}
	})

	rc := NewRapidClient(mockClient)

	_, err := rc.FetchAllMovies()

	suite.NotNil(err)
}

func (suite *RapidClientTestSuite) TestFetchAllMovies_InvalidJSON() {
	invalidJSON := `[{id:1,title:"Movie 1"}]`

	mockClient := newTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(invalidJSON)),
			Header:     make(http.Header),
		}
	})

	rc := NewRapidClient(mockClient)

	_, err := rc.FetchAllMovies()

	suite.NotNil(err)
}
