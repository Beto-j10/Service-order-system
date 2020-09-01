# Service order system

Create, assign, edit and track service orders.

## development environment

- Launch PostgreSQL locally: 

`Docker run --name postgres --rm -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:12.4-alpine`

- Run core:

 `go run main.go`

### initial data

#### Users

Name: "Anderson", Email: "anderson@dominio.com", Password: "012"

Name: "Linda", Email: "linda@dominio.com", Password: "234"

Name: "Lucia", Email: "lucia@dominio.com", Password: "456"

Name: "Georgi", Email: "georgi@dominio.com", Password: "678"

#### Technicians

Name: "David", Email: "david@harper.com", Password: "210", Code: 3323

Name: "Norman", Email: "norman@harper.com", Password: "432", Code: 8112

Name: "Lizeth", Email: "lizeth@harper.com", Password: "654", Code: 4211

Name: "Azul", Email: "azul@harper.com", Password: "876", Code: 1017

Name: "Lorena", Email: "lorena@harper.com", Password: "098", Code: 2201


#### Notes

- Database has no persistence.
