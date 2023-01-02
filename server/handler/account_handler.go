package handler

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pankaj-katyare-wiz/airway-cargo-shipping-tracking/server/model"
	"github.com/pankaj-katyare-wiz/airway-cargo-shipping-tracking/server/repository"
	"github.com/pankaj-katyare-wiz/airway-cargo-shipping-tracking/server/response"
	"github.com/pankaj-katyare-wiz/airway-cargo-shipping-tracking/server/utils"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	DB      *sqlx.DB
	queries *repository.Queries
}

func NewAccountHandler(DB *sqlx.DB) *AccountHandler {
	return &AccountHandler{
		DB:      DB,
		queries: repository.New(DB),
	}
}

func (handler AccountHandler) CreateAccount(context *gin.Context) {

	var account model.Account

	if err := context.ShouldBind(&account); err != nil {
		response.ErrorResponse(context, http.StatusBadRequest, "Required fields are empty")
		fmt.Println("error", err)
		return
	}

	state, err := handler.queries.CreateAccountDetails(context, repository.CreateAccountDetailsParams{
		ID:          uuid.New().String(),
		Name:        sql.NullString{String: account.Name, Valid: true},
		Email:       sql.NullString{String: account.Email, Valid: true},
		CompanyName: sql.NullString{String: account.CompanyName, Valid: true},
		Mobile:      sql.NullString{String: account.Mobile, Valid: true},
		Roles:       sql.NullString{String: account.Roles, Valid: true},
		City:        sql.NullString{String: account.City, Valid: true},
		Password:    sql.NullString{String: utils.GetHash(account.Password), Valid: true},
	})

	if err != nil {
		fmt.Println("error", err)
		response.ErrorResponse(context, http.StatusBadRequest, "Error inserting account")
		return
	}

	response.SuccessResponse(context, map[string]interface{}{
		"code":    "success",
		"message": "Request account Created successfuly",
		"data":    state,
	})
}

func (handler AccountHandler) Login(context *gin.Context) {

	var login model.LoginRequest

	if err := context.ShouldBind(&login); err != nil {
		response.ErrorResponse(context, http.StatusBadRequest, "Required fields are empty")
		fmt.Println("error", err)
		return
	}

	if login.Email == "" {
		response.ErrorResponse(context, http.StatusNotFound, "ID not specified")
		return
	}
	if login.Password == "" {
		response.ErrorResponse(context, http.StatusNotFound, "password not specified")
		return
	}

	state, err := handler.queries.GetAccountDetails(context, login.Email)

	if err != nil {
		response.ErrorResponse(context, http.StatusNotFound, "Error on get account details")
		fmt.Println("error", err)
		return
	}

	if state.Email.String != login.Email && !utils.IsHashValid(state.Password.String, login.Password) {
		response.ErrorResponse(context, http.StatusNotFound, "id or password are wrong")
		return
	}

	response.SuccessResponse(context, map[string]interface{}{
		"code":    "success",
		"message": "login successful",
		"data":    state,
	})
}

func (handler AccountHandler) GetAccountByID(context *gin.Context) {

	id := context.Request.URL.Query().Get("id")
	if id == "" {
		response.ErrorResponse(context, http.StatusNotFound, "ID not specified")
		return
	}

	state, err := handler.queries.GetAccountDetails(context, id)

	if err != nil {
		response.ErrorResponse(context, http.StatusNotFound, err.Error())
		return
	}

	response.SuccessResponse(context, map[string]interface{}{
		"code":    "success",
		"message": "account Data",
		"data":    state,
	})
}

func (handler AccountHandler) UpdateAccountDetails(context *gin.Context) {

	var account model.Account

	if err := context.ShouldBind(&account); err != nil {
		response.ErrorResponse(context, http.StatusBadRequest, "Required fields are empty")
		return
	}

	data, err := handler.queries.GetAccountDetails(context, account.Id)

	if err != nil {
		response.ErrorResponse(context, http.StatusNotFound, "Error on get account details")
		fmt.Println("error", err)
		return
	}

	fmt.Println("old data :", data)

	err = handler.queries.UpdateAccountDetails(context, repository.UpdateAccountDetailsParams{
		ID:          account.Id,
		Name:        sql.NullString{String: account.Name, Valid: true},
		CompanyName: sql.NullString{String: account.CompanyName, Valid: true},
		Email:       sql.NullString{String: account.Email, Valid: true},
		Mobile:      sql.NullString{String: account.Mobile, Valid: true},
		Roles:       sql.NullString{String: account.Roles, Valid: true},
		City:        sql.NullString{String: account.City, Valid: true},
		Password:    sql.NullString{String: utils.GetHash(account.Password), Valid: true},
	})

	if err != nil {
		response.ErrorResponse(context, http.StatusNotFound, "Error on Update account details")
		fmt.Println("error", err)
		return
	}
	// TODO return, nothing to update
	response.SuccessResponse(context, map[string]interface{}{
		"code":    "success",
		"message": "Updated suceessfully",
	})
}

func (handler AccountHandler) GetAllAccount(context *gin.Context) {

	quotes, err := handler.queries.ListAccountDetails(context)

	if err != nil {
		response.SuccessResponse(context, map[string]interface{}{
			"code":    "success",
			"message": "Error int get all account",
			"error":   err,
		})
		return
	}

	response.SuccessResponse(context, map[string]interface{}{
		"code":    "success",
		"message": "Fetched all account list",
		"data":    quotes,
	})
}
