package expense

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func GetExpenseHandler(c echo.Context) error {
	id := c.Param("id")
	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses WHERE id = $1")
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't prepare query expense statment."})
	}

	row := stmt.QueryRow(id)
	e := Expense{}
	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, Err{Message: "expense not found"})
	case nil:
		return c.JSON(http.StatusOK, e)
	default:
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, Err{Message: "can't scan expense."})
	}
}

func GetAllExpensesHandler(c echo.Context) error {
	stmt, err := db.Prepare("SELECT id, title, amount, note, tags FROM expenses order by id")
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't prepare query all expenses statment."})
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, Err{Message: "Can't query all expenses."})
	}

	expenses := []Expense{}

	for rows.Next() {
		e := Expense{}
		err := rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
		if err != nil {
			log.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, Err{Message: "Can't scan expense."})
		}
		expenses = append(expenses, e)
	}

	return c.JSON(http.StatusOK, expenses)
}
