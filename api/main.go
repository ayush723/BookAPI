package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"utils"
	// "net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

//Book holds the information for a book
type Book struct{
	ID int `json:"id"`
	Name string `json:"name"`
	Author string `json:"author"`
	Year string `json:"year"`
}

//setting driver visible to wholeprogram
var db *sql.DB


func main()  {
	var err error
	//connecting to the database
	db, err = sql.Open("postgres", "postgres://ayush:envy@localhost/book") //("driver", "driver://user:passsword@host/dbname")
	if err!= nil{
		panic(err)
	}
	//initializing tables
	utils.Initialize(db)
	defer db.Close()
	//
	r := gin.Default()

	//adding routes to REST verbs
		r.GET("/books/",getBooks)
		r.GET("/books/:book-id/", getBook)
		r.POST("/books/",addBook)
		r.PUT("/books/:book-id",updateBook)
		r.DELETE("/books/:book-id/",deleteBook)
		r.Run(":8000")
		
	}

//GET http://localhost:8000/books/
//getBooks gets all the book 
func getBooks(c *gin.Context)  {
	

	var books []Book
	var book Book
	rows, err := db.Query("select * from books_new;")
	if err!=nil{
		 panic(err)
	}
	defer rows.Close()
	for rows.Next(){
		err = rows.Scan(&book.ID, &book.Name, &book.Author, &book.Year)
		books = append(books,book)
	}


	if err!=nil{
		log.Println(err)
		c.JSON(500, gin.H{
			"error":err.Error(),
		})
	} else{
		c.JSON(200, gin.H{
			"result": books,
		})
	}
}
//GET http://localhost:8000/books/:book-id
//getBook gets a specific book by using id
func getBook(c *gin.Context)  {
		
		id := c.Param("book-id")
		var tmp Book
		err := db.QueryRow("select * from books_new where id =$1",id).Scan(&tmp.ID, &tmp.Name,&tmp.Author,&tmp.Year)
		fmt.Println(tmp)
		if err!= nil{
			log.Println(err)
			c.JSON(500, gin.H{
				"error":err.Error(),
			})
			} else{
				c.JSON(200,gin.H{
					"result":tmp,
				})
			}
}
//POST http://localhost:8000/books/
//addBook adds new book into the database
func addBook(c *gin.Context)  {

		var book Book

		if err := c.BindJSON(&book); err ==nil{
			_, err := db.Exec("insert into books_new (ID, NAME, AUTHOR, YEAR) VALUES ($1,$2,$3,$4)",book.ID,book.Name, book.Author,book.Year)
			if err == nil{
				
				c.JSON(200,gin.H{
					"result":book,
					}) 
				}else{
					c.JSON(500,gin.H{ 
						"error":err.Error(),
					})
				}
			}
		}
	
//PUT http://localhost:8000/books/:book-id
//updateBook updates any existing book using the id
func updateBook(c *gin.Context)  {
			
			id := c.Param("book-id")
			var book Book
			json.NewDecoder(c.Request.Body).Decode(&book)

			
			result, err := db.Exec("update books_new set id=$1, name=$2, author=$3, year=$4 where id=$5 returning id",&book.ID,&book.Name,&book.Author,&book.Year,id)
			if err!=nil{
				c.JSON(500,gin.H{
					"error":err.Error(),
				})
			}
			rowsUpdated,err := result.RowsAffected()
			if err!=nil{
				c.JSON(500,gin.H{
					"error":err.Error(),
				})
				} else {
					c.JSON(200, gin.H{
						"result":book,
						"rowsaffected":rowsUpdated,
					})
				}
				if book.ID == 0 || book.Author == "" || book.Name == "" || book.Year == "" {
					
					c.JSON( 400, gin.H{
						"error":"enter missing fields",
					})
				
				}
}
		
//DELETE http://localhost:8000/books/:book-id
//deleteBook deletes existing book using id
func deleteBook(c *gin.Context)  {

	
	id := c.Param("book-id")
	
	_, err := db.Query("delete from books_new where id=$1",id)
	if err!= nil{
		panic(err)
		} else{
			c.JSON(200, gin.H{
				"result":`succesfully deleted`,
			})
		}
}