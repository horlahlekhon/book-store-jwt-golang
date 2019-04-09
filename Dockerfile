FROM golang:1.12

LABEL MAINTAINER="lekan.adebari@ubanquity.com"


# Install beego and the bee dev tool
RUN go get github.com/astaxie/beego && go get github.com/beego/bee

RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/jinzhu/gorm
RUN go get github.com/joho/godotenv
RUN go get github.com/gorilla/mux
RUN go get github.com/lib/pq
RUN go get golang.org/x/crypto/bcrypt

COPY . /go/src/book-store
# Expose the application on port 8080
EXPOSE 9000

# Set the entry point of the container to the bee command that runs the
# application and watches for changes
CMD ["bee", "run"]