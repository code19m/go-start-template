package http

import (
	"go-start-template/pkg/errx"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func bindAndValidate(c *gin.Context, reqBody interface{}) error {
	err := c.ShouldBindJSON(reqBody)
	if err == nil {
		return nil
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		appErr := errx.ErrValidation
		for _, fe := range errs {
			appErr = appErr.WithDetail(fe.Field(), fe.Tag())
		}
		return appErr
	}

	return err
}

// func handleBindErr(c *gin.Context, err error) bool {
// 	if c.ContentType() != gin.MIMEJSON {
// 		c.JSON(http.StatusUnsupportedMediaType, gin.H{
// 			"type":    apperr.UnsupportedPayloadType,
// 			"message": "Invalid content type",
// 		})
// 		c.Error(err) //nolint: errcheck
// 		return true
// 	}

// 	if err == nil {
// 		return false
// 	}

// 	// Handle ValidationErrors
// 	if errs, ok := err.(validator.ValidationErrors); ok { //nolint: errorlint
// 		invalidArgs := parseValidationErrors(errs)
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"type":    apperr.Validation,
// 			"message": "Invalid input params",
// 			"args":    invalidArgs,
// 		})
// 		c.Error(err) //nolint: errcheck
// 		return true
// 	}

// 	c.JSON(http.StatusBadRequest, gin.H{
// 		"type":    apperr.Validation,
// 		"message": err.Error(),
// 	})
// 	c.Error(err) //nolint: errcheck
// 	return true
// }

// func handleAppErr(c *gin.Context, err error) bool {
// 	if err == nil {
// 		return false
// 	}

// 	status := getStatus(err)
// 	if status == http.StatusInternalServerError {
// 		err = apperr.NewAppError(apperr.Internal, err.Error())
// 		c.Error(err) //nolint: errcheck
// 	}

// 	if !c.Writer.Written() {
// 		c.JSON(status, err)
// 		c.Error(err) //nolint: errcheck
// 	}

// 	return true
// }

// type invalidArgument struct {
// 	Field string `json:"field"`
// 	Tag   string `json:"tag"`
// }

// // parseValidationErrors converts BindJson errors to user friendly format
// func parseValidationErrors(errs validator.ValidationErrors) []invalidArgument {
// 	var invalidArgs []invalidArgument

// 	for _, err := range errs {
// 		invalidArgs = append(invalidArgs, invalidArgument{
// 			err.Field(),
// 			err.Tag(),
// 		})
// 	}
// 	return invalidArgs
// }

// func getStatus(err error) int {
// 	var e *apperr.AppError
// 	if errors.As(err, &e) {
// 		switch e.Type {
// 		case apperr.Authorization:
// 			return http.StatusUnauthorized
// 		case apperr.Forbidden:
// 			return http.StatusForbidden
// 		case apperr.Validation:
// 			return http.StatusBadRequest
// 		case apperr.NotFound:
// 			return http.StatusNotFound
// 		case apperr.Conflict:
// 			return http.StatusConflict
// 		case apperr.PayloadTooLarge:
// 			return http.StatusRequestEntityTooLarge
// 		case apperr.UnsupportedPayloadType:
// 			return http.StatusUnsupportedMediaType
// 		case apperr.TimeOut:
// 			return http.StatusGatewayTimeout
// 		case apperr.Internal:
// 			//
// 		}
// 	}
// 	return http.StatusInternalServerError
// }
