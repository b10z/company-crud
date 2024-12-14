#INCLUDES------------------------------------------------------------------------------------
include app.env

#APPLICATION OPTIONS-------------------------------------------------------------------------
app.start:
	make app.stop
	docker compose -f ./deployments/docker-compose.yml build
	docker compose -f ./deployments/docker-compose.yml up -d

app.start.clean:
	make sql.migrate
	make app.start

app.stop:
	docker compose -f ./deployments/docker-compose.yml stop

#TESTS OPTIONS-------------------------------------------------------------------------------
tests.test-clear:
	@docker image rm -f sevice-test-build || (echo "Image sevice-test-build didn't exist."; exit 0)

tests.test-build:
	make tests.test-clear
	@docker build \
		--tag sevice-test-build \
		-f ./deployments/test/Dockerfile .


#SQL DB OPTIONS------------------------------------------------------------------------------
sql.migrate:
	make tests.test-build
	@docker run --network="host" --volume .:/app --workdir /app \
    	sevice-test-build /bin/bash -c "PGPASSWORD=${DB_PASSWORD} psql -h localhost -U ${DB_USERNAME} -d ${DB_NAME} -a -f ./scripts/drop.sql \
    									PGPASSWORD=${DB_PASSWORD} psql -h localhost -U ${DB_USERNAME} -d ${DB_NAME} -a -f ./scripts/init.sql";\
	make tests.test-clear

sql.migrate.populated:
	make tests.test-build
	@docker run --network="host" --volume .:/app --workdir /app \
    	sevice-test-build /bin/bash -c "PGPASSWORD=${DB_PASSWORD} psql -h localhost -U ${DB_USERNAME} -d ${DB_NAME} -a -f ./scripts/drop.sql \
    									PGPASSWORD=${DB_PASSWORD} psql -h localhost -U ${DB_USERNAME} -d ${DB_NAME} -a -f ./scripts/init.sql \
    									PGPASSWORD=${DB_PASSWORD} psql -h localhost -U ${DB_USERNAME} -d ${DB_NAME} -a -f ./scripts/populate.sql";\
	make tests.test-clear

sql.init:
	make tests.test-build
	@docker run --network="host" --volume .:/app --workdir /app \
    	sevice-test-build /bin/bash -c "PGPASSWORD=${DB_PASSWORD} psql -h localhost -U ${DB_USERNAME} -d ${DB_NAME} -a -f ./scripts/init.sql";\
	make tests.test-clear

sql.drop:
	make tests.test-build
	@docker run --network="host" --volume .:/app --workdir /app \
    	sevice-test-build /bin/bash -c "PGPASSWORD=${DB_PASSWORD} psql -h localhost -U ${DB_USERNAME} -d ${DB_NAME} -a -f ./scripts/drop.sql";\
    make tests.test-clear

sql.populate:
	make tests.test-build
	@docker run --network="host" --volume .:/app --workdir /app \
    	sevice-test-build /bin/bash -c "PGPASSWORD=${DB_PASSWORD} psql -h localhost -U ${DB_USERNAME} -d ${DB_NAME} -a -f ./scripts/populate.sql";\
    make tests.test-clear