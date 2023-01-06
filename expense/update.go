package expense

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func UpdateExpenseHandler(c echo.Context) error {

	id := c.Param("id")
	e := Expense{}
	err := c.Bind(&e)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, Err{Message: "ERROR, Can't bind object."})
	}

	stmt, err := db.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1;")
	if err != nil {
		log.Println("can't prepare statment update", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: "ERROR, Can't prepare statment to update object."})
	}

	if _, err := stmt.Exec(id, e.Title, e.Amount, e.Note, pq.Array(&e.Tags)); err != nil {
		log.Println("error execute update ", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: "ERROR, Can't update object."})
	}

	return c.JSON(http.StatusOK, e)
}
