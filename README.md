## Flex Pay Backend

Project done for faculty purposes. Backend repo is related to frontend repo https://github.com/DenisKantic/flex_pay_frontend_angular

## Technologies used

### Backend
- Golang version 1.22
- Gin framework for Golang
- .env variables
- SMTP email settings (I used from my work email, **hetzner** hosting) 

### Database

- PostgreSQL (with procedures and functions)
 
### Docker

- Docker (used with dockerfile to spin up postgres database via command **docker compose up**)


## STEPS FOR STARTING THE PROJECT

- Clone this project 
- Copy the database creation table, procedures etc.. from **database_notes.txt** file into the postgres
- Start the Golang backend from your IDE (I'm using GoLand) or you can modify "docker-compose.yaml" to also include
Golang to start via Docker 
- use docker command to start already configured postgres inside **docker-compose.yaml** file with **docker compose up --build**
- After successfully starting Golang and Postgres, you can use Postman to test it out. **Send Login and Register data as objects inside postman**