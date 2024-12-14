# xm-company-crud


---

## 🔗 [SWAGGER UI](http://localhost:8000/swagger/)

***Swagger requires the builder-back-monolith service to run.**

---

# Makefile Commands:

Prerequisites: `Docker, make`

| Command                 | Description                                                                                                                             |
|-------------------------|-----------------------------------------------------------------------------------------------------------------------------------------|
| `app.start`             | Builds and run the project's application and services.                                                                                  |
| `app.start.clean`       | Migrate new DB schema and builds and run the project's application and services. (🤚This action wipes all the current data from the db) |
| `app.stop`              | Stops running project's application and services.                                                                                       |
| `swagger.init`          | Inits the swagger docs. This action is required when there are APIs changes/updates in order to maintain the SWAGGER UI up to date.     |
| `linter.run`            | Run linter checks located in `.golangci.yml` file.                                                                                      |
| `sql.migrate`           | Delete previous DB data and apply new schema. (🤚This action wipes all the current data from the db)                                    |
| `sql.migrate.populated` | Delete previous DB data and apply new schema and populates the DB with random data.                                                     |
| `sql.init`              | Restore the db to the initial schema and set triggers.                                                                                  |
| `sql.populate`          | Populate the `companies` table with random data for testing.                                                                            |
| `sql.drop`              | Drop all the tables. All tables (and their data) will be wiped and `sql.init` is required to have a functional db again.                |

---


---
