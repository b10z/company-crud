# xm-company-crud


---

## 🔗 [SWAGGER UI](http://localhost:8000/swagger/)
***JWT Token without expiration date**: **`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.PwN9mqs6JDOROs42oqojiJ0iGEzOtLejuVrDPITuxqw`** 

***Swagger requires the Company CRUD service to run.**

---

# Instructions:
Prerequisites: `Docker, make`

- Clone the repo. 
- `cd` into the project
- Run `make app.start`. 

These steps will launch the **service**, an initialized **DB**, and a **Kafka broker** with initialized topics.


# Makefile Commands:

| Command                 | Description                                                                                                                             |
|-------------------------|-----------------------------------------------------------------------------------------------------------------------------------------|
| `app.start`             | Builds and run the project's application and services.                                                                                  |
| `app.start.clean`       | Migrate new DB schema and builds and run the project's application and services. (🤚This action wipes all the current data from the db) |
| `app.stop`              | Stops running project's application and services.                                                                                       |
| `swagger.init`          | Inits the swagger docs. This action is required when there are APIs changes/updates in order to maintain the SWAGGER UI up to date.     |
| `sql.migrate`           | Delete previous DB data and apply new schema. (🤚This action wipes all the current data from the db)                                    |
| `sql.migrate.populated` | Delete previous DB data and apply new schema and populates the DB with random data.                                                     |
| `sql.init`              | Restore the db to the initial schema and set triggers.                                                                                  |
| `sql.populate`          | Populate the `companies` table with random data for testing.                                                                            |
| `sql.drop`              | Drop all the tables. All tables (and their data) will be wiped and `sql.init` is required to have a functional db again.                |

---


---

# Short mention of the REST-endpoints: 

- POST - `/companies`
- GET - `/companies/{company_name}` 
- DELETE - `/companies/{company_name}`
- PATCH - `/companies/{company_name}`

For more details, please refer to SWAGGER.

# Information about testing:

An Integration-like test was implemented `cmd/company_crud/main_test.go`.
This test, utilizes a custom suite that uses a temporal DB and a Kafka broker. Then multiple cases run using actual client calls to the service, also the data are verified from the service and directly from the DB.

