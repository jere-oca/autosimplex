// Draft: este handler se usará en futuros Issues
func SimplexHandler(c *gin.Context) {
	var req models.SimplexRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "JSON inválido"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"received": req})
}
