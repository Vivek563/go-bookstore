package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"test.iamgak.net/models"
	"test.iamgak.net/validator"
	"time"
)

type Message struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// home page
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	resp := Message{
		Status:  true,
		Message: "Welcome to our Bookstore, Website",
	}

	app.sendJSONResponse(w, 200, &resp)
}

// all the review listing
func (app *application) ReviewListing(w http.ResponseWriter, r *http.Request) {
	bks, err := app.models.Review.ReviewListing()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sendJSONResponse(w, 200, bks)
}

// logged user review
func (app *application) MyReview(w http.ResponseWriter, r *http.Request) {
	bks, err := app.models.Review.MyReview(app.user_id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sendJSONResponse(w, 200, bks)
}

// search by isbn only
func (app *application) ReviewSearch(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	isbn := params.ByName("isbn")
	bks, err := app.models.Review.GetReviewByIsbn(isbn)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sendJSONResponse(w, 200, bks)
}

func (app *application) DeleteReview(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	review_id, err := strconv.ParseInt(params.ByName("id"), 10, 32)
	if err != nil {
		app.notFound(w)
		return
	}

	err = app.models.Review.DeleteReview(review_id, app.user_id)
	if err != nil {
		app.notFound(w)
		return
	}

	app.models.Users.ActivityLog("review_deleted", app.user_id)
	resp := app.sendMessage(true, "Review Deleted")
	app.sendJSONResponse(w, 200, resp)
}

// create new review

func (app *application) AddReview(w http.ResponseWriter, r *http.Request) {
	var CreateReview *models.Review
	err := json.NewDecoder(r.Body).Decode(&CreateReview)
	if err != nil {
		app.errorLog.Print(err)
		app.CustomError(w, "header should be application/json and all the field should be in proper format", 400)
		return
	}

	CreateReview.Uid = app.user_id
	validator := &validator.Validator{
		Errors: make(map[string]string),
	}

	validator.CheckField(validator.NotBlank(CreateReview.Isbn), "isbn", "Please, fill the isbn field")
	validator.CheckField(validator.NotBlank(CreateReview.Descriptions), "descriptions", "Please, fill the descriptions field")
	validator.CheckField(validator.NotBlank(CreateReview.Title), "title", "Please, fill the title field")

	if validator.Errors["isbn"] == "" {
		validator.CheckField(validator.MaxChars(CreateReview.Title, 20), "isbn", "Please, fill the ISBN shorter than 20")
	}

	if validator.Errors["descriptions"] == "" {
		validator.CheckField(validator.MaxChars(CreateReview.Descriptions, 100), "isbn", "Please, fill the ISBN shorter than 100")
	}

	if validator.Errors["title"] == "" {
		validator.CheckField(validator.MaxChars(CreateReview.Title, 50), "title", "Please, fill the TITLE shorter than 50")
	}

	if validator.Valid() {
		book_id := app.models.Books.BookExist(CreateReview.Isbn)
		if book_id {
			validator.Errors["isbn"] = "isbn already Exist"
		}
	}

	if !validator.Valid() {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	err = app.models.Review.CreateReview(CreateReview)
	if err != nil {
		app.CustomError(w, "Server Issue", 500)
		return
	}

	app.models.Users.ActivityLog("review_created", app.user_id)
	resp := app.sendMessage(true, "Review Saved")
	app.sendJSONResponse(w, 200, resp)
}

// bookListing related handlers
func (app *application) BookListing(w http.ResponseWriter, r *http.Request) {

	bks, err := app.models.Books.BooksListing()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sendJSONResponse(w, 200, bks)
}

// book info based on isbn
func (app *application) BookInfo(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	isbn := params.ByName("isbn")
	validator := &validator.Validator{
		Errors: make(map[string]string),
	}
	validator.CheckField(validator.NotBlank(isbn), "ISBN", "Empty, ISBN field")

	if !validator.Valid() {
		app.sendJSONResponse(w, 200, validator.Errors)
		return
	}

	info, err := app.models.Books.GetBookByIsbn(isbn)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sendJSONResponse(w, 200, info)
}

// add book in db
func (app *application) AddBook(w http.ResponseWriter, r *http.Request) {
	var bookRegister *models.Book
	err := json.NewDecoder(r.Body).Decode(&bookRegister)
	if err != nil {
		app.errorLog.Print(err)
		app.CustomError(w, "Unsupported or empty fields", 400)
		return
	}

	validator := &validator.Validator{
		Errors: make(map[string]string),
	}

	validator.CheckField(validator.NotBlank(bookRegister.ISBN), "isbn", "Please, fill the isbn field")
	validator.CheckField(validator.NotBlank(bookRegister.Genre), "genre", "Please, fill the genre field")
	validator.CheckField(validator.NotBlank(bookRegister.Descriptions), "descriptions", "Please, fill the descriptions field")
	validator.CheckField(validator.NotBlank(bookRegister.Author), "author", "Please, fill the author field")
	validator.CheckField(validator.NotBlank(bookRegister.Title), "title", "Please, fill the title field")

	if validator.Errors["isbn"] == "" {
		validator.CheckField(validator.MaxChars(bookRegister.Title, 20), "isbn", "Please, fill the ISBN shorter than 20")
	}

	if validator.Errors["descriptions"] == "" {
		validator.CheckField(validator.MaxChars(bookRegister.Descriptions, 100), "isbn", "Please, fill the ISBN shorter than 100")
	}

	if validator.Errors["author"] == "" {
		validator.CheckField(validator.MaxChars(bookRegister.Author, 50), "author", "Please, fill the AUTHOR shorter than 50")
	}
	if validator.Errors["title"] == "" {
		validator.CheckField(validator.MaxChars(bookRegister.Title, 50), "title", "Please, fill the TITLE shorter than 50")
	}

	if validator.Valid() {
		book_id := app.models.Books.BookExist(bookRegister.ISBN)
		if book_id {
			validator.Errors["isbn"] = "isbn already Exist"
		}
	}

	if !validator.Valid() {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	err = app.models.Books.CreateBook(bookRegister)
	if err != nil {
		app.CustomError(w, "Server Issue1", 500)
		return
	}

	app.models.Users.ActivityLog("Book Listed", app.user_id)
	resp := app.sendMessage(true, "Book Record Saved, Sucessfully")
	app.sendJSONResponse(w, 200, resp)
}

// register user
func (app *application) UserRegister(w http.ResponseWriter, r *http.Request) {

	var creds *models.UserRegister
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.errorLog.Print(err)
		app.CustomError(w, "Unsupported or empty fields", 400)
		return
	}

	validator := &validator.Validator{
		Errors: make(map[string]string),
	}

	validator.CheckField(validator.NotBlank(creds.Email), "email", "Please, fill the email field")
	validator.CheckField(validator.NotBlank(creds.Password), "password", "Please, fill the password field")
	validator.CheckField(validator.NotBlank(creds.RepeatPassword), "repeatPassword", "Please, fill the repeat password field")

	validator.ValidPassword(creds.Password)

	if validator.Errors["email"] == "" {
		validator.CheckField(validator.ValidEmail(creds.Email), "email", "Invalid Email Format")
	}
	if app.models.Users.EmailExist(creds.Email) != 0 {
		validator.CheckField(false, "email", "Email already registered")
	}

	if validator.Errors["repeatPassword"] == "" && validator.Errors["password"] == "" {
		if creds.Password != creds.RepeatPassword {
			validator.Errors["repeatPassword"] = "Password not matched"
		}
	}

	if !validator.Valid() {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	uri := app.generateHash(r.RemoteAddr, r.URL.Port())
	uid, err := app.models.Users.InsertUser(creds.Email, creds.Password, uri)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.models.Users.ActivityLog("Account Created", uid)
	resp := Message{
		Status:  true,
		Message: "Registration Successfull ",
	}

	app.sendJSONResponse(w, 200, resp)
}

// after registration user need to validate using link save in users table respective email
func (app *application) UserActivation(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	uri := params.ByName("uri")
	if !app.models.Users.ValidURI(uri) {
		app.notFound(w)
		return
	}

	err := app.models.Users.AccountActivate(uri)
	if err != nil {
		app.notFound(w)
		return
	}

	resp := app.sendMessage(true, "Your account has been Verified.")
	app.sendJSONResponse(w, 200, resp)
}

// login
func (app *application) UserLogin(w http.ResponseWriter, r *http.Request) {
	var creds *models.UserLogin
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.errorLog.Print(err)
		app.CustomError(w, "Unsupported or empty fields", 400)
		return
	}

	validator := &validator.Validator{
		Errors: make(map[string]string),
	}

	validator.CheckField(validator.NotBlank(creds.Email), "email", "Please, fill the email field")
	validator.CheckField(validator.NotBlank(creds.Password), "password", "Please, fill the password field")

	if validator.Errors["email"] == "" {
		validator.CheckField(validator.ValidEmail(creds.Email), "email", "Invalid Email Format")
	}

	if !validator.Valid() {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	uid := app.models.Users.EmailExist(creds.Email)
	if uid == 0 {
		resp := app.sendMessage(false, "Incorrect Credentials")
		app.sendJSONResponse(w, 200, resp)
		return
	}

	uid, err = app.models.Users.Login(creds)
	if err != nil {
		app.CustomError(w, "Internal Server Error", 500)
		return
	}

	if uid == 0 {
		resp := app.sendMessage(false, "Incorrect Credentials")
		app.sendJSONResponse(w, 200, resp)
		return
	}

	hashed := app.generateHash(r.RemoteAddr, r.URL.Port())
	err = app.models.Users.SetLoginToken(hashed, uid)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.models.Users.ActivityLog("logged_in", app.user_id)
	cookie := &http.Cookie{
		Name:    "ldata",
		Value:   hashed,
		Expires: time.Now().Add(1 * time.Hour),
		Path:    "/",
	}

	http.SetCookie(w, cookie)

	resp := app.sendMessage(true, "Login Successfull")
	app.sendJSONResponse(w, 200, resp)
}

func (app *application) UserLogout(w http.ResponseWriter, r *http.Request) {
	app.models.Users.ActivityLog("log_out", app.user_id)
	resp := app.sendMessage(true, "Logout Successfull")
	app.sendJSONResponse(w, 200, resp)
}

// after forget password it create uri in db like /new_password/db_uri
func (app *application) NewPasswordPost(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())
	uri := params.ByName("uri")

	uid, err := app.models.Users.ForgetPasswordUri(uri)
	if err != nil {
		app.notFound(w)
		return
	}

	var creds *models.UserNewPassword
	err = json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.errorLog.Print(err)
		app.CustomError(w, "Unsupported or empty fields", 400)
		return
	}

	validator := &validator.Validator{
		Errors: make(map[string]string),
	}

	validator.CheckField(validator.NotBlank(creds.Password), "password", "Please, fill the password field")

	validator.CheckField(validator.NotBlank(creds.RepeatPassword), "repeatPassword", "Please, fill the password field")
	if validator.Errors["repeatPassword"] == "" && validator.Errors["password"] == "" {
		if creds.Password != creds.RepeatPassword {
			validator.Errors["repeatPassword"] = "Password not matched"
		}
	}

	if !validator.Valid() {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	err = app.models.Users.NewPassword(creds.Password, uid)
	if err != nil {
		app.CustomError(w, "Internal Server Error"+err.Error(), 500)
		return
	}

	resp := app.sendMessage(true, "Password Changed Successfully")
	app.sendJSONResponse(w, 200, resp)
}

// it will create a uri if email is valid and active user
func (app *application) ForgetPasswordPost(w http.ResponseWriter, r *http.Request) {

	var creds *models.ForgetPassword
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.errorLog.Print(err)
		app.CustomError(w, "Empty field / Unsupported Content-Type", 400)
		return
	}

	validator := &validator.Validator{
		Errors: make(map[string]string),
	}

	validator.CheckField(validator.NotBlank(creds.Email), "email", "Please, fill the email field")

	if !validator.Valid() {
		app.sendJSONResponse(w, 200, validator)
		return
	}

	uid := app.models.Users.EmailExist(creds.Email)
	if uid > 0 {
		token := app.generateHash(creds.Email, r.RemoteAddr)
		uri := app.generateHash(token, r.RemoteAddr)
		err := app.models.Users.ForgetPassword(uid, uri)
		if err != nil {
			app.CustomError(w, "Internal Server Error", 400)
		}
	}

	// whether email is registered or not we will show same success message
	resp := app.sendMessage(true, "If email is registered email you will get link on your email")
	app.sendJSONResponse(w, 200, resp)
}

// create a token for login, user_verification
func (app *application) generateHash(addr, port string) string {
	data := addr + port + strconv.FormatInt(time.Now().Unix(), 10)
	hasher := sha1.New()
	hasher.Write([]byte(data))
	hash := hasher.Sum(nil)
	// Convert hash bytes to a hexadecimal string
	hashStr := fmt.Sprintf("%x", hash)

	return hashStr
}

// to print json message
func (app *application) sendJSONResponse(w http.ResponseWriter, statusCode int, message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}

// create a struct Message to send data in sendJSONResponse
func (app *application) sendMessage(status bool, message string) Message {
	return Message{
		Status:  status,
		Message: message,
	}
}
