package main

import (
	"errors"
	"log"
	"net/http"
	"net/mail"
	"os"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


type Server struct {
	DB *gorm.DB
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	db, err := Connect()
	if err != nil {
		log.Fatal(err)
	}

	server := &Server{
		DB: db,
	}	
	
	e := echo.New()

	// Routes
	e.POST("/register", server.createUser)


	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Fatal(e.Start(":" + port))
	
}

func (s *Server) createUser(c echo.Context) error {
	newUser := User {
		Username: c.FormValue("username"),
		Email: c.FormValue("email"),
	}
	// Check if email is valid
	_, err := mail.ParseAddress(newUser.Email)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, "invalid email format")
	}
	// Check if username already exists
	var userExists User
	if err := s.DB.Where("username = ?", newUser.Username).First(&userExists).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		return c.JSON(http.StatusBadRequest, "username already exists")
	}

	// Check if email already exists
	if err := s.DB.Where("email = ?", newUser.Email).First(&userExists).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	} else {
		return c.JSON(http.StatusBadRequest, "email already exists")
	}

	password := c.FormValue("password")
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	newUser.Password = string(hashedPassword)

	err = s.DB.Create(&newUser).Error
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	response := map[string]interface{}{
        "id":       newUser.ID,
        "username": newUser.Username,
        "email":    newUser.Email,
    }
	return c.JSON(http.StatusOK, response)

}


