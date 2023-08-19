run:
	export JWT_SECRET=someSecreteee21
	nodemon -e 'go, gohtml' -x 'go run main.go' --signal SIGTERM

seed:
	go run scripts/seed.go

test:
	go test -v ./...