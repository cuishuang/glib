package glib

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

/* ================================================================================
 * Http
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Http Get请求
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func HttpGet(url string, args ...string) (string, error) {
	requestUrl := url
	if len(args) == 1 {
		params := args[0]
		requestUrl = fmt.Sprintf("%s?%s", url, params)
	}

	resp, err := http.Get(requestUrl)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Http POST请求
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func HttpPost(url, params string, args ...string) (string, error) {
	cookie := ""
	if len(args) == 1 {
		cookie = args[0]
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(params))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if len(cookie) > 0 {
		req.Header.Set("Cookie", cookie)
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 上传文件
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func HttpPostFile(url, filename, fileTag string, params map[string]string) (int, map[string][]string, string, error) {
	if !filepath.IsAbs(filename) {
		filename, _ = filepath.Abs(filename)
	}
	file, err := os.Open(filename)
	if err != nil {
		return 0, nil, "", err
	}
	defer file.Close()

	if fileTag == "" {
		fileTag = "file"
	}

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	formFile, err := bodyWriter.CreateFormFile(fileTag, filepath.Base(filename))
	if err != nil {
		return 0, nil, "", err
	}

	//写入file数据
	_, err = io.Copy(formFile, file)
	if err != nil {
		return 0, nil, "", err
	}

	//写入参数
	for key, val := range params {
		_ = bodyWriter.WriteField(key, val)
	}

	contentType := bodyWriter.FormDataContentType()
	err = bodyWriter.Close()
	if err != nil {
		return 0, nil, "", err
	}

	//http post
	resp, err := http.Post(url, contentType, bodyBuffer)
	if err != nil {
		return 0, nil, "", err
	}
	defer resp.Body.Close()

	//状态码
	statusCode, _ := strconv.Atoi(resp.Status)

	header := resp.Header

	//获取响应数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, "", err
	}

	return statusCode, header, string(body), err
}
