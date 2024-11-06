package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sanitize/data"
)

type Controller struct {
	Database *data.SanitizeDB
}

// NewController example
func NewController(db *data.SanitizeDB) *Controller {
	return &Controller{Database: db}
}

// ListWords godoc
//
// @Summary		Sanitized Words
// @Description	Returns all the current words that will be used to sanitize text.
// @Tags			CRUD
// @Accept		json
// @Produce		json
// @Success		200	{object}   controller.SanitizeWord
// @Error        500
// @Router		/words [get]
func (c *Controller) ListWords(ctx *gin.Context) {
	operation, err := doCrudOperation(SELECT, SanitizeWord{}, c.Database)
	if err != nil {
		log.Printf("Error in doCrudOperation: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
	} else {
		ctx.JSON(200, operation)
	}
}

// AddWords godoc
//
// @Summary		Add Sanitized Words
// @Description	Provides the ability to add sanitized words. Returns a list of words that was successfully added.
// @Tags			CRUD
// @Accept		json
// @Produce		json
// @Param		account	body    controller.SanitizeWord	true "Add Sanitized Word"
// @Success		200	{object}   controller.SanitizeWord
// @Error       500
// @Router		/words [put]
func (c *Controller) AddWords(ctx *gin.Context) {
	var request SanitizeWord
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if len(request.Words) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	operation, err := doCrudOperation(INSERT, request, c.Database)
	if err != nil {
		log.Printf("Error in doCrudOperation: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
	} else {
		ctx.JSON(200, operation)
	}
}

// DeleteWords godoc
//
// @Summary		Remove Sanitized Words
// @Description	Provides the ability to add sanitized words. Returns a list of words that was successfully deleted.
// @Tags			CRUD
// @Accept		json
// @Produce		json
// @Param		account	body    controller.SanitizeWord	true "Remove Sanitized Word"
// @Success		200	{object}   controller.SanitizeWord
// @Error       500
// @Router		/words [delete]
func (c *Controller) DeleteWords(ctx *gin.Context) {
	var request SanitizeWord
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if len(request.Words) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	operation, err := doCrudOperation(DELETE, request, c.Database)
	if err != nil {
		log.Printf("Error in doCrudOperation: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
	} else {
		ctx.JSON(200, operation)
	}
}

// UpdateWords godoc
//
// @Summary		Update Sanitized Words
// @Description	Provides the ability to update an existing word(s). The first word in the list is the value that should be
// updated and the second is the value the first will update to. There is no limitation of the number of expect that all request
// must contain the first and second word. Returns all values that was updated to.
// @Tags		CRUD
// @Accept		json
// @Produce		json
// @Param		account	body    controller.SanitizeWord	true "Update Sanitized Word"
// @Success		200	{object}   controller.SanitizeWord
// @Error       500
// @Router		/words [post]
func (c *Controller) UpdateWords(ctx *gin.Context) {
	var request SanitizeWord
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if len(request.Words) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	operation, err := doCrudOperation(UPDATE, request, c.Database)
	if err != nil {
		log.Printf("Error in doCrudOperation: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
	} else {
		ctx.JSON(200, operation)
	}
}

// Sanitize godoc
//
// @Summary		Sanitize
// @Description	Provides the ability to sanitize strings based on the stored values. Returns the sanitized list in sequence the
// requests occurred.
// @Tags		Sanitize
// @Accept		json
// @Produce		json
// @Param		sanitize	body    controller.Sanitize	true "Sanitize Request"
// @Success		200	{object}   controller.Sanitize
// @Error       500
// @Router		/sanitize [post]
func (c *Controller) Sanitize(ctx *gin.Context) {
	var request Sanitize

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if len(request.Sentences) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	result, err := doSanitize(request, c.Database)
	if err != nil {
		log.Printf("Error in doSanitize: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
	} else {
		ctx.JSON(200, result)
	}
}
