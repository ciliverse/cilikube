package handlers

import (
	"net/http"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/gin-gonic/gin"
)

func ListServices(c *gin.Context) {
	services, err := service.ListServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, services)
}

func CreateService(c *gin.Context) {
	var svc service.Service
	if err := c.ShouldBindJSON(&svc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.CreateService(&svc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, svc)
}

func GetService(c *gin.Context) {
	name := c.Param("name")
	svc, err := service.GetService(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, svc)
}

func UpdateService(c *gin.Context) {
	name := c.Param("name")
	var svc service.Service
	if err := c.ShouldBindJSON(&svc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.UpdateService(name, &svc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, svc)
}

func DeleteService(c *gin.Context) {
	name := c.Param("name")
	if err := service.DeleteService(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
