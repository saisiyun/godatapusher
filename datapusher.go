// Copyright (c) 2018 - Saisiyun Co., Ltd.
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package godatapusher

import (
	"encoding/json"
	"errors"
	"time"
	"fmt"
	"bytes"
	"net/http"
)

type Config struct {
	AccessKey  string
	Event string
	Url string
}

var dataPushConfig *Config

// Initial config
func Init(config *Config) error{
	if config == nil {
		return errors.New("miss config")
	}
	if config.AccessKey == "" {
		return errors.New("miss accessKey")
	}
	dataPushConfig = config
	return nil
}

// Post data to DataPusher via http POST
func Post(ob interface{}) error {
	if ob == nil {
		return errors.New("miss object")
	}
	if dataPushConfig.Event == "" {
		return errors.New("miss event")
	}
	dat, err := json.Marshal(ob)
	if err != nil {
		return err
	}

	go postDate(dat, dataPushConfig.Event)
	return nil
}

// Post data to DataPusher via http POST
func PostWithEvent(ob interface{}, event string) error {
	if ob == nil {
		return errors.New("miss object")
	}
	if event == "" {
		return errors.New("miss event")
	}
	dat, err := json.Marshal(ob)
	if err != nil {
		return err
	}

	go postDate(dat, event)
	return nil
}

// Post data to DataPusher via http POST
func postDate(dat []byte, event string) error {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	host := "http://123.207.72.235:3000"
	if dataPushConfig.Url != "" {
		host = dataPushConfig.Url
	}
	invokeUrl := fmt.Sprintf("%s/v1/project/%s/events/%s", host, dataPushConfig.AccessKey, event)
	req, err := http.NewRequest("POST", invokeUrl, bytes.NewBuffer(dat))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

