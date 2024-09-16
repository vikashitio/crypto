package main

import (
	"fmt"
	"html/template"
	"log"
	"template/database"
	"template/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

func init() {
	database.ConnectDb() // for connect Db define in function.go
}

func main() {

	// set engine for frontend file
	//engine := html.New("./views", ".html").Funcs(funcMap)

	// Load templates with custom functions
	// Define custom functions
	funcMap := template.FuncMap{
		"sub": func(a, b int) int {
			return a - b
		},
		"add": func(a, b int) int {
			return a + b
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},
	}

	// Load templates with custom functions
	engine := html.New("./views", ".html")
	engine.AddFuncMap(funcMap)

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})
	// Apply CORS middleware globally
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",                            // Allow all origins or specify certain domains
		AllowMethods:     "GET, POST, PATCH, DELETE",     // Allow only GET requests
		AllowHeaders:     "Origin, Content-Type, Accept", // Allow necessary headers
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))
	app.Static("/views", "./views")   // allow folder for get public files like js/css etc
	app.Static("/assets", "./assets") // allow folder for get public files like js/css etc
	app.Use(logger.New())

	routes.InitRoutes(app) // for Set page routes path define in routes.go

	// Handle 404 errors with a custom HTML page
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).Render("404", fiber.Map{})
	})

	// check .env file is exist or not
	if err := godotenv.Load(".env"); err != nil {
		fmt.Printf("ENV not Found")
		return
	}

	log.Fatal(app.Listen(":3000"))
	//app.Listen(":3000") // Set port for run site

}
