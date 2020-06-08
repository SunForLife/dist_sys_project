module app

go 1.13

require (
	github.com/jinzhu/gorm v1.9.12
	google.golang.org/grpc v1.29.1
	lib v0.0.0-00010101000000-000000000000
)

replace lib => ../lib
