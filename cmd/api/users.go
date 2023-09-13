package main

import (
	"errors"
	"fmt"
	"go-api/internal/models"
	"go-api/internal/serializer"
	"go-api/internal/validator"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the request body into the anonymous struct.
	err := serializer.DeserializeFromJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &models.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if validator.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert the user data into the database.
	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// After the user record has been created in the database, generate a new activation token for the user.
	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, models.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Launch a goroutine which runs an anonymous function that sends the welcome email.
	app.background(func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		err = app.mailer.Send(user.Email, "user_welcome.gohtml", data)
		if err != nil {
			app.logger.Error(err, nil)
		}
	})

	err = serializer.SerializeToJson(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := serializer.DeserializeFromJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the plaintext token provided by the client.
	v := validator.New()
	if validator.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Retrieve the details of the user associated with the token using the
	// GetForToken() method (which we will create in a minute). If no matching record
	// is found, then we let the client know that the token they provided is not valid.
	user, err := app.models.Users.GetForToken(models.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Update the user's activation status.
	user.Activated = true

	// Save the updated user record in our database, checking for any edit conflicts in
	// the same way that we did for our movie records.
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// If everything went successfully, then we delete all activation tokens for the user.
	err = app.models.Tokens.DeleteAllForUser(models.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Add the "movies:read" permission for the user now that he is activated.
	err = app.models.Permissions.AddForUser(user.ID, "movies:read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send the updated user details to the client in a JSON response.
	err = serializer.SerializeToJson(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Password       string `json:"password"`
		TokenPlaintext string `json:"token"`
	}

	err := serializer.DeserializeFromJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	validator.ValidatePasswordPlaintext(v, input.Password)
	validator.ValidateTokenPlaintext(v, input.TokenPlaintext)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Retrieve the details of the user associated with the password reset token,
	// returning an error message if no matching record was found.
	user, err := app.models.Users.GetForToken(models.ScopePasswordReset, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			v.AddError("token", "invalid or expired password reset token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Set the new password for the user.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Save the updated user record in our database, checking for any edit conflicts
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// If everything was successful, then delete all password reset tokens for the user.
	err = app.models.Tokens.DeleteAllForUser(models.ScopePasswordReset, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send the user a confirmation message.
	env := envelope{"message": "your password was successfully reset"}
	err = serializer.SerializeToJson(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) uploadUserDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	const MaxUploadSize = 1024 * 1024 * 32 // 32MB
	const MaxUploadFileSize = 1024 * 1024  // 1MB

	var input struct {
		Files []*multipart.FileHeader `json:"files"`
	}

	// check if user authenticated
	//user := app.contextGetUser(r)
	//if user.IsAnonymous() {
	//	app.authenticationRequiredResponse(w, r)
	//	return
	//}

	// 32 MB is the default used by FormFile()
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("the uploaded files are too big. Please use files less than %dMB in size", MaxUploadSize))
		return
	}

	// Get a reference to the fileHeaders.
	// They are accessible only after ParseMultipartForm is called
	input.Files = r.MultipartForm.File["file"]

	for _, fileHeader := range input.Files {
		// Restrict the size of each uploaded file to 1MB.
		// To prevent the aggregate size from exceeding
		// a specified value, use the http.MaxBytesReader() method
		// before calling ParseMultipartForm()
		if fileHeader.Size > MaxUploadFileSize {
			app.badRequestResponse(w, r, fmt.Errorf("the uploaded file is too big. Please use file less than %dMB in size", MaxUploadFileSize))
			return
		}

		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				app.logger.Fatal(err, nil)
			}
		}(file)

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" {
			app.badRequestResponse(w, r, fmt.Errorf("the provided file format is not allowed. Please upload a JPEG or PNG image"))
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		err = os.MkdirAll("./storage", os.ModePerm)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		f, err := os.Create(fmt.Sprintf("./storage/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				app.logger.Fatal(err, nil)
			}
		}(f)

		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Send the user a confirmation message.
	env := envelope{"message": "Documents uploaded successfully"}
	err := serializer.SerializeToJson(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
