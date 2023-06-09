package middlewares

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/mrehanabbasi/appraisal-system-backend/constants"
	log "github.com/mrehanabbasi/appraisal-system-backend/logger"
	"github.com/mrehanabbasi/appraisal-system-backend/models"
	"github.com/mrehanabbasi/appraisal-system-backend/utils"
)

func VerifyToken() gin.HandlerFunc {
	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	if err != nil {
		log.Error("failed to parse the issuer url: ", err.Error())
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to parse the issuer url"})
		}
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH0_AUDIENCE")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &models.TossClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Error("failed to set up the jwt validator: ", err.Error())
		return func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "failed to set up the jwt validator"})
		}
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Error("encountered error while validating jwt: ", err.Error())
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(c *gin.Context) {
		encounteredError := true
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			c.Request = r
			c.Next()
		}

		middleware.CheckJWT(handler).ServeHTTP(c.Writer, c.Request)

		if encounteredError {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid jwt token"})
			return
		}
	}
}

func ValidateJWTClaims(c *gin.Context) {
	// Getting claims
	claims, ok := c.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	if !ok {
		err := errors.New("failed to validate jwt claims")
		log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tossClaims, ok := claims.CustomClaims.(*models.TossClaims)
	if !ok {
		err := errors.New("failed to cast custom jwt claims to the desired type")
		log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	empID, _ := strconv.ParseUint(tossClaims.EmpID, 10, 16)
	designationID, _ := strconv.ParseUint(tossClaims.Designation, 10, 16)
	supID, _ := strconv.ParseUint(tossClaims.Supervisor, 10, 16)
	roleID, _ := strconv.ParseUint(tossClaims.Role, 10, 16)
	supName, err := utils.GetSupervisorName(uint16(supID))
	if err != nil {
		log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tokenInfo := models.TokenInfo{
		EmpID:         uint16(empID),
		EmailAddr:     tossClaims.EmailAddr,
		DesignationID: uint16(designationID),
		Designation:   tossClaims.EmpDesignation,
		SupervisorID:  uint16(supID),
		Supervisor:    supName,
		Department:    tossClaims.Department,
		EmpImagePath:  tossClaims.EmpImagePath,
		EmpRoleID:     uint16(roleID),
	}

	c.Set(constants.TOKEN_DATA, tokenInfo)
}
