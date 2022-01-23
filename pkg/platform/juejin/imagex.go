package juejin

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/juju/errors"
	"github.com/tidwall/gjson"
)

const (
	amzDateISO8601TimeFormat = "20060102T150405Z"
	shortTimeFormat          = "20060102"
	algorithm                = "AWS4-HMAC-SHA256"
	serviceName              = "imagex"
	serviceID                = "k3u1fbpfcp"
	version                  = "2018-08-01"
	uploadURLFormat          = "https://%s/%s"

	RegionCNNorth = "cn-north-1"

	actionApplyImageUpload  = "ApplyImageUpload"
	actionCommitImageUpload = "CommitImageUpload"

	polynomialCRC32 = 0xEDB88320
)

var (
	newLine = []byte{'\n'}
)

type ImageX struct {
	AccessKey string
	SecretKey string
	Region    string
	Client    *http.Client

	Token   string
	Version string
	BaseURL string
}

type UploadToken struct {
	AccessKeyID     string `json:"AccessKeyID"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"SessionToken"`
}

func (c *Client) UploadImage(region, path string) (string, error) {
	uploadToken, err := c.GetUploadToken()
	if err != nil {
		return "", errors.Trace(err)
	}

	ix := &ImageX{
		AccessKey: uploadToken.AccessKeyID,
		SecretKey: uploadToken.SecretAccessKey,
		Token:     uploadToken.SessionToken,
		Region:    region,
	}

	applyRes, err := ix.ApplyImageUpload()
	if err != nil {
		return "", errors.Trace(err)
	}

	storeInfo := gjson.Get(applyRes, "Result.UploadAddress.StoreInfos.0")
	storeURI := storeInfo.Get("StoreUri").String()
	storeAuth := storeInfo.Get("Auth").String()
	uploadHost := gjson.Get(applyRes, "Result.UploadAddress.UploadHosts.0").String()
	uploadURL := fmt.Sprintf(uploadURLFormat, uploadHost, storeURI)
	if err := ix.Upload(uploadURL, path, storeAuth); err != nil {
		return "", errors.Trace(err)
	}

	sessionKey := gjson.Get(applyRes, "Result.UploadAddress.SessionKey").String()
	if _, err = ix.CommitImageUpload(sessionKey); err != nil {
		return "", errors.Trace(err)
	}

	imageURL, err := c.GetImageURL(storeURI)
	if err != nil {
		return "", errors.Trace(err)
	}
	return imageURL.MainURL, nil
}

type ImageURL struct {
	BackupURL string `json:"backup_url"`
	MainURL   string `json:"main_url"`
}

func (c *Client) GetImageURL(uri string) (*ImageURL, error) {
	endpoint := "/imagex/get_img_url"
	params := &url.Values{
		"uri": []string{uri},
	}
	raw, err := c.Get(endpoint, params)
	if err != nil {
		return nil, errors.Trace(err)
	}
	data := gjson.Get(raw, "data").String()
	if data == "" {
		return nil, errors.Errorf("invalid response: %s", raw)
	}
	var result *ImageURL
	err = json.Unmarshal([]byte(data), &result)
	return result, errors.Trace(err)
}

func (c *Client) GetUploadToken() (*UploadToken, error) {
	endpoint := "/imagex/gen_token"
	params := &url.Values{
		"client": []string{"web"},
	}
	raw, err := c.Get(endpoint, params)
	if err != nil {
		return nil, errors.Trace(err)
	}
	var token *UploadToken
	err = json.Unmarshal([]byte(gjson.Get(raw, "data.token").String()), &token)
	return token, errors.Trace(err)
}

func (ix *ImageX) ApplyImageUpload() (string, error) {
	rawurl := fmt.Sprintf("https://imagex.bytedanceapi.com/?Action=%s&Version=%s&ServiceId=%s",
		actionApplyImageUpload, version, serviceID)
	req, err := http.NewRequest(http.MethodGet, rawurl, nil)
	if err != nil {
		return "", errors.Trace(err)
	}

	if err := ix.signRequest(req); err != nil {
		return "", errors.Trace(err)
	}

	res, err := ix.getClient().Do(req)
	if err != nil {
		return "", errors.Trace(err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Trace(err)
	}
	raw := string(b)
	if res.StatusCode != 200 || gjson.Get(raw, "ResponseMetadata.Error").Exists() {
		return "", errors.Errorf("raw: %s, response: %+v", raw, res)
	}
	return raw, nil
}

func (ix *ImageX) CommitImageUpload(sessionKey string) (string, error) {
	rawurl := fmt.Sprintf("https://imagex.bytedanceapi.com/?Action=%s&Version=%s&SessionKey=%s&ServiceId=%s",
		actionCommitImageUpload, version, sessionKey, serviceID)
	req, err := http.NewRequest(http.MethodPost, rawurl, nil)
	if err != nil {
		return "", errors.Trace(err)
	}

	if err := ix.signRequest(req); err != nil {
		return "", errors.Trace(err)
	}

	res, err := ix.getClient().Do(req)
	if err != nil {
		return "", errors.Trace(err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Trace(err)
	}
	raw := string(b)
	if res.StatusCode != 200 || gjson.Get(raw, "ResponseMetadata.Error").Exists() {
		return "", errors.Errorf("raw: %s, response: %+v", raw, res)
	}
	return raw, nil
}

func (ix *ImageX) getClient() *http.Client {
	if ix.Client == nil {
		return http.DefaultClient
	}
	return ix.Client
}

func (ix *ImageX) signKeys(t time.Time) []byte {
	h := makeHMac([]byte("AWS4"+ix.SecretKey), []byte(t.Format(shortTimeFormat)))
	h = makeHMac(h, []byte(ix.Region))
	h = makeHMac(h, []byte(serviceName))
	h = makeHMac(h, []byte("aws4_request"))
	return h
}

func (ix *ImageX) writeRequest(w io.Writer, r *http.Request) error {
	r.Header.Set("host", r.Host)

	w.Write([]byte(r.Method))
	w.Write(newLine)
	writeURI(w, r)
	w.Write(newLine)
	writeQuery(w, r)
	w.Write(newLine)
	writeHeader(w, r)
	w.Write(newLine)
	w.Write(newLine)
	writeHeaderList(w, r)
	w.Write(newLine)
	err := writeBody(w, r)
	return errors.Trace(err)
}

func (ix *ImageX) writeStringToSign(w io.Writer, t time.Time, r *http.Request) error {
	w.Write([]byte(algorithm))
	w.Write(newLine)
	w.Write([]byte(t.Format(amzDateISO8601TimeFormat)))
	w.Write(newLine)

	w.Write([]byte(ix.creds(t)))
	w.Write(newLine)

	h := sha256.New()
	if err := ix.writeRequest(h, r); err != nil {
		return errors.Trace(err)
	}
	fmt.Fprintf(w, "%x", h.Sum(nil))
	return nil
}

func (ix *ImageX) creds(t time.Time) string {
	return t.Format(shortTimeFormat) + "/" + ix.Region + "/" + serviceName + "/aws4_request"
}

func (ix *ImageX) signRequest(req *http.Request) error {
	t := time.Now().UTC()
	req.Header.Set("x-amz-date", t.Format(amzDateISO8601TimeFormat))

	req.Header.Set("x-amz-security-token", ix.Token)

	k := ix.signKeys(t)
	h := hmac.New(sha256.New, k)

	if err := ix.writeStringToSign(h, t, req); err != nil {
		return errors.Trace(err)
	}

	auth := bytes.NewBufferString(algorithm)
	auth.Write([]byte(" Credential=" + ix.AccessKey + "/" + ix.creds(t)))
	auth.Write([]byte{',', ' '})
	auth.Write([]byte("SignedHeaders="))
	writeHeaderList(auth, req)
	auth.Write([]byte{',', ' '})
	auth.Write([]byte("Signature=" + fmt.Sprintf("%x", h.Sum(nil))))

	req.Header.Set("authorization", auth.String())
	return nil
}

func writeURI(w io.Writer, r *http.Request) {
	path := r.URL.RequestURI()
	if r.URL.RawQuery != "" {
		path = path[:len(path)-len(r.URL.RawQuery)-1]
	}
	slash := strings.HasSuffix(path, "/")
	path = filepath.Clean(path)
	if path != "/" && slash {
		path += "/"
	}
	w.Write([]byte(path))
}
func writeQuery(w io.Writer, r *http.Request) {
	var a []string
	for k, vs := range r.URL.Query() {
		k = url.QueryEscape(k)
		for _, v := range vs {
			if v == "" {
				a = append(a, k)
			} else {
				v = url.QueryEscape(v)
				a = append(a, k+"="+v)
			}
		}
	}
	sort.Strings(a)
	for i, s := range a {
		if i > 0 {
			w.Write([]byte{'&'})
		}
		w.Write([]byte(s))
	}
}

func writeHeader(w io.Writer, r *http.Request) {
	i, a := 0, make([]string, len(r.Header))
	for k, v := range r.Header {
		sort.Strings(v)
		a[i] = strings.ToLower(k) + ":" + strings.Join(v, ",")
		i++
	}
	sort.Strings(a)
	for i, s := range a {
		if i > 0 {
			w.Write(newLine)
		}
		io.WriteString(w, s)
	}
}

func writeHeaderList(w io.Writer, r *http.Request) {
	i, a := 0, make([]string, len(r.Header))
	for k := range r.Header {
		a[i] = strings.ToLower(k)
		i++
	}
	sort.Strings(a)
	for i, s := range a {
		if i > 0 {
			w.Write([]byte{';'})
		}
		w.Write([]byte(s))
	}
}

func writeBody(w io.Writer, r *http.Request) (err error) {
	var b []byte
	// If the payload is empty, use the empty string as the input to the SHA256 function
	// http://docs.amazonwebservices.com/general/latest/gr/sigv4-create-canonical-request.html
	if r.Body == nil {
		b = []byte("")
	} else {
		b, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return errors.Trace(err)
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	}

	h := sha256.New()
	h.Write(b)
	fmt.Fprintf(w, "%x", h.Sum(nil))
	return nil
}

func makeHMac(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

// Upload upload image to ByteDance Storage, support local file and web resource
func (ix *ImageX) Upload(uploadURL, path, auth string) error {
	var data []byte
	if isValidURL(path) {
		resp, err := http.Get(path)
		if err != nil {
			return errors.Trace(err)
		}
		defer resp.Body.Close()
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Trace(err)
		}
	} else {
		var err error
		data, err = ioutil.ReadFile(path)
		if err != nil {
			return errors.Trace(err)
		}
	}

	crc32, err := hashFileCRC32(bytes.NewBuffer(data))
	if err != nil {
		return errors.Trace(err)
	}

	req, err := http.NewRequest(http.MethodPost, uploadURL, bytes.NewBuffer(data))
	if err != nil {
		return errors.Trace(err)
	}
	req.Header.Add("authorization", auth)
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("content-crc32", crc32)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Trace(err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Trace(err)
	}
	raw := string(b)
	if gjson.Get(raw, "success").Int() != 0 {
		return errors.Errorf("raw: %s, response: %+v", raw, res)
	}
	return nil
}

// hashFileCRC32 generate CRC32 hash of a file
// Refer https://mrwaggel.be/post/generate-crc32-hash-of-a-file-in-golang-turorial/
func hashFileCRC32(r io.Reader) (string, error) {
	tablePolynomial := crc32.MakeTable(polynomialCRC32)
	hash := crc32.New(tablePolynomial)
	if _, err := io.Copy(hash, r); err != nil {
		return "", errors.Trace(err)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// https://golangcode.com/how-to-check-if-a-string-is-a-url/
func isValidURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}

	u, err := url.Parse(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}
