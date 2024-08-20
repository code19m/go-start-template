package http

import (
	"go-start-template/internal/domain"
	"go-start-template/pkg/errx/errto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createMyModelReqBody struct {
	Name string `json:"name" binding:"required"`
	Age  int32  `json:"age"  binding:"required"`
}

// @Router /my-model [post]
// @Tags my-model
// @Param payload body createMyModelReqBody true "_"
func (h *HttpServer) createMyModelHandler(c *gin.Context) {
	var reqBody createMyModelReqBody

	err := bindAndValidate(c, &reqBody)
	if err != nil {
		errto.HTTP(c.Writer, err)
		return
	}

	id, err := h.myModelSrv.Create(c, domain.CreateMyModelParams{
		Name: reqBody.Name,
		Age:  reqBody.Age,
	})
	if err != nil {
		errto.HTTP(c.Writer, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"id":      id,
	})
}
