run: docker
	@sleep 1
	@go run .

docker:
	@docker compose up -d --wait

clean:
	@docker compose down -v
