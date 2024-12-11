# Book Store System

Welcome to the Book Store System project! This system allows users to log in/out,insert book in book listing, review books, and perform other related actions.

## Features
- User authentication: Users can register, log in, and log out securely.
- Book review: Users can read and write reviews for books.
- Book rating: Users can rate books.
- Book search: Users can search for books by isbn.
- User profile: Users can view and update their few profile information.

## Technologies Used
- GoLang: BackEnd development
- MySQL: Database management
- Redis: For Fast Retrieval and put off some load from database
- ByCrypt: Password Management
- Regex: Few places like email
- httprouter: For routing servmux was good (but i need pattern like `localhost:8000/user/activation/token/` ,`localhost:8000/new_password/reset-token/` token changed everytime even i can do it in home route but it make it look shitty so at the end of the day used this man -> httprouter)
- Secureheader:  Added Few headers For security Related
- TLS:  Added TLS Certificate Encryption 
<!-- - JWT (JSON Web Tokens): Authentication mechanism -->

## Getting Started

### Prerequisites
- GoLang installed on your system. You can download it from [here](https://golang.org/dl/).
- MySQL installed on your system. You can download it from [here](https://www.mysql.com/download/).
- Redis installed on your system. You can read & download it from [here](https://redis.io/docs/latest/operate/oss_and_stack/install/install-redis/install-redis-on-linux/).
- Before sending data set Header Content-Type: "application/json" in Postman / cURL or anyother for POST & other Non-Get methods.

### Installation
1. Clone the repository:
```bash
    git clone github.com/iamgak/go-bookstore
```

2. Navigate to the project directory:
```bash
    cd go-bookstore
```

3. Set up the database:
    - Create a MySQL database named `go-bookstore`.
    - Import the database schema from `database/schema.sql`.


4. Install dependencies:
```bash
    go mod tidy
```
5. Create TLS key & Ceritificate:

Search, where Go file is saved then goto src then crypto then tls then copy whole path something like below

```bash
    mkdir tls; cd tls;
    go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
```

After that if you have .gitignore than add manually or 

```bash
    cd go-bookstore;
    echo 'tls/' << .gitignore
```

6. Run the server:

```bash
    go run cmd/cli
```

    optional if you want to change port number and root info
```bash
    go run cmd/cli -addr=":8000" -dsn="root:@/bookstore?parseTime=true"
 ```

7. Access the application:
    Open your web browser and navigate to `https://localhost:8000`.

### Usage
- Register a new user account By PostMethod `https://localhost:8000/user/register`.
- Validate account By GetMethod `https://localhost:8000/user/activation/activation-token` activation-token is in user table after registration in activation_token attribute.
- Login/Logout in with existing credentials By PostMethod `https://localhost:8000/user/login` / `https://localhost:8000/user/logout`.
- Browse all reviews By GetMethod `https://localhost:8000/review/listing`.
- Delete reviews By GetMethod `https://localhost:8000/book/delete/id`.
- Browse reviews By ISBN or Author or Book GetMethod `https://localhost:8000/review/search/934-3434` or `https://localhost:8000/review/search/murakami` .
- Insert new book info By PostMethod After Login `https://localhost:8000/book/create` .
- Write reviews for books you've read By PostMethod After Login `https://localhost:8000/review/create` .
- Browse All Book By GetMethod `https://localhost:8000/book/listing`
- Request Forget Password By PostMethod  `https://localhost:8000/user/forget_password/`token send on given email if registered but in this case you will copy it from table forget_passw during new password.
- Change Forget Password By PostMethod  `https://localhost:8000/user/new_password/reset-token` token from forget_password .

## Infuture
I want to split it into multiple microservices and put it in a single docker file and add few more context and channels.

## Contributing
Contributions are welcome! If you'd like to contribute to this project, please fork the repository and submit a pull request with your changes.

