
start-podman:
	podman machine start

run-kafka:
	podman-compose up -d

down-kafka:
	podman-compose down

run-consumer:
	go run consumer/main.go

run-producer:
	go run producer/main.go