# AttendEase (WIP)
Streamline event management with AttendEase. Effortlessly generate and send unique QR-coded tickets to attendees. Implement a QR code scanner for seamless check-ins. 

### Features in progress
- [x] **1. Create an organisation account and allow login**
  - [x] Design and implement organization account creation functionality.
  - [x] Implement login functionality for organization accounts.

- [x] **2. Create event and upload CSV file with name and email of participants**
  - [x] Develop event creation functionality.
  - [x] Implement CSV file upload feature for participant details.

- [ ] **3. Send invitation emails to all participants with QR ticket and event info**
  - [ ] Set up email integration for sending invitations.
  - [ ] Generate QR tickets for participants.
  - [ ] Include event information in invitation emails.

- [x] **4. Create employee account inside the organisation for scanning QR**
  - [x] Implement employee account creation within the organization.
  - [x] Set up login functionality for employees.

- [ ] **5. Scan QR to mark attendance**
  - [ ] Develop QR code scanning functionality.
  - [x] Implement attendance marking logic.

- [ ] **6. Get list of attended and missed in form of Excel**
  - [ ] Create functionality to generate attendance reports.
  - [ ] Export reports to Excel format.

  

## Dependencies

- [godotenv](https://github.com/joho/godotenv): Used for loading environment variables from a `.env` file.
- [mongo-driver](https://go.mongodb.org/mongo-driver/mongo): MongoDB driver for Golang.
- [jwt](https://github.com/golang-jwt/jwt/v5): Golang implementation of JSON Web Tokens (JWT).
- [bcrypt](https://golang.org/x/crypto/bcrypt): A library for hashing and comparing passwords using bcrypt algorithm.
- [excelize](github.com/xuri/excelize) : A library for manipulating excel sheets

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