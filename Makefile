run:
	go run cmd/main.go

test:
	go test ./... -count=1

migrate-liquibase:
	liquibase --changeLogFile=db/changelog.xml --url="jdbc:postgresql://localhost:5432/movie_db" --username=postgres --password=Sanjit update
