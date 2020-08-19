#Creating the database structure

CREATE TABLE products
(
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id)
)

# Get the required files

go get -u github.com/gorilla/mux 
go get -u github.com/lib/pq

# Create the app files
┌── app.go
├── main.go
├── main_test.go
├── model.go
├── go.sum
└── go.mod