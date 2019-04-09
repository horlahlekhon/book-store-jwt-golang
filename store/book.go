package store

import (
	"database/sql"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Book struct {
	Id int `json:"id"`
	Name *string `json:"name"`
	Isbn *uint `json:"isbn"`
	Price *int `json:"price"`

}

func Database() *sql.DB {


	uri := fmt.Sprintf("host=database dbname=book_store user=postgres sslmode=disable password=postgres")

	db, error := sql.Open("postgres", uri)

	if error != nil {
		fmt.Println(error)
	}
	return db

}

func (book *Book)CreateBook()  interface{} {
	statemnt := `INSERT INTO book (id, name, isbn, price) VALUES($1,$2,$3,$4) RETURNING id`

	id := 0

	errore := Database().QueryRow(statemnt, &book.Id,&book.Name, &book.Isbn,&book.Price).Scan(&id)

	if errore != nil {
		fmt.Println(errore)
		panic(errore.Error())
	}

	row , errorr :=  Database().Query("SELECT id, name, isbn, price FROM book WHERE id=$1", id)

	if errorr != nil{
		fmt.Println(errorr)
		panic(errorr.Error())
	}
	var bok Book
	for row.Next(){
		if err := row.Scan(&bok.Id, &bok.Name, &bok.Price, &bok.Isbn); err != nil {
			fmt.Println(err)
		}
	}

	return bok


}

func GetBooks() []*Book {
	fmt.Println("here")
	rows, error := Database().Query("SELECT * FROM book")

	if error != nil {
		fmt.Println(error)
		panic(error.Error())
	}

	books := []*Book{}
	
	for rows.Next() {
		var id int
		var Name string
		var isbn uint
		var price int
		if error := rows.Scan(&id,&Name,&isbn,&price); error != nil {
			panic(error.Error())

		}
		books = append(books, &Book{id, &Name,&isbn, &price})
	}
	//fmt.Println(books)
	return books
}

func UpdateBook(book Book) map[string]interface{}  {
	statement := `UPDATE book SET name=$1, isbn=$2, price=$3 WHERE id=$4`

	result, error := Database().Exec(statement, book.Name,book.Isbn, book.Price, book.Id)

	if error != nil {
		fmt.Println(error)
		panic(error)
	}
	count , error := result.RowsAffected()

	if error != nil {
		fmt.Println(error)
		panic(error)
	}
	updatedBook := GetBook(book.Id)

	resp := map[string]interface{}{
		"data updated" : count,
		"data" : updatedBook,
	}

	return resp

}

func GetBook(id int) Book {

	statement := `SELECT * FROM book WHERE id=$1`

	result , error := Database().Query(statement,id)

	if error != nil{
		fmt.Println(error)
		panic(error)
	}
	book := Book{}
	for result.Next() {
		var id int
		var name string
		var isbn uint
		var price int
		if err := result.Scan(&id, &name, &isbn, &price); err != nil {
			fmt.Println(err)
			panic(error)
		}
		book = Book{id, &name , &isbn, &price}
	}


	return book
}

func DeleteBookById(id int) map[string]int64 {
	statement := `DELETE FROM book WHERE id=$1`

	result , error := Database().Exec(statement, id)
	if error!=nil {
		fmt.Println(error)
		panic(error)
	}

	count , erro := result.RowsAffected()
	if erro != nil {
		fmt.Println(error)
		panic(error)
	}

	return map[string]int64{"data affected": count}
}





