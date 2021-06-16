package plagiarismchecker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const textRuUrl = "http://api.text.ru/post"

type Checker struct {
	userKey string
	visible bool
}

type AddTextResponse struct {
	TextUid   string `json:"text_uid,omitempty"`
	ErrorCode int    `json:"error_code,omitempty"`
	ErrorDesc string `json:"error_desc,omitempty"`
}

type GetResultResponse struct {
	TextUnique string `json:"text_unique,omitempty"`
	ErrorCode  int    `json:"error_code,omitempty"`
	ErrorDesc  string `json:"error_desc,omitempty"`
}

func New(userKey string, visible bool) *Checker {
	return &Checker{
		userKey: userKey,
		visible: visible,
	}
}

func (c *Checker) AddText(text string) (string, error) {
	data := url.Values{}
	data.Add("text", text)
	data.Add("userkey", c.userKey)
	if c.visible {
		data.Add("visible", "vis_on")
	}

	response, err := http.PostForm(textRuUrl, data)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var addTextResponse AddTextResponse
	if err := decoder.Decode(&addTextResponse); err != nil {
		return "", err
	}

	if addTextResponse.ErrorCode != 0 {
		return "", fmt.Errorf("%d: %s", addTextResponse.ErrorCode, addTextResponse.ErrorDesc)
	}

	return addTextResponse.TextUid, nil
}

func (c *Checker) GetResult(uid string) (float32, error) {
	data := url.Values{}
	data.Add("uid", uid)
	data.Add("userkey", c.userKey)

	for {
		response, err := http.PostForm(textRuUrl, data)
		if err != nil {
			return 0, err
		}

		defer response.Body.Close()

		decoder := json.NewDecoder(response.Body)
		var getResultResponse GetResultResponse
		if err := decoder.Decode(&getResultResponse); err != nil {
			return 0, err
		}

		if getResultResponse.ErrorCode == 181 {
			time.Sleep(time.Second)
			continue
		}

		if getResultResponse.ErrorCode != 0 {
			return 0, fmt.Errorf("%d: %s", getResultResponse.ErrorCode, getResultResponse.ErrorDesc)
		}

		unique, err := strconv.ParseFloat(getResultResponse.TextUnique, 32)
		if err != nil {
			return 0, err
		}
		return float32(unique), nil
	}
}
