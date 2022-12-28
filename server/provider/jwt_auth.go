package provider

import (
	"log"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pankaj-katyare-wiz/airway-cargo-shipping-tracking/server/model"
	"github.com/pankaj-katyare-wiz/airway-cargo-shipping-tracking/server/repository"
	"github.com/pankaj-katyare-wiz/airway-cargo-shipping-tracking/server/request"
	"github.com/pankaj-katyare-wiz/airway-cargo-shipping-tracking/server/service"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const identityKey = "id"

type Success struct {
	Code   int    `json:"code" example:"200"`
	Expire string `json:"expire"`
	Token  string `json:"token"`
}

var once sync.Once

var mw *jwtAuthMiddleware

func NewJwtAuth(db *sqlx.DB) JwtAuthMiddleware {
	once.Do(func() {
		var err error

		mw = &jwtAuthMiddleware{
			databaseDriver: db,
		}

		mw.authMiddleware, err = jwt.New(mw.prepareMiddleware())

		if err != nil {
			log.Fatal("JWT error")
		}
	})

	return mw
}

type JwtAuthMiddleware interface {
	Middleware() *jwt.GinJWTMiddleware
	Refresh(c *gin.Context)
}

type jwtAuthMiddleware struct {
	databaseDriver *sqlx.DB
	authMiddleware *jwt.GinJWTMiddleware
}

func (mw *jwtAuthMiddleware) Middleware() *jwt.GinJWTMiddleware {
	return mw.authMiddleware
}

func (mw *jwtAuthMiddleware) prepareMiddleware() *jwt.GinJWTMiddleware {
	jwtSettings, err := service.NewJwtEnvVars()

	if err != nil {
		log.Fatal(err)
	}

	middleware := &jwt.GinJWTMiddleware{
		Realm:                 jwtSettings.Realm(),
		Key:                   []byte(jwtSettings.Secret()),
		Timeout:               jwtSettings.Expiration(),
		MaxRefresh:            jwtSettings.RefreshTime(),
		IdentityKey:           identityKey,
		PayloadFunc:           addUserIDToClaims,
		IdentityHandler:       extractIdentityKeyFromClaims,
		Authorizator:          mw.isUserValid,
		Authenticator:         mw.authenticate,
		HTTPStatusMessageFunc: takeAppropriateErrorMessage,
		TimeFunc:              time.Now,
	}

	return middleware
}

// authenticate godoc
//	@Summary		Authenticate a user
//	@Description	Perform user login
//	@ID				user-login
//	@Tags			User Actions
//	@Accept			json
//	@Produce		json
//	@Param			params	body		request.AuthRequest	true	"User's credentials"
//	@Success		200		{object}	Success
//	@Failure		401		{object}	response.Error
//	@Router			/login [post]
func (mw jwtAuthMiddleware) authenticate(c *gin.Context) (interface{}, error) {
	var authRequest request.AuthRequest
	var user model.Account

	if err := c.ShouldBind(&authRequest); err != nil {
		return user, jwt.ErrMissingLoginValues
	}

	userRepository := repository.NewUsersRepository(mw.databaseDriver)

	user, _ = userRepository.FindUserByEmail(authRequest.Email)

	if user.Name == "" || (bcrypt.CompareHashAndPassword([]byte(user.Name), []byte(authRequest.Password)) != nil) {
		return user, jwt.ErrFailedAuthentication
	}

	return user, nil
}

// refresh godoc
//	@Summary		Refresh token
//	@Description	Refresh user's token
//	@ID				refresh-token
//	@Tags			User Actions
//	@Produce		json
//	@Success		200	{object}	Success
//	@Failure		401	{object}	response.Error
//	@Security		ApiKeyAuth
//	@Router			/refresh [get]
func (mw jwtAuthMiddleware) Refresh(c *gin.Context) {
	mw.Middleware().RefreshHandler(c)
}

func (mw jwtAuthMiddleware) isUserValid(data interface{}, _ *gin.Context) bool {
	userID, ok := data.(float64)

	if !ok {
		return false
	}

	userRepository := repository.NewUsersRepository(mw.databaseDriver)

	return userRepository.FindUserByID(int(userID)).ID != 0
}

func extractIdentityKeyFromClaims(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)

	return claims[identityKey].(float64)
}

func addUserIDToClaims(data interface{}) jwt.MapClaims {
	if user, ok := data.(model.Account); ok {
		return jwt.MapClaims{
			identityKey: user.Id,
		}
	}

	return jwt.MapClaims{}
}

func takeAppropriateErrorMessage(err error, _ *gin.Context) string {
	switch err {
	case jwt.ErrMissingLoginValues:
		return "Email and password are required"
	case jwt.ErrFailedAuthentication:
		return "Invalid email or password"
	}

	return err.Error()
}
