package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"rest/response"
	"strings"
)

type Address struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}
type User struct {
	Cpf     string  `json:"cpf"`
	Name    string  `json:"name"`
	Age     int64   `json:"age"`
	Cep     string  `json:"cep"`
	Address Address `json:"address"`
}

func GetUsers(Client *sql.DB) ([]User, error) {

	rows, err := Client.Query("SELECT * from Users")

	var users = []User{}
	for rows.Next() {
		var usr User
		if err := rows.Scan(&usr.Cpf, &usr.Name, &usr.Age, &usr.Cep); err != nil {
			return []User{}, err
		}
		{
			address := getAddressByCep(usr.Cep)
			usr.Address = address
			users = append(users, usr)
		}
	}
	if err := rows.Err(); err != nil {
		return []User{}, err
	}
	return users, err
}

func GetUserByCpf(Client *sql.DB, cpf string) (User, error) {
	var user User
	row := Client.QueryRow("SELECT * from Users where cpf=" + cpf)
	err := row.Scan(&user.Cpf, &user.Name, &user.Age, &user.Cep)
	address := getAddressByCep(user.Cep)
	user.Address = address
	return user, err

}

func getAddressByCep(cep string) Address {
	var address Address
	resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	erro := json.Unmarshal(body, &address)
	if erro != nil {
		panic(err)
	}

	return address
}
func validateCep(cep string) bool {
	resp, _ := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode != 200 || strings.Contains(string(body), "erro") {
		return false
	} else {
		return true
	}
}
func CreateUser(Client *sql.DB, user User) response.ResponseHttp {
	validCep := validateCep(user.Cep)
	if validCep != true {
		return response.NewResponseHttp(400, "Invalid CEP")
	}
	_, err := Client.Query(fmt.Sprintf("insert into Users VALUES ('%s','%s','%d','%s');", user.Cpf, user.Name, user.Age, user.Cep))

	if err != nil {
		return response.NewResponseHttp(500, "Internal server error")

	}
	return response.NewResponseHttp(200, "Created succesfully")

}
