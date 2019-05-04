all: start

start: bins
		docker-compose --file ./docker/docker-compose.yaml up

start-detached: bins
		docker-compose --file ./docker/docker-compose.yaml up -d

bins:
		dep ensure -v
		go build

restart:
		docker-compose --file ./docker/docker-compose.yaml restart

stop:
		docker-compose --file ./docker/docker-compose.yaml stop

clean: stop
		docker-compose --file ./docker/docker-compose.yaml down --rmi 'local'
		go clean -x
		rm -r --force vendor/

logs:
		docker-compose --file ./docker/docker-compose.yaml logs

debug: bins
		docker-compose --file ./docker/docker-compose.debug.yaml up
