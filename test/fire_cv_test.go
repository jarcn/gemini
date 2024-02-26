package test

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestDown(t *testing.T) {
	url := "https://open.larksuite.com/open-apis/drive/v1/medias/DO5QbzEIqoKalXxyTc2uhyLvsib/download?extra=%7B%22bitablePerm%22%3A%7B%22tableId%22%3A%22tblxPGXWxc3luvS9%22%2C%22rev%22%3A18146%7D%7D"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer u-fESpk1dol0GbcoLQD.5Ar8Bhmxdxl1NxoUw0h1Q003ke")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
