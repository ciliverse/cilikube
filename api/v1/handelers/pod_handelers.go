package handlers

import (
	"net/http"

	"github.com/ciliverse/cilikube/internal/service"
	"github.com/gin-gonic/gin"
)

func ListPods(c *gin.Context) {
	pods, err := service.ListPods()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pods)
}

func CreatePod(c *gin.Context) {
	var pod service.Pod
	if err := c.ShouldBindJSON(&pod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.CreatePod(&pod); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pod)
}

func GetPod(c *gin.Context) {
	name := c.Param("name")
	pod, err := service.GetPod(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pod)
}

func UpdatePod(c *gin.Context) {
	name := c.Param("name")
	var pod service.Pod
	if err := c.ShouldBindJSON(&pod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := service.UpdatePod(name, &pod); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, pod)
}

func DeletePod(c *gin.Context) {
	name := c.Param("name")
	if err := service.DeletePod(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
