# Marketplace Api
REST API that handles entities and operations in a marketplace website, made using Go, Fiber(Go Framework), Gorm(Go ORM), MySQL/MariaDB, and JWT.
This project has clean architecture folder structure, which is based on [this github repository](https://github.com/bxcodec/go-clean-arch). The database diagram for this API can be seen at [this link](https://drive.google.com/file/d/1GPkRrlSdIww3BnxKaPDkcjue4Q5G79v-/view?usp=sharing). [Postman](https://www.postman.com/) can be used to access the API, the API Endpoints can be seen on [this Postman collection link](https://www.postman.com/vickonovianto/workspace/public-workspace/collection/457088-a5eccf56-e002-4483-b5fc-b29169cc9208?action=share&creator=457088). Before accessing API using Postman, we must change collection variable `local` into appropriate URL along with the `URL_PREFIX` we fill in step 5 below, for example the default value of variable `local` is `localhost:1213/api/v1`. After that, create two global variable with type secret in Postman, which are `userToken` and `adminToken`. `adminToken` is needed to access endpoints at `Category` except `Get All Category`. This API uses `Authorization: Bearer Token`.

## How to run the code
1. Create a new database for this API. No need to manually create other tables in the new database because the tables will be created automatically after executing `go run .`(Step 5).
2. Copy and rename file `example.env` into `.env`.
3. Open file `.env` and change `PORT`, `SECRET`, `DATABASE_URL`, and `API_PREFIX` into the appropriate port, secret for generating JWT token, database url, and api prefix.
4. Open terminal, go into root directory of this code, and run `go mod tidy`.
5. Then run `go run .`.
6. After creating a new user via the API, change manually the column `is_admin` with value `1` in table `user` to change a user into an admin in the database, because only an admin can create, get by ID, update, and delete categories. 
6. Press `Ctrl + C` to terminate the API.