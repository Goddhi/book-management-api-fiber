package main

///
import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/goddhi/book-api-fiber/models"
	"github.com/goddhi/book-api-fiber/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

// the logic function
func(r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{} /// initializing the instance

	err := context.BodyParser(&book) // mapping request data from json  to Go structs
	if err != nil{
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{
		"messages": "request failed",
	})
		return err
	}

	/// creating book to the database
	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message":"could not create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
	"message": "book has been added",
	})
	return nil

}

func (r *Repository) Getbooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not get books"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "books fetched successfully", "data": bookModels, })
	return nil
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {

	id := context.Params("id")
	bookModel := &models.Books{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
		"message": "id can't be empty",
	})
		return nil
	}
	fmt.Println("the ID is", id)


	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
		"message": "book id cant be fetched",
	})
		return err
	}
	
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book id fetched successfully", 
		"data": bookModel,
})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id") // gettimg value of the params
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "id can not be empty"})
		return nil
	}

	err := r.DB.Delete(bookModel, id)
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"messages": "cannot delete book id"})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "book id ferched successfully"})
	return nil
}


func (r *Repository) SetupRoutes(app *fiber.App){
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.Getbooks)
}
//
type Book struct {
	Author string `json:"author"`
	Title string  `json:"title"`
	Publisher string `json:"publisher"`
}


func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
 		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_DBNAME"),
		SSLMode: os.Getenv("DB_SSLMODE"),


	}

	db, err := storage.NewConnection(config) // set up connectivity with the database
	if err != nil {
		log.Fatal("can't connect with the database") 
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB : db,   //
	}
	app := fiber.New() //  creates a new instance of the Fiber web application
	// This function's purpose is to register various routes (URL patterns and their associated handler functions) with the application, enabling it to process incoming requests correctly
	r.SetupRoutes(app) // This line calls a function named r.SetupRoutes to configure the application's routes
	app.Listen(":8080")


}
//