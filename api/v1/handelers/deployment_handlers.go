package handlers

import (
	"net/http"

	"github.com /gin-gonic/gin"
	"github.com/ciliverse/cilikube/internal/service"
)

func ListDeployments(c *gin.Context) {
	deployments, err := service.ListDeployments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, deployments)
}

func CreateDeployment(c *gin.Context) {
	var deployment service.Deployment
	if err := c.ShouldBindJSON(&deployment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.CreateDeployment(&deployment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, deployment)
}

func GetDeployment(c *gin.Context) {
	name := c.Param("name")
	deployment, err := service.GetDeployment(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, deployment)
}

func UpdateDeployment(c *gin.Context) {
	name := c.Param("name")
	var deployment service.Deployment
	if err := c.ShouldBindJSON(&deployment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.UpdateDeployment(name, &deployment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, deployment)
}

func DeleteDeployment(c *gin.Context) {
	name := c.Param("name")
	if err := service.DeleteDeployment(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
