# Go Refresh Token API

This project is a RESTful API built using Go that implements user authentication with JWT (JSON Web Tokens) and refresh tokens. It allows users to register, log in, and manage items with standard CRUD operations.

## Features

- User Registration
- User Login with JWT and Refresh Tokens
- CRUD Operations for Items
- Structured JSON Responses
- Input Validation

## Technologies Used

- Go
- GORM (for database interactions)
- Gorilla Mux (for routing)
- JWT (for token generation and validation)
- PostgreSQL (or any other database of your choice)

## Getting Started

### Prerequisites

Make sure you have the following installed:

- Go (version 1.16 or higher)
- PostgreSQL (or any other database you plan to use)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/go-refresh-token.git
   cd go-refresh-token
   CopyInsert
   Install the required dependencies:
   ```

go mod tidy
CopyInsert
(Optional) Update the database connection settings in database.go to match your environment.

Running the Application
Start the server:

go run main.go
CopyInsert
The server will start on localhost:8080 (by default). You can change the port in the main.go file if needed.

API Endpoints
User Registration
URL: POST /api/register
Body:
{
"username": "your_username",
"password": "your_password"
}
CopyInsert
User Login
URL: POST /api/login

Body:

{
"username": "your_username",
"password": "your_password"
}
CopyInsert
Response:

{
"status": "success",
"message": "Login successful",
"token": "your_access_token",
"refreshToken": "your_refresh_token"
}
CopyInsert
Get All Items
URL: GET /api/items
Response:
{
"status": true,
"message": "Items retrieved successfully",
"items": [...]
}
CopyInsert
Get Single Item
URL: GET /api/items/{id}
Response:
{
"status": true,
"message": "Item retrieved successfully",
"item": {...}
}
CopyInsert
Create Item
URL: POST /api/items

Body:

{
"name": "item_name",
"price": 100
}
CopyInsert
Response:

{
"status": true,
"message": "Item created successfully",
"item": {...}
}
CopyInsert
Update Item
URL: PUT /api/items/{id}

Body:

{
"name": "updated_item_name",
"price": 150
}
CopyInsert
Response:

{
"status": true,
"message": "Item updated successfully",
"item": {...}
}
CopyInsert
Delete Item
URL: DELETE /api/items/{id}
Response:
{
"status": false,
"message": "Item Not Found"
}
CopyInsert
Project Structure
go-refresh-token/
├── controllers/ # Contains the API controller logic
├── database/ # Database connection and models
├── utils/ # Helper functions (e.g., for password hashing and token generation)
├── main.go # Entry point for the application
├── go.mod # Go module file
└── README.md # Project documentation
CopyInsert
Contributing
Contributions are welcome! Please create an issue or submit a pull request if you have suggestions or improvements.

License
This project is licensed under the MIT License - see the LICENSE file for details.
