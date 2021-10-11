INTEGRATION_TEST_PATH?=./test/integration

# set of env variables that you need for testing
ENV_LOCAL_TEST=\
		POSTGRES_PASSWORD=mysecretpassword \
		POSTGRES_DB=postgres \
		POSTGRES_HOST=postgres \
		POSTGRES_USER=postgres

# start docker components set in docker-compose.yaml
docker.start:
		docker-compose up -d --remove-orphans postgres;

# shutting down docker components
docker.stop:
		docker-compose down;

# trigger integration test
test.integration:
		$(ENV_LOCAL_TEST) \
		go test -tags=integration $(INTEGRATION_TEST_PATH) -count=1

# trigger integration test with verbose mode
test.integration.debug:
		$(ENV_LOCAL_TEST) \
		go test -tags=integration $(INTEGRATION_TEST_PATH) -count=1 -v

