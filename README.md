# AttendEase 
Streamline event management with AttendEase. Effortlessly generate and send unique QR-coded tickets to attendees. Implement a QR code scanner for seamless check-ins. 

### Features in progress
- [x] **1. Create an organisation account and allow login**
  - [x] Design and implement organization account creation functionality.
  - [x] Implement login functionality for organization accounts.

- [x] **2. Create event and upload CSV file with name and email of participants**
  - [x] Develop event creation functionality.
  - [x] Implement CSV file upload feature for participant details.

- [x] **3. Send invitation emails to all participants with QR ticket and event info**
  - [x] Set up email integration for sending invitations.
  - [x] Generate QR tickets for participants.
  - [x] Include event information in invitation emails.

- [x] **4. Create employee account inside the organisation for scanning QR**
  - [x] Implement employee account creation within the organization.
  - [x] Set up login functionality for employees.

- [x] **5. Scan QR to mark attendance**
  - [x] Implement attendance marking logic.

- [x] **6. Get list of attended and missed in form of Excel**
  - [x] Create functionality to generate attendance reports.
  - [x] Export reports to Excel format.

  
## Routes

Health Check:
  - Endpoint: `/health`
  - Handler: controllers.HandleHealth


Private Endpoint (Requires Admin Login):
  - Endpoint: `/private`
  - Handler: controllers.HandlePrivate
  - Middleware: middleware.AdminLoginRequired

Admin Authentication and Authorization:
  - Sign in: `/api/admin/login`
  - Sign up: `/api/admin/signup`
  - Refresh token:`/api/admin/refresh`
  - Logout: `/api/logout`

Event Management:
  - Create event: `/api/events` (POST)
  - Get all events: `/api/events` (GET)
  - Get a specific event by ID: `/api/events/{eventid}` (GET)
  - Upload participants for an event: `/api/events/{eventid}/participants` (POST)
  - Middleware for these routes: middleware.AdminLoginRequired

Employee Management:
  - Create employee: `/api/employees` (POST)
  - Get all employees: `/api/employees` (GET)
  - Employee login: `/api/employees/login` (POST)
  - Middleware for these routes: middleware.AdminLoginRequired

Scanning:
  - Scan a participant: `/api/events/scan` (POST)
  - Middleware: middleware.LoginRequired

Participant Management in an Event:
  - Get all participants of an event: `/api/events/{eventid}/participants` (GET)
  - Get all participants of an event in a file: `/api/events/{eventid}/participants/file` (GET)
  - Send email to all participants of an event: `/api/events/{eventid}/send` (GET)
  - Middleware: middleware.AdminLoginRequired


[<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://god.gw.postman.com/run-collection/30160615-b18996b7-589a-4589-8772-db4820a41cd3?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D30160615-b18996b7-589a-4589-8772-db4820a41cd3%26entityType%3Dcollection%26workspaceId%3D854bfee7-441b-4733-a866-80484917b281)

## Dependencies

- [godotenv](https://github.com/joho/godotenv): Used for loading environment variables from a `.env` file.
- [mongo-driver](https://go.mongodb.org/mongo-driver/mongo): MongoDB driver for Golang.
- [jwt](https://github.com/golang-jwt/jwt/v5): Golang implementation of JSON Web Tokens (JWT).
- [bcrypt](https://golang.org/x/crypto/bcrypt): A library for hashing and comparing passwords using bcrypt algorithm.
- [excelize](github.com/xuri/excelize) : A library for manipulating excel sheets
- [amqp091-go](github.com/rabbitmq/amqp091-go) : A library to implement queue using rabbitMQ
- [go-qrcode](github.com/skip2/go-qrcode) : A library to generate QR code

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
### Running the main application

```bash
go run ./cmd/main/
```

### Start RabbitMQ 
- I am using Docker to run it
```bash
docker run --rm -it -p 15672:15672 -p 5672:5672 rabbitmq:3-management
```

### Running email sender with workers
```bash
go run ./cmd/email_queue_consumer/
```
By default, the server will start on port 8080.