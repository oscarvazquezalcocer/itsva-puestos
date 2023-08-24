package handlers

import (
	"fmt"
	"itsva-puestos/models"
	"itsva-puestos/services"
	"itsva-puestos/utils"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func List(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var puestos []models.Puesto
	db.Find(&puestos)

	var nombresJefes = make(map[uint]string)
	for _, puesto := range puestos {
		if puesto.IDJefe != 0 {
			var puestoPadre models.Puesto
			db.First(&puestoPadre, puesto.IDJefe)
			nombresJefes[puesto.IDJefe] = puestoPadre.Nombre
		}
	}

	puestosWithDetails, err := utils.GetPuestosWithDetails(puestos)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}

	fmt.Println(puestosWithDetails)
	c.HTML(http.StatusOK, "list.html", gin.H{"puestos": puestosWithDetails, "nombresJefes": nombresJefes})
}

func ShowForm(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	funciones, err := services.GetListFuncionFromAPI()

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var jefes []models.Puesto
	db.Find(&jefes)

	c.HTML(http.StatusOK, "create.html", gin.H{"jefes": jefes, "funciones": funciones})

}

func Create(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var newPuesto models.Puesto
	if err := c.ShouldBind(&newPuesto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Checamos que el puesto que se reciba sea valido
	if newPuesto.IDJefe != 0 {
		var parent models.Puesto
		result := db.First(&parent, newPuesto.IDJefe)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ParentID no válido"})
			return
		}
	}

	if err := db.Create(&newPuesto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//c.JSON(http.StatusOK, newUser)
	c.Redirect(http.StatusSeeOther, "/")
}

func Show(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	var puesto models.Puesto
	result := db.First(&puesto, id)
	if result.Error != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"message": "Puesto no encontrado"})
		return
	}

	// Si es el director no enviamos a los jefes
	if puesto.ID == 1 {
		c.HTML(http.StatusOK, "show.html", gin.H{"puesto": puesto})
		return

	}

	funciones, err := services.GetListFuncionFromAPI()

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	var jefes []models.Puesto
	db.Find(&jefes)

	c.HTML(http.StatusOK, "show.html", gin.H{"puesto": puesto, "jefes": jefes, "funciones": funciones})
}

func Update(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	var puesto models.Puesto
	result := db.First(&puesto, id)
	if result.Error != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"message": "Puesto no encontrado"})
		return
	}
	if err := c.ShouldBind(&puesto); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{"message": err.Error()})
		return
	}
	fmt.Println(puesto.Repetible)

	db.Save(&puesto)
	c.Redirect(http.StatusSeeOther, "/")
}

func Delete(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	var puesto models.Puesto
	result := db.First(&puesto, id)
	if result.Error != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{"message": "Puesto no encontrado"})
		return
	}
	db.Delete(&puesto)
	c.Redirect(http.StatusSeeOther, "/")
}

func TreeView(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var puestos []models.Puesto
	db.Find(&puestos)

	topLevel := utils.RenderTree(puestos, 0) // 0 representa el jefe raíz

	//c.JSON(http.StatusOK, gin.H{"tree": topLevel})
	c.HTML(http.StatusOK, "tree.html", gin.H{"tree": topLevel})
}
