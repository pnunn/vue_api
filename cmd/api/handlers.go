package main

import (
  "errors"
  "net/http"
  "time"
)

type jsonResponse struct {
  Error   bool        `json:"error"`
  Message string      `json:"message"`
  Data    interface{} `json:"data,omitempty"`
}

type envelope map[string]interface{}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
  type credentials struct {
    UserName string `json:"email"`
    Password string `json:"password"`
  }

  var creds credentials
  var payload jsonResponse

  err := app.readJSON(w, r, &creds)
  if err != nil {
    app.errorLog.Println(err)
    payload.Error = true
    payload.Message = "invalid json supplied or missing"
    _ = app.writeJSON(w, http.StatusBadRequest, payload)
  }

  // TODO authenticate
  app.infoLog.Println(creds.UserName, creds.Password)

  // look up the user by email
  user, err := app.models.User.GetByEmail(creds.UserName)
  if err != nil {
    app.errorJSON(w, errors.New("invalid username/passwond"))
    return
  }

  // validate the user's password
  validPassword, err := user.PasswordMatches(creds.Password)
  if err != nil || !validPassword {
    app.errorJSON(w, errors.New("invalid username/password"))
    return
  }

  // we have a valid user, generate token
  token, err := app.models.Token.GenerateToken(user.ID, 24*time.Hour)
  if err != nil {
    app.errorJSON(w, err)
    return
  }

  // save to the database
  err = app.models.Token.Insert(*token, *user)
  if err != nil {
    app.errorJSON(w, err)
    return
  }

  // send back a response
  payload = jsonResponse{
    Error:   false,
    Message: "logged in",
    Data:    envelope{"token": token},
  }

  //app.infoLog.Println(token)
  err = app.writeJSON(w, http.StatusOK, payload)
  if err != nil {
    app.errorLog.Println(err)
  }
}

func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
  var requestPayload struct {
    Token string `json:"token"`
  }

  err := app.readJSON(w, r, &requestPayload)
  if err != nil {
    app.errorLog.Println(err)
    app.errorJSON(w, errors.New("invalid json 1"))

    return
  }

  err = app.models.Token.DeleteByToken(requestPayload.Token)
  if err != nil {
    app.errorJSON(w, errors.New("invalid json 2"))
    return
  }

  payload := jsonResponse{
    Error:   false,
    Message: "logged out",
  }

  _ = app.writeJSON(w, http.StatusOK, payload)
}
