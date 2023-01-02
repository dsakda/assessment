package main

import (
	"log"
	"os"

	"github.com/dsakda/assessment/expense"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	err := godotenv.Load("dev.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	expense.InitDB()

	e := echo.New()

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		AuthScheme: "November",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == "10, 2009", nil
		},
	}))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/expenses", expense.GetAllExpensesHandler)
	e.POST("/expenses", expense.CreateExpenseHandler)
	e.GET("/expenses/:id", expense.GetExpenseHandler)
	e.PUT("/expenses/:id", expense.UpdateExpenseHandler)

	e.Logger.Fatal(e.Start(os.Getenv("PORT")))
}
