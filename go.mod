module booksapi

go 1.16

replace utils => ./utils

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/lib/pq v1.10.2
	utils v0.0.0-00010101000000-000000000000
)
