module auth

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/jinzhu/gorm v1.9.12
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	github.com/stretchr/testify v1.5.1
	google.golang.org/grpc v1.29.1
	lib v0.0.0-00010101000000-000000000000
)

replace lib => ../lib
