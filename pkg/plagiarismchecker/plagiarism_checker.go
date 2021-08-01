package plagiarismchecker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const textRuUrl = "http://api.text.ru/post"

type Checker struct {
	userKey       string
	visible       bool
	exceptDomains []string
}

type Uid = string

const (
	stillInProcessErrorCode = 181
	maxLoop                 = 100
)

var stillInProcessError = errors.New("Processing")

type AddTextResponse struct {
	TextUid   Uid    `json:"text_uid,omitempty"`
	ErrorCode int    `json:"error_code,omitempty"`
	ErrorDesc string `json:"error_desc,omitempty"`
}

type GetResultResponse struct {
	TextUnique string `json:"text_unique,omitempty"`
	ErrorCode  int    `json:"error_code,omitempty"`
	ErrorDesc  string `json:"error_desc,omitempty"`
}

func New(userKey string, visible bool, exceptDomains []string) *Checker {
	return &Checker{
		userKey:       userKey,
		visible:       visible,
		exceptDomains: exceptDomains,
	}
}

func (c *Checker) AddText(text string) (Uid, error) {
	data := url.Values{}
	data.Add("text", text)
	data.Add("userkey", c.userKey)
	if c.visible {
		data.Add("visible", "vis_on")
	}

	if len(c.exceptDomains) > 0 {
		data.Add("exceptdomain", strings.Join(c.exceptDomains, ","))
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

func (c *Checker) GetResult(uid Uid) (float32, error) {
	form := url.Values{}
	form.Add("uid", uid)
	form.Add("userkey", c.userKey)

	for i := 0; i < maxLoop; i++ {
		uniq, err := c.getResult(form)
		if err == stillInProcessError {
			time.Sleep(time.Second)
			continue
		}

		return uniq, err
	}

	return 0, errors.New("Timeout exceeded")
}

func (c *Checker) getResult(form url.Values) (float32, error) {
	response, err := http.PostForm(textRuUrl, form)
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	var getResultResponse GetResultResponse
	if err := decoder.Decode(&getResultResponse); err != nil {
		return 0, err
	}

	if getResultResponse.ErrorCode == stillInProcessErrorCode {
		return 0, stillInProcessError
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
