## COSMOS ODYSSEY 
## How to run
### Install:
1. PostgreSQL
    - Download and install PostgreSql from [here](https://www.postgresql.org/download/)
2. Go
    - Download and install GO from [here](https://go.dev/doc/install)
3. Visual Studio Code
    - Download and install Visual Studio Code from [here](https://code.visualstudio.com/download)
    - Install Go extension by Go Team at Google
4. Go migrate
    - Follow the installation guide from [here](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
### Migrate database
- In the project folder open the file **api/database.go**
- At the start of the file, change variables **port**, **user**, **dbname** and **password** to your PostgreSql settings
- In the following command, replace all variables between <> with corresponding variables from the **database.go** file. Then navigate to project folder and run the migration with the command.

    `migrate -path db/migration -database "postgresql://<user>:<password>@localhost:<port>/<database_name>?sslmode=disable" -verbose up` 

### Run the code
- In Visual Studio Code, navigate to **cmd** folder of the project and run the code by command   `go run main.go`
- Open http://localhost:3000/ page on your browser