package repository

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MovieRepositoryTestSuite struct {
	suite.Suite
	db             *sqlx.DB
	mockedDB       *sqlx.DB
	mockDB         sqlmock.Sqlmock
	testRepository CartRepository
}

func TestMovieRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MovieRepositoryTestSuite))
}

func (suite *MovieRepositoryTestSuite) SetupTest() {
	testDb, err := sqlx.Open("godror", fmt.Sprintf("host=localhost port=5432 user=postgres password=Sanjit dbname=movie_db sslmode=disable"))
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}
	var sqlDB *sql.DB
	sqlDB, suite.mockDB, err = sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		panic(fmt.Sprintf("failed to open mock sql connection: %v", err))
	}
	suite.mockedDB = sqlx.NewDb(sqlDB, "godror")
	suite.db = testDb

	suite.testRepository = NewCartRepository(suite.db)
}

func (suite *MovieRepositoryTestSuite) TearDownTest() {
	suite.db.Close()
}
