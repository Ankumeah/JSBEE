package apis

import (
	"github.com/Ankumeah/JSBEE/backend/internal/database"

	"github.com/gin-gonic/gin"

	"errors"
	"net/http"
)

func handleError(c *gin.Context, err error) bool {
	if errors.Is(err, database.ErrInvalid) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return false
	} else if errors.Is(err, database.ErrExists) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return false
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return false
	}

	return true
}
