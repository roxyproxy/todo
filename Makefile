INTEGRATION_TEST_PATH?=./test/integration

# start docker components set in docker-compose.yaml
docker.start:
		docker-compose up --build -d

# shutting down docker components
docker.stop:
		docker-compose down;

# trigger integration test
test.integration:
		go test -tags=integration $(INTEGRATION_TEST_PATH) -count=1

# trigger integration test with verbose mode
test.integration.debug:
		go test -tags=integration $(INTEGRATION_TEST_PATH) -count=1 -v

