DB_DOCKER_CONTAINER=product_namer_db
BINARY_NAME=product_namer

postgres:
	docker run --name string ${DB_DOCKER_CONTAINER} -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

run:
	go run E:\golang_projects\namerGPT\cmd\main\main.go