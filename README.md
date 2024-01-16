# AttendEase (WIP)
Streamline event management with AttendEase. Effortlessly generate and send unique QR-coded tickets to attendees. Implement a QR code scanner for seamless check-ins. 

### Features
1. Create an organisation account and allow login
2. Create event and upload csv file with name and email of participants
3. Send invitaton emails to all participants with QR ticket and event info
4. Create employee account inside the organisation for scanning QR
5. Scan QR to mark attendence 
6. Get list of attended and missed in form of excel

## Dependencies

- [godotenv](https://github.com/joho/godotenv): Used for loading environment variables from a `.env` file.
- [mongo-driver](https://go.mongodb.org/mongo-driver/mongo): MongoDB driver for Golang.
- [jwt](https://github.com/golang-jwt/jwt/v5): Golang implementation of JSON Web Tokens (JWT).
- [bcrypt](https://golang.org/x/crypto/bcrypt): A library for hashing and comparing passwords using bcrypt algorithm.

## Installation

1. Clone the repository:

```bash
git clone https://github.com/aswinbennyofficial/attendease.git
```


2. Install dependencies:

```go
go mod tidy
```


3. Configure your environment variables by renaming `.env.example` into `.env`


## Usage
### Running the application

```bash
go run ./cmd/main/
```
By default, the server will start on port 8080.