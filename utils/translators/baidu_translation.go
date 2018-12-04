// Copyright 2018 The Pouch Robot Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package translators

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// NewBaiduTranslator returns a translator of Baidu API
func NewBaiduTranslator(options BaiduTranslatorOptions) Translator {
	translator := baiduTranslator{
		Appid: options.Appid,
		Key:   options.Key,
		From:  "auto",
		To:    "en",
	}
	if options.From != "" {
		translator.From = options.From
	}
	if options.To != "" {
		translator.To = options.To
	}
	return &translator
}

type baiduTranslator struct {
	Appid string
	Key   string
	From  string
	To    string
}

// BaiduTranslatorOptions is what you need to get a new baidu translator
type BaiduTranslatorOptions struct {
	Appid string
	Key   string
	From  string
	To    string
}

type baiduTranslateResultUnit struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

type baiduTranslateResult struct {
	ErrorCode    int                        `json:"error_code,omitempty"`
	ErrorMsg     string                     `json:"error_msg,omitempty"`
	TransResults []baiduTranslateResultUnit `json:"trans_result,omitempty"`
	From         string                     `json:"from"`
	To           string                     `json:"to"`
}

// translate a single line for there could be mixed language which will affect language detection
func (bt baiduTranslator) translateLine(text string) (string, error) {
	req, err := http.NewRequest("GET", "http://api.fanyi.baidu.com/api/trans/vip/translate", nil)
	if err != nil {
		return "", err
	}

	from := "auto"
	if bt.From != "" {
		from = bt.From
	}
	to := "en"
	if bt.To != "" {
		to = bt.To
	}

	client := &http.Client{}
	salt := strconv.Itoa(rand.Intn(100000))
	sign := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s%s%s", bt.Appid, text, salt, bt.Key))))
	q := req.URL.Query()
	q.Add("q", text)
	q.Add("from", from)
	q.Add("to", to)
	q.Add("salt", salt)
	q.Add("appid", bt.Appid)
	q.Add("sign", sign)

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result baiduTranslateResult
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if result.ErrorCode != 0 {
		return "", errors.New(result.ErrorMsg)
	}
	transResults := result.TransResults
	if transResults == nil {
		return "", errors.New("Unknown error")
	}
	// if detect the origin language is what we want, return empty string
	if result.From == bt.To {
		return "", nil
	}
	var ret string
	for _, res := range transResults {
		ret += res.Dst + "\r\n"
	}
	return ret, nil
}

// Translate translate what came in, return empty string if error occurred or no need to translate
func (bt baiduTranslator) Translate(text string) string {
	// get each line first
	re := regexp.MustCompile(`\r\n`)
	t := re.ReplaceAllString(text, "\n")
	lines := strings.Split(t, "\n")
	needTrans := false
	ret := ""
	for _, line := range lines {
		if line == "" {
			ret += line
			continue
		}
		trans, err := bt.translateLine(line)
		if err != nil {
			logrus.Error(err)
		} else if trans != "" {
			ret += line + "\r\n// " + trans
			needTrans = true
		} else {
			ret += line
		}
		ret += "\r\n"
	}
	if needTrans {
		return strings.TrimSpace(ret)
	}
	return ""
}
