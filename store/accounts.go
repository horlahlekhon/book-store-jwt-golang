package store

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
)

type token struct {
	UserID string `json:"token_id"`
	UserRole string `json:"user_role"`
	jwt.StandardClaims
}

type Account struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Password string `json:"password"`
	token string `json:"token"`
	Role string `json:"role"`
}

const(
	PRIVATE_KEY = "keys/jwt_rsa.key"
	PUBLIC_KEY = "keys/jwt_rsa.key.pub"
)

func SignKey(keyPath string)[]byte  {

	SignKey, error := ioutil.ReadFile(keyPath)
	if error != nil {
		log.Fatal("Error reading private key")
		return nil
	}
	return SignKey
}

func VerifyKey( keyPath string) []byte {


	VerifyKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("Error reading public key")
		return nil
	}
	return VerifyKey
}


func GetAccount(name string) (*Account,  string) {

	if name == "" {
		return nil, "Please enter a non Empty name that is valid"
	}

	statement := `SELECT * FROM accounts WHERE name = $1`
	result , err := Database().Query(statement, name)
	if err != nil {
		log.Fatal(err)
		return nil, err.Error()
	}

	 account := Account{}

	for result.Next(){
		var id int
		var name string
		var password string
		var token string
		var role string
		if error := result.Scan(&id,&name, &password, &token, &role); error != nil {
			fmt.Println(error)
			return nil, error.Error()
		}
		account = Account{id, name, password, token, role}
	}


	 return &account, ""
}

func (account *Account) Validate() (string ,bool ) {

	if len(account.Password) < 6 {
		return "password is short", false
	}

	if account.Name == "" {
		return "kindly find in your heart to add a name", false
	}

	account , error := GetAccount(account.Name)
	if error == "There is a user with that Name, pls choose another" && account != nil {
		return error, false
	}


	return "Requirement is passed", true

}

func RegisterAccount(account Account) (bool,map[string]interface{})  {
	if resp, ok := account.Validate(); !ok {
		return  false, map[string]interface{}{"response": resp}
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	//create a new token for the new user, we can set additionl private claims inside the token and pass it to the jwt.NewWithClaims like that to
	//create the token string
	 tk := token{UserID:account.Name, UserRole:account.Role}


	 //here we passed the token which includes our claims to the newWithClaim jwt function and specify an algorithm to be used for signing the jwt
	 //and remember to not use the literal signing i.e GetSigningMethod("RSA"), but this : jwt.SigningMethodRS256
	 signer := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	 //this will get the token string taht has been signed with the key we provided and the method we put, it will return the signed string token
	 tokenString, error := signer.SignedString([]byte("lekan"))
	 fmt.Println("token string ", tokenString)

	if error != nil {
		fmt.Println("error ", error)
		return false, map[string]interface{}{"response":error}
	}
	 //put the token inside the account as its token
	account.token = tokenString


	/////persist
	statement := `INSERT INTO accounts(name, password ,token, role) VALUES($1,$2,$3,$4)`

	_, errore := Database().Exec(statement, &account.Name, &account.Password, &account.token,&account.Role)
	if errore != nil {
		return false, map[string]interface{}{"response": errore}
	}

	return true, map[string]interface{}{"response":  "succesfully created your account"}
	
}

func Login(name, password string) (bool,string, string ) {


	account , err := GetAccount(name)
	if err != "" && account == nil {
		return false , "Account is not found", err
	}

	error := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if error != nil && error == bcrypt.ErrMismatchedHashAndPassword  {
		log.Fatal(error.Error())
		return false , " invalid password, or incorrect", ""
	}

	//works well , lets log u in
	password = ""

	//generate jwt for new session
	tk := token{UserID: name, UserRole:account.Role}

	signer := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)

	tokenString , errr := signer.SignedString([]byte("lekan"))
	if errr != nil {
		log.Fatal("error ==>  " +errr.Error())
		return false,  "error signing key",""
	}

	account.token = tokenString

	account.Password = ""
	return true, account.token, account.Name

}
func GetAllAccounts() []*Account {
	statement := `SELECT * FROM accounts`

	result , err := Database().Query(statement)
	if err != nil {
		fmt.Println("ann error occur getting the data , " + err.Error())
	}

	accounts := []*Account{}

	for result.Next() {
		var id int
		var name string
		var password string
		var token string
		var role string

		if error := result.Scan(&id, &name, &password, &token, &role); error != nil{
			fmt.Println("error while scanning the database reuslts , " + error.Error())
		}

		accounts = append(accounts, &Account{id, name, password, token, role})
	}

	return accounts
}