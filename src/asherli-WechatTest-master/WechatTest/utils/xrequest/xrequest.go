/*
 *    Copyright 2016-2018 Li ZongZe
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package xrequest

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/moul/http2curl"
	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
)

type Request *http.Request
type Response *http.Response

const (
	POST    = "POST"
	GET     = "GET"
	HEAD    = "HEAD"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
)

type basicAuth struct {
	Username string
	Password string
}

type retryAble struct {
	RetryAbleStatus []int
	RetryTime       time.Duration
	RetryCount      int
	Attempt         int
	Enable          bool
}

type xRequest struct {
	Url               string
	Header            map[string]string
	Method            string
	TargetType        string
	ForceType         string
	Data              map[string]interface{}
	SliceData         []interface{}
	FormData          url.Values
	QueryData         url.Values
	FileData          []File
	BounceToRawString bool
	RawString         string
	Client            *http.Client
	Transport         *http.Transport
	Cookies           []*http.Cookie
	Errors            []error
	BasicAuth         basicAuth
	Debug             bool
	CurlCommand       bool
	logger            Logger
	RetryAble         retryAble
}

var DisableTransportSwap = false

func New() *xRequest {
	cookieJarOptions := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, _ := cookiejar.New(&cookieJarOptions)

	debug := os.Getenv("X_REQUEST_DEBUG") == "1"

	s := &xRequest{
		TargetType:        "JSON",
		Data:              make(map[string]interface{}),
		Header:            make(map[string]string),
		RawString:         "",
		SliceData:         []interface{}{},
		FormData:          url.Values{},
		QueryData:         url.Values{},
		FileData:          make([]File, 0),
		BounceToRawString: false,
		Client: &http.Client{
			Jar: jar,
		},
		Transport:   &http.Transport{},
		Cookies:     make([]*http.Cookie, 0),
		Errors:      nil,
		BasicAuth:   basicAuth{},
		Debug:       debug,
		CurlCommand: false,
		logger:      log.New(os.Stderr, "[XRequest]", log.LstdFlags),
	}
	s.Transport.DisableKeepAlives = true
	return s
}

func (s *xRequest) IsDebug(enable bool) *xRequest {
	s.Debug = enable
	return s
}

func (s *xRequest) IsCurlCommand(enable bool) *xRequest {
	s.CurlCommand = enable
	return s
}

func (s *xRequest) SetLogger(logger Logger) *xRequest {
	s.logger = logger
	return s
}

func (s *xRequest) Custom(method, targetUrl string) *xRequest {
	switch method {
	case POST:
		return s.Post(targetUrl)
	case GET:
		return s.Get(targetUrl)
	case HEAD:
		return s.Head(targetUrl)
	case PUT:
		return s.Put(targetUrl)
	case DELETE:
		return s.Delete(targetUrl)
	case PATCH:
		return s.Patch(targetUrl)
	case OPTIONS:
		return s.Options(targetUrl)
	default:
		s.Method = method
		s.Url = targetUrl
		s.Errors = nil
		return s
	}
}

func (s *xRequest) Get(targetUrl string) *xRequest {
	s.Method = GET
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *xRequest) Post(targetUrl string) *xRequest {
	s.Method = POST
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *xRequest) Head(targetUrl string) *xRequest {
	s.Method = HEAD
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *xRequest) Put(targetUrl string) *xRequest {
	s.Method = PUT
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *xRequest) Delete(targetUrl string) *xRequest {
	s.Method = DELETE
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *xRequest) Patch(targetUrl string) *xRequest {
	s.Method = PATCH
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *xRequest) Options(targetUrl string) *xRequest {
	s.Method = OPTIONS
	s.Url = targetUrl
	s.Errors = nil
	return s
}

func (s *xRequest) AddHeader(param string, value string) *xRequest {
	s.Header[param] = value
	return s
}

func (s *xRequest) Retry(retryerCount int, retryerTime time.Duration, statusCode ...int) *xRequest {
	for _, code := range statusCode {
		statusText := http.StatusText(code)
		if len(statusText) == 0 {
			s.Errors = append(s.Errors, errors.New("状态码： '"+strconv.Itoa(code)+"' 不存在"))
		}
	}

	s.RetryAble = struct {
		RetryAbleStatus []int
		RetryTime       time.Duration
		RetryCount      int
		Attempt         int
		Enable          bool
	}{
		RetryAbleStatus: statusCode,
		RetryTime:       retryerTime,
		RetryCount:      retryerCount,
		Attempt:         0,
		Enable:          true,
	}
	return s
}

func (s *xRequest) SetBasicAuth(username string, password string) *xRequest {
	s.BasicAuth = basicAuth{
		Username: username,
		Password: password,
	}
	return s
}

func (s *xRequest) AddCookie(c *http.Cookie) *xRequest {
	s.Cookies = append(s.Cookies, c)
	return s
}

func (s *xRequest) AddCookies(cookies []*http.Cookie) *xRequest {
	s.Cookies = append(s.Cookies, cookies...)
	return s
}

var Types = map[string]string{
	"HTML":       "text/html",
	"JSON":       "application/json",
	"XML":        "application/xml",
	"TEXT":       "text/plain",
	"URLENCODED": "application/x-www-form-urlencoded",
	"FORM":       "application/x-www-form-urlencoded",
	"FORM-DATA":  "application/x-www-form-urlencoded",
	"MULTIPART":  "multipart/form-data",
}

func (s *xRequest) SetForceType(typeStr string) *xRequest {
	if _, ok := Types[strings.ToUpper(typeStr)]; ok {
		s.ForceType = typeStr
	} else {
		s.Errors = append(s.Errors, errors.New("类型参数： 不存在类型 \""+typeStr+"\""))
	}
	return s
}

func (s *xRequest) QueryString(content interface{}) *xRequest {
	switch v := reflect.ValueOf(content); v.Kind() {
	case reflect.String:
		s.queryString(v.String())
	case reflect.Struct:
		s.queryStruct(v.Interface())
	case reflect.Map:
		s.queryMap(v.Interface())
	default:
	}
	return s
}

func (s *xRequest) queryStruct(queryString interface{}) *xRequest {
	if marshalContent, err := json.Marshal(queryString); err != nil {
		s.Errors = append(s.Errors, err)
	} else {
		var val map[string]interface{}
		if err := json.Unmarshal(marshalContent, &val); err != nil {
			s.Errors = append(s.Errors, err)
		} else {
			for k, v := range val {
				k = strings.ToLower(k)
				var queryVal string
				switch t := v.(type) {
				case string:
					queryVal = t
				case float64:
					queryVal = strconv.FormatFloat(t, 'f', -1, 64)
				case time.Time:
					queryVal = t.Format(time.RFC3339)
				default:
					j, err := json.Marshal(v)
					if err != nil {
						continue
					}
					queryVal = string(j)
				}
				s.QueryData.Add(k, queryVal)
			}
		}
	}
	return s
}

func (s *xRequest) queryString(queryString string) *xRequest {
	var val map[string]string
	if err := json.Unmarshal([]byte(queryString), &val); err == nil {
		for k, v := range val {
			s.QueryData.Add(k, v)
		}
	} else {
		if queryData, err := url.ParseQuery(queryString); err == nil {
			for k, queryValues := range queryData {
				for _, queryValue := range queryValues {
					s.QueryData.Add(k, string(queryValue))
				}
			}
		} else {
			s.Errors = append(s.Errors, err)
		}
	}
	return s
}

func (s *xRequest) queryMap(queryString interface{}) *xRequest {
	return s.queryStruct(queryString)
}

func (s *xRequest) RequestParams(key string, value string) *xRequest {
	s.QueryData.Add(key, value)
	return s
}

func (s *xRequest) ConnectTimeout(timeout time.Duration) *xRequest {
	s.Transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, timeout)
		if err != nil {
			s.Errors = append(s.Errors, err)
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(timeout))
		return conn, nil
	}
	return s
}

func (s *xRequest) TLSClientConfig(config *tls.Config) *xRequest {
	s.Transport.TLSClientConfig = config
	return s
}

func (s *xRequest) SetProxy(proxyUrl string) *xRequest {
	parsedProxyUrl, err := url.Parse(proxyUrl)
	if err != nil {
		s.Errors = append(s.Errors, err)
	} else if proxyUrl == "" {
		s.Transport.Proxy = nil
	} else {
		s.Transport.Proxy = http.ProxyURL(parsedProxyUrl)
	}
	return s
}

func (s *xRequest) RedirectPolicy(policy func(req Request, via []Request) error) *xRequest {
	s.Client.CheckRedirect = func(r *http.Request, v []*http.Request) error {
		vv := make([]Request, len(v))
		for i, r := range v {
			vv[i] = Request(r)
		}
		return policy(Request(r), vv)
	}
	return s
}

func (s *xRequest) Send(content interface{}) *xRequest {
	switch v := reflect.ValueOf(content); v.Kind() {
	case reflect.String:
		s.SendString(v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: // includes rune
		s.SendString(strconv.FormatInt(v.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: // includes byte
		s.SendString(strconv.FormatUint(v.Uint(), 10))
	case reflect.Float64:
		s.SendString(strconv.FormatFloat(v.Float(), 'f', -1, 64))
	case reflect.Float32:
		s.SendString(strconv.FormatFloat(v.Float(), 'f', -1, 32))
	case reflect.Bool:
		s.SendString(strconv.FormatBool(v.Bool()))
	case reflect.Struct:
		s.SendStruct(v.Interface())
	case reflect.Slice:
		s.SendSlice(makeSliceOfReflectValue(v))
	case reflect.Array:
		s.SendSlice(makeSliceOfReflectValue(v))
	case reflect.Ptr:
		s.Send(v.Elem().Interface())
	case reflect.Map:
		s.SendMap(v.Interface())
	default:
		return s
	}
	return s
}

func makeSliceOfReflectValue(v reflect.Value) (slice []interface{}) {

	kind := v.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return slice
	}

	slice = make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		slice[i] = v.Index(i).Interface()
	}

	return slice
}

func (s *xRequest) SendSlice(content []interface{}) *xRequest {
	s.SliceData = append(s.SliceData, content...)
	return s
}

func (s *xRequest) SendMap(content interface{}) *xRequest {
	return s.SendStruct(content)
}

func (s *xRequest) SendStruct(content interface{}) *xRequest {
	if marshalContent, err := json.Marshal(content); err != nil {
		s.Errors = append(s.Errors, err)
	} else {
		var val map[string]interface{}
		d := json.NewDecoder(bytes.NewBuffer(marshalContent))
		d.UseNumber()
		if err := d.Decode(&val); err != nil {
			s.Errors = append(s.Errors, err)
		} else {
			for k, v := range val {
				s.Data[k] = v
			}
		}
	}
	return s
}

func (s *xRequest) SendString(content string) *xRequest {
	if !s.BounceToRawString {
		var val interface{}
		d := json.NewDecoder(strings.NewReader(content))
		d.UseNumber()
		if err := d.Decode(&val); err == nil {
			switch v := reflect.ValueOf(val); v.Kind() {
			case reflect.Map:
				for k, v := range val.(map[string]interface{}) {
					s.Data[k] = v
				}
				// add to SliceData
			case reflect.Slice:
				s.SendSlice(val.([]interface{}))
				// bounce to rawstring if it is arrayjson, or others
			default:
				s.BounceToRawString = true
			}
		} else if formData, err := url.ParseQuery(content); err == nil {
			for k, formValues := range formData {
				for _, formValue := range formValues {
					// make it array if already have key
					if val, ok := s.Data[k]; ok {
						var strArray []string
						strArray = append(strArray, string(formValue))
						// check if previous data is one string or array
						switch oldValue := val.(type) {
						case []string:
							strArray = append(strArray, oldValue...)
						case string:
							strArray = append(strArray, oldValue)
						}
						s.Data[k] = strArray
					} else {
						// make it just string if does not already have same key
						s.Data[k] = formValue
					}
				}
			}
			s.TargetType = "form"
		} else {
			s.BounceToRawString = true
		}
	}
	// Dump all contents to RawString in case in the end user doesn't want json or form.
	s.RawString += content
	return s
}

type File struct {
	Filename  string
	Fieldname string
	Data      []byte
}

// SendFile function works only with type "multipart". The function accepts one mandatory and up to two optional arguments. The mandatory (first) argument is the file.
// The function accepts a path to a file as string:
//
//      gorequest.New().
//        Post("http://example.com").
//        Type("multipart").
//        SendFile("./example_file.ext").
//        End()
//
// File can also be a []byte slice of a already file read by eg. ioutil.ReadFile:
//
//      b, _ := ioutil.ReadFile("./example_file.ext")
//      gorequest.New().
//        Post("http://example.com").
//        Type("multipart").
//        SendFile(b).
//        End()
//
// Furthermore file can also be a os.File:
//
//      f, _ := os.Open("./example_file.ext")
//      gorequest.New().
//        Post("http://example.com").
//        Type("multipart").
//        SendFile(f).
//        End()
//
// The first optional argument (second argument overall) is the filename, which will be automatically determined when file is a string (path) or a os.File.
// When file is a []byte slice, filename defaults to "filename". In all cases the automatically determined filename can be overwritten:
//
//      b, _ := ioutil.ReadFile("./example_file.ext")
//      gorequest.New().
//        Post("http://example.com").
//        Type("multipart").
//        SendFile(b, "my_custom_filename").
//        End()
//
// The second optional argument (third argument overall) is the fieldname in the multipart/form-data request. It defaults to fileNUMBER (eg. file1), where number is ascending and starts counting at 1.
// So if you send multiple files, the fieldnames will be file1, file2, ... unless it is overwritten. If fieldname is set to "file" it will be automatically set to fileNUMBER, where number is the greatest exsiting number+1.
//
//      b, _ := ioutil.ReadFile("./example_file.ext")
//      gorequest.New().
//        Post("http://example.com").
//        Type("multipart").
//        SendFile(b, "", "my_custom_fieldname"). // filename left blank, will become "example_file.ext"
//        End()
//
func (s *xRequest) SendFile(file interface{}, args ...string) *xRequest {

	filename := ""
	fieldname := "file"

	if len(args) >= 1 && len(args[0]) > 0 {
		filename = strings.TrimSpace(args[0])
	}
	if len(args) >= 2 && len(args[1]) > 0 {
		fieldname = strings.TrimSpace(args[1])
	}
	if fieldname == "file" || fieldname == "" {
		fieldname = "file" + strconv.Itoa(len(s.FileData)+1)
	}

	switch v := reflect.ValueOf(file); v.Kind() {
	case reflect.String:
		pathToFile, err := filepath.Abs(v.String())
		if err != nil {
			s.Errors = append(s.Errors, err)
			return s
		}
		if filename == "" {
			filename = filepath.Base(pathToFile)
		}
		data, err := ioutil.ReadFile(v.String())
		if err != nil {
			s.Errors = append(s.Errors, err)
			return s
		}
		s.FileData = append(s.FileData, File{
			Filename:  filename,
			Fieldname: fieldname,
			Data:      data,
		})
	case reflect.Slice:
		slice := makeSliceOfReflectValue(v)
		if filename == "" {
			filename = "filename"
		}
		f := File{
			Filename:  filename,
			Fieldname: fieldname,
			Data:      make([]byte, len(slice)),
		}
		for i := range slice {
			f.Data[i] = slice[i].(byte)
		}
		s.FileData = append(s.FileData, f)
	case reflect.Ptr:
		if len(args) == 1 {
			return s.SendFile(v.Elem().Interface(), args[0])
		}
		if len(args) >= 2 {
			return s.SendFile(v.Elem().Interface(), args[0], args[1])
		}
		return s.SendFile(v.Elem().Interface())
	default:
		if v.Type() == reflect.TypeOf(os.File{}) {
			osfile := v.Interface().(os.File)
			if filename == "" {
				filename = filepath.Base(osfile.Name())
			}
			data, err := ioutil.ReadFile(osfile.Name())
			if err != nil {
				s.Errors = append(s.Errors, err)
				return s
			}
			s.FileData = append(s.FileData, File{
				Filename:  filename,
				Fieldname: fieldname,
				Data:      data,
			})
			return s
		}

		s.Errors = append(s.Errors, errors.New("SendFile currently only supports either a string (path/to/file), a slice of bytes (file content itself), or a os.File!"))
	}

	return s
}

func changeMapToURLValues(data map[string]interface{}) url.Values {
	var newUrlValues = url.Values{}
	for k, v := range data {
		switch val := v.(type) {
		case string:
			newUrlValues.Add(k, val)
		case bool:
			newUrlValues.Add(k, strconv.FormatBool(val))
		case json.Number:
			newUrlValues.Add(k, string(val))
		case int:
			newUrlValues.Add(k, strconv.FormatInt(int64(val), 10))
		case float64:
			newUrlValues.Add(k, strconv.FormatFloat(float64(val), 'f', -1, 64))
		case float32:
			newUrlValues.Add(k, strconv.FormatFloat(float64(val), 'f', -1, 64))
		case []string:
			for _, element := range val {
				newUrlValues.Add(k, element)
			}
		case []int:
			for _, element := range val {
				newUrlValues.Add(k, strconv.FormatInt(int64(element), 10))
			}
		case []bool:
			for _, element := range val {
				newUrlValues.Add(k, strconv.FormatBool(element))
			}
		case []float64:
			for _, element := range val {
				newUrlValues.Add(k, strconv.FormatFloat(float64(element), 'f', -1, 64))
			}
		case []float32:
			for _, element := range val {
				newUrlValues.Add(k, strconv.FormatFloat(float64(element), 'f', -1, 64))
			}
		case []interface{}:

			if len(val) <= 0 {
				continue
			}

			switch val[0].(type) {
			case string:
				for _, element := range val {
					newUrlValues.Add(k, element.(string))
				}
			case bool:
				for _, element := range val {
					newUrlValues.Add(k, strconv.FormatBool(element.(bool)))
				}
			case json.Number:
				for _, element := range val {
					newUrlValues.Add(k, string(element.(json.Number)))
				}
			}
		default:
			// TODO add ptr, arrays, ...
		}
	}
	return newUrlValues
}

// End is the most important function that you need to call when ending the chain. The request won't proceed without calling it.
// End function returns Response which matchs the structure of Response type in Golang's http package (but without Body data). The body data itself returns as a string in a 2nd return value.
// Lastly but worth noticing, error array (NOTE: not just single error value) is returned as a 3rd value and nil otherwise.
//
// For example:
//
//    resp, body, errs := gorequest.New().Get("http://www.google.com").End()
//    if (errs != nil) {
//      fmt.Println(errs)
//    }
//    fmt.Println(resp, body)
//
// Moreover, End function also supports callback which you can put as a parameter.
// This extends the flexibility and makes GoRequest fun and clean! You can use GoRequest in whatever style you love!
//
// For example:
//
//    func printBody(resp gorequest.Response, body string, errs []error){
//      fmt.Println(resp.Status)
//    }
//    gorequest.New().Get("http://www..google.com").End(printBody)
//
func (s *xRequest) End(callback ...func(response Response, body string, errs []error)) (Response, string, []error) {
	var bytesCallback []func(response Response, body []byte, errs []error)
	if len(callback) > 0 {
		bytesCallback = []func(response Response, body []byte, errs []error){
			func(response Response, body []byte, errs []error) {
				callback[0](response, string(body), errs)
			},
		}
	}

	resp, body, errs := s.EndBytes(bytesCallback...)
	bodyString := string(body)

	return resp, bodyString, errs
}

// EndBytes should be used when you want the body as bytes. The callbacks work the same way as with `End`, except that a byte array is used instead of a string.
func (s *xRequest) EndBytes(callback ...func(response Response, body []byte, errs []error)) (Response, []byte, []error) {
	var (
		errs []error
		resp Response
		body []byte
	)

	for {
		resp, body, errs = s.getResponseBytes()
		if errs != nil {
			return nil, nil, errs
		}
		if s.isRetryableRequest(resp) {
			resp.Header.Set("Retry-Count", strconv.Itoa(s.RetryAble.Attempt))
			break
		}
	}

	respCallback := *resp
	if len(callback) != 0 {
		callback[0](&respCallback, body, s.Errors)
	}
	return resp, body, nil
}

func (s *xRequest) isRetryableRequest(resp Response) bool {
	if s.RetryAble.Enable && s.RetryAble.Attempt < s.RetryAble.RetryCount && contains(resp.StatusCode, s.RetryAble.RetryAbleStatus) {
		time.Sleep(s.RetryAble.RetryTime)
		s.RetryAble.Attempt++
		return false
	}
	return true
}

func contains(respStatus int, statuses []int) bool {
	for _, status := range statuses {
		if status == respStatus {
			return true
		}
	}
	return false
}

// EndStruct should be used when you want the body as a struct. The callbacks work the same way as with `End`, except that a struct is used instead of a string.
func (s *xRequest) EndStruct(v interface{}, callback ...func(response Response, v interface{}, body []byte, errs []error)) (Response, []byte, []error) {
	resp, body, errs := s.EndBytes()
	if errs != nil {
		return nil, body, errs
	}
	fmt.Println(string(body))
	err := json.Unmarshal(body, &v)
	if err != nil {
		s.Errors = append(s.Errors, err)
		return resp, body, s.Errors
	}
	respCallback := *resp
	if len(callback) != 0 {
		callback[0](&respCallback, v, body, s.Errors)
	}
	return resp, body, nil
}

func (s *xRequest) getResponseBytes() (Response, []byte, []error) {
	var (
		req  *http.Request
		err  error
		resp Response
	)
	// check whether there is an error. if yes, return all errors
	if len(s.Errors) != 0 {
		return nil, nil, s.Errors
	}
	// check if there is forced type
	switch strings.ToUpper(s.ForceType) {
	case "JSON", "FORM", "XML", "TEXT", "MULTIPART":
		s.TargetType = s.ForceType
		// If forcetype is not set, check whether user set Content-Type header.
		// If yes, also bounce to the correct supported TargetType automatically.
	default:
		for k, v := range Types {
			if s.Header["Content-Type"] == v {
				s.TargetType = k
			}
		}
	}

	// if slice and map get mixed, let's bounce to rawstring
	if len(s.Data) != 0 && len(s.SliceData) != 0 {
		s.BounceToRawString = true
	}

	// Make Request
	req, err = s.MakeRequest()
	if err != nil {
		s.Errors = append(s.Errors, err)
		return nil, nil, s.Errors
	}

	// Set Transport
	if !DisableTransportSwap {
		s.Client.Transport = s.Transport
	}

	// Log details of this request
	if s.Debug {
		dump, err := httputil.DumpRequest(req, true)
		s.logger.SetPrefix("[http] ")
		if err != nil {
			s.logger.Println("Error:", err)
		} else {
			s.logger.Printf("HTTP Request: %s", string(dump))
		}
	}

	// Display CURL command line
	if s.CurlCommand {
		curl, err := http2curl.GetCurlCommand(req)
		s.logger.SetPrefix("[curl] ")
		if err != nil {
			s.logger.Println("Error:", err)
		} else {
			s.logger.Printf("CURL command line: %s", curl)
		}
	}

	// Send request
	resp, err = s.Client.Do(req)
	if err != nil {
		s.Errors = append(s.Errors, err)
		return nil, nil, s.Errors
	}
	defer resp.Body.Close()

	// Log details of this response
	if s.Debug {
		dump, err := httputil.DumpResponse(resp, true)
		if nil != err {
			s.logger.Println("Error:", err)
		} else {
			s.logger.Printf("HTTP Response: %s", string(dump))
		}
	}

	body, _ := ioutil.ReadAll(resp.Body)
	// Reset resp.Body so it can be use again
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return resp, body, nil
}

func (s *xRequest) MakeRequest() (*http.Request, error) {
	var (
		req           *http.Request
		contentType   string
		contentReader io.Reader
		err           error
	)

	if s.Method == "" {
		return nil, errors.New("No method specified")
	}
	switch strings.ToUpper(s.TargetType) {
	case "JSON":
		var contentJson []byte
		if s.BounceToRawString {
			contentJson = []byte(s.RawString)
		} else if len(s.Data) != 0 {
			contentJson, _ = json.Marshal(s.Data)
		} else if len(s.SliceData) != 0 {
			contentJson, _ = json.Marshal(s.SliceData)
		}
		if contentJson != nil {
			contentReader = bytes.NewReader(contentJson)
			contentType = "application/json"
		}
	case "FORM", "FROM-DATA", "URLENCODED":
		var contentForm []byte
		if s.BounceToRawString || len(s.SliceData) != 0 {
			contentForm = []byte(s.RawString)
		} else {
			formData := changeMapToURLValues(s.Data)
			contentForm = []byte(formData.Encode())
		}
		if len(contentForm) != 0 {
			contentReader = bytes.NewReader(contentForm)
			contentType = "application/x-www-form-urlencoded"
		}
	case "TEXT":
		if len(s.RawString) != 0 {
			contentReader = strings.NewReader(s.RawString)
			contentType = "text/plain"
		}
	case "XML":
		if len(s.RawString) != 0 {
			contentReader = strings.NewReader(s.RawString)
			contentType = "application/xml"
		}
	case "MULTIPART":
		var (
			buf = &bytes.Buffer{}
			mw  = multipart.NewWriter(buf)
		)

		if s.BounceToRawString {
			fieldName, ok := s.Header["data_fieldname"]
			if !ok {
				fieldName = "data"
			}
			fw, _ := mw.CreateFormField(fieldName)
			fw.Write([]byte(s.RawString))
			contentReader = buf
		}

		if len(s.Data) != 0 {
			formData := changeMapToURLValues(s.Data)
			for key, values := range formData {
				for _, value := range values {
					fw, _ := mw.CreateFormField(key)
					fw.Write([]byte(value))
				}
			}
			contentReader = buf
		}

		if len(s.SliceData) != 0 {
			fieldName, ok := s.Header["json_fieldname"]
			if !ok {
				fieldName = "data"
			}
			// copied from CreateFormField() in mime/multipart/writer.go
			h := make(textproto.MIMEHeader)
			fieldName = strings.Replace(strings.Replace(fieldName, "\\", "\\\\", -1), `"`, "\\\"", -1)
			h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`, fieldName))
			h.Set("Content-Type", "application/json")
			fw, _ := mw.CreatePart(h)
			contentJson, err := json.Marshal(s.SliceData)
			if err != nil {
				return nil, err
			}
			fw.Write(contentJson)
			contentReader = buf
		}

		// add the files
		if len(s.FileData) != 0 {
			for _, file := range s.FileData {
				fw, _ := mw.CreateFormFile(file.Fieldname, file.Filename)
				fw.Write(file.Data)
			}
			contentReader = buf
		}

		// close before call to FormDataContentType ! otherwise its not valid multipart
		mw.Close()

		if contentReader != nil {
			contentType = mw.FormDataContentType()
		}
	default:
		// let's return an error instead of an nil pointer exception here
		return nil, errors.New("TargetType '" + s.TargetType + "' could not be determined")
	}

	if req, err = http.NewRequest(s.Method, s.Url, contentReader); err != nil {
		return nil, err
	}

	if len(contentType) != 0 {
		req.Header.Set("Content-Type", contentType)
	}

	for k, v := range s.Header {
		req.Header.Set(k, v)
		// Setting the host header is a special case, see this issue: https://github.com/golang/go/issues/7682
		if strings.EqualFold(k, "host") {
			req.Host = v
		}
	}
	// Add all querystring from Query func
	q := req.URL.Query()
	for k, v := range s.QueryData {
		for _, vv := range v {
			q.Add(k, vv)
		}
	}
	req.URL.RawQuery = q.Encode()

	// Add basic auth
	if s.BasicAuth != struct{ Username, Password string }{} {
		req.SetBasicAuth(s.BasicAuth.Username, s.BasicAuth.Password)
	}

	// Add cookies
	for _, cookie := range s.Cookies {
		req.AddCookie(cookie)
	}

	return req, nil
}

// AsCurlCommand returns a string representing the runnable `curl' command
// version of the request.
func (s *xRequest) AsCurlCommand() (string, error) {
	req, err := s.MakeRequest()
	if err != nil {
		return "", err
	}
	cmd, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return "", err
	}
	return cmd.String(), nil
}
