package main

import (
	"fmt"
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

const postURL = "http://finalistx.com/email.php"

func SendPost(data *CourseSign) error {
	
	params := fmt.Sprintf("name=%s&email=%s&tel=%s&course=%s", data.Name, data.Email, data.Telephone, data.Course)

	buf := bytes.NewBufferString(params)

	resp, err := http.Post(
		postURL,
		"application/x-www-forum-urlencoded",
		buf,
	)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	bd, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("ioutil: %s", bd)

	if resp.StatusCode != 200 {
		return errors.New("not 200 return")
	}

	return nil  
}