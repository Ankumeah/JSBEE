package apis

import (
	a "github.com/Ankumeah/JSBEE/backend/internal/app"
	"github.com/gin-gonic/gin"
)

func Apis(r *gin.RouterGroup, app *a.App) {
	ping(r)
	account(r, app)
}
