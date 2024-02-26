package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	endpoint := "http://localhost:8080/"
	fmt.Println("Введите длинный URL")
	reader := bufio.NewReader(os.Stdin)
	long, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	long = strings.TrimSuffix(long, "\n")
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(long))
	if err != nil {
		panic(err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	fmt.Println("Статус-код ", response.Status)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}
