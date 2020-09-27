run:
	docker-compose up

stop:
	docker-compose stop

tests:
	go test

tests-cover:
	go test -cover