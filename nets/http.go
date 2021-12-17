// Package nets
/*
 * Version: 1.0.0
 * Copyright (c) 2021. Pashifika
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package nets

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/pashifika/util/files"
)

var (
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client  = &http.Client{Transport: tr}
	bigSize = int64(1024 * 1024 * 10)

	// HTTP headers
	//acceptRangeHeader   = "Accept-Ranges"
	contentLengthHeader = "Content-Length"
)

func IsUrl(URL string) (*url.URL, error) {
	u, err := url.ParseRequestURI(URL)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, errors.New("invalid URI for request, url:" + URL)
	}
	return u, nil
}

// HttpDownload is auto join the urlPaths to URL parameter
//goland:noinspection GoUnusedExportedFunction
func HttpDownload(URL, localPath string, urlPaths ...string) error {
	u, err := IsUrl(URL)
	if err != nil {
		return err
	}
	if len(urlPaths) != 0 {
		u.Path = path.Join(append([]string{u.Path}, urlPaths...)...)
	}
	if err = files.MkdirIfNotExist(localPath); err != nil {
		return err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	//noinspection ALL
	defer resp.Body.Close()

	// get download size
	clen := resp.Header.Get(contentLengthHeader)
	if clen == "" {
		clen = "1"
	}

	size, err := strconv.ParseInt(clen, 10, 64)
	if err != nil {
		return err
	}

	if size >= bigSize {
		err = files.BufferToFile(localPath, resp.Body)
	} else {
		var buf []byte
		buf, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = files.ByteToFile(localPath, buf)
	}

	return err
}
