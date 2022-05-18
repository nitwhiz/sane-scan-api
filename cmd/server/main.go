package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nitwhiz/sane-scan-api/pkg/scanimage"
	"net/http"
	"os"
)

type ScanRequest struct {
	Format     string  `form:"format"`
	Resolution int     `form:"resolution"`
	Mode       string  `form:"mode"`
	Gamma      float64 `form:"gamma"`
}

func main() {
	gin.SetMode(os.Getenv(gin.EnvGinMode))

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/scan", func(c *gin.Context) {
		var req ScanRequest

		if err := c.ShouldBind(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		si := scanimage.New()

		si.Device = os.Getenv("SCAN_DEVICE")

		si.Format = req.Format
		si.Resolution = req.Resolution
		si.Mode = req.Mode
		si.Gamma = req.Gamma

		res, err := si.Scan()

		if err != nil {
			if _, ok := err.(*scanimage.ParameterError); ok {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}
		} else {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Content-Type", si.GetMimeType())
			c.String(200, res.String())
		}
	})

	err := r.Run("0.0.0.0:3000")

	if err != nil {
		panic(err)
	}
}
