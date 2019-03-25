package posts

import (
	"fmt"
	
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"../../../database/models"
	"../../../lib/common"
)

type Post = models.Post
type User = models.User

type JSON = common.JSON

// create ...
func create(c *gin.Context) {
  db := c.MustGet("db").(*gorm.DB)
	type RequestBody struct {
		Text string `json:"text" binding:"required"`
	}
	var requestBody RequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		fmt.Println("No se encontro ningun valor en el body =>", err)
		c.AbortWithStatus(400)
		return
	}

	user := c.MustGet("user").(User)
	post := Post{Text: requestBody.Text, User: user}
	db.NewRecord(post)
	db.Create(&post)
	c.JSON(200, post.Serialize())
}

func list(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	cursor := c.Query("cursor")
	recent := c.Query("recent")

	var posts []Post

	if cursor == "" {
		if err := db.Preload("User").Limit(10).Order("id desc").Find(&posts).Error; err != nil {
			c.AbortWithStatus(500)
			return
		}
	} else {
		condition := "id < ?"
		if recent == "1" {
			condition = "id > ?"
		}
		if err := db.Preload("User").Limit(10).Order("id desc").Where(condition, cursor).Find(&posts).Error; err != nil {
			fmt.Println("Ocurrio un error al listar ==>", err)
			c.AbortWithStatus(500)
			return
		}
	}

	length := len(posts)
	serialized := make([]JSON, length, length)

	for i := 0; i < length; i++ {
		serialized[i] = posts[i].Serialize()
	}

	c.JSON(200, serialized)
}

// read ...
func read(c *gin.Context)  {
  db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	var post Post

	if err := db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&post).Error; err != nil {
		fmt.Println("No se encontro el post", err)
		c.AbortWithStatus(404)
		return
	}
	c.JSON(200, post.Serialize())
}

// remove ...
func remove(c *gin.Context)  {
  db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	user := c.MustGet("user").(User)

	var post Post
	if err := db.Where("id = ?", id).First(&post).Error; err != nil {
		fmt.Println("Ocurrio un error al eliminar ==>", err)
		c.AbortWithStatus(404)
		return
	}

	if post.UserID != user.ID {
		fmt.Println("No se encontro al usuario ")
		c.AbortWithStatus(403)
		return
	}

	db.Delete(&post)
	c.Status(204)
}

// update ...
func update(c *gin.Context)  {
  db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	user := c.MustGet("user").(User)

	type RequestBody struct  {
		Text string `json:"text" binding:"required"`
	}

	var requestBody RequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		fmt.Println("Ocurrio un error al actualizar ==>", err)
		c.AbortWithStatus(400)
		return
	}

	var post Post
	if err := db.Preload("User").Where("id = ?", id).First(&post).Error; err != nil {
		fmt.Println("Ocurrio un error al acutalizar el post user => ", err)
		c.AbortWithStatus(404)
		return
	}

	if post.UserID != user.ID {
		fmt.Println("el usuario asociado al post no se encontro ")
		c.AbortWithStatus(403)
		return
	}

	post.Text = requestBody.Text
	db.Save(&post)
	c.JSON(200, post.Serialize())
}
