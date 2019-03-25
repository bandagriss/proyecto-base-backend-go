package auth

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"../../../database/models"
	"../../../lib/common"
	"golang.org/x/crypto/bcrypt"
)

// User es el alias para models.User
type User = models.User

func hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// checkHash ...
func checkHash(password string, hash string) bool  {
  err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateToken ...
func generateToken(data common.JSON) (string, error)  {
	// token es valido por 7 dias
	date := time.Now().Add(time.Hour * 24 * 7)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": data,
		"exp": date.Unix(),
	})

	// get path from root dir
	pwd, _ := os.Getwd()
	keyPath := pwd + "/jwtsecret.key"

	key, readErr := ioutil.ReadFile(keyPath)
	if readErr != nil {
		return "Ocurrio un error al leer el archivo jwtsecret ==>", readErr
	}
	tokenString, err := token.SignedString(key)
	return tokenString, err
}

func register(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	type RequestBody struct  {
		Username     string `json:"username" binding:"required"`
		DisplayName  string `json:"display_name" binding:"required"`
		Password  string `json:"password" binding:"required"`
	}

	var body RequestBody
	if err := c.BindJSON(&body); err != nil {
		fmt.Println("Error en el Body ====>", err)
		c.AbortWithStatus(400)
		return
	}

	var exists User
	if err := db.Where("username = ?", body.Username).First(&exists).Error; err == nil {
		fmt.Println("El usuario ya existe =>", err)
		c.AbortWithStatus(409)
		return
	}

	hash, hashErr := hash(body.Password)
	if hashErr != nil {
		fmt.Println("Ocurrio un error al generar la contraseña", hashErr)
		c.AbortWithStatus(500)
		return
	}

	// creando el usuario
	user := User {
		Username:     body.Username,
		DisplayName:  body.DisplayName,
		PasswordHash: hash,
	}

	db.NewRecord(user)
	db.Create(&user)

	serialized := user.Serialize()
	token, _ := generateToken(serialized)
	c.SetCookie("token", token, 60*60*24*7, "/", "", false, true)

	c.JSON(200, common.JSON{
		"user": user.Serialize(),
		"token": token,
	})
}

// login ...
func login(c *gin.Context)  {
	db := c.MustGet("db").(*gorm.DB)

	type RequestBody struct  {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var body RequestBody
	if err := c.BindJSON(&body); err != nil {
		fmt.Println("Ocurrió un error en el body de ingreso", err)
		c.AbortWithStatus(400)
	}

	var user User
	if err := db.Where("username = ?", body.Username).First(&user).Error; err != nil {
		fmt.Println("Usuario no encontrado =>", err)
		c.AbortWithStatus(400)
		return
	}

	if !checkHash(body.Password, user.PasswordHash) {
		fmt.Println("La contraseña es incorrecta")
		c.AbortWithStatus(401)
		return
	}

	serialized := user.Serialize()
	token, _ := generateToken(serialized)

	c.SetCookie("token", token, 60*60*24*7, "/", "", false, true)

	c.JSON(200, common.JSON{
		"user": user.Serialize(),
		"token": token,
	})
}

// check ...
func check(c *gin.Context)  {
  userRaw, ok := c.Get("user")
	if !ok {
		fmt.Println("No se logró obtener al usuario", ok)
		c.AbortWithStatus(401)
		return
	}

	user := userRaw.(User)

	tokenExpire := int64(c.MustGet("token_expire").(float64))
	now := time.Now().Unix()
	diff := tokenExpire - now

	fmt.Println(diff)
	if diff < 60*60*24*3 {
		// renovar token
		token, _ := generateToken(user.Serialize())
		c.SetCookie("token", token, 60*60*24*7, "/", "", false, true)
		c.JSON(200, common.JSON{
			"token": token,
			"user": user.Serialize(),
		})
		return
	}

	c.JSON(200, common.JSON{
		"token": nil,
		"user": user.Serialize(),
	})

}
