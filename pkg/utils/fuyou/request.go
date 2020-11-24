package fuyou

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"

	"github.com/sirupsen/logrus"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/constant"
	transUtil "dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/transform"
)

type Request struct {
	logger *logrus.Entry
}

func NewRequest(logger *logrus.Entry) *Request {

	return &Request{logger: logger}
}

func (req *Request) SendRequest(data map[string]interface{}, url string) (responseData map[string]string, err error) {

	sign, err := Sign(data)
	if err != nil {
		req.logger.WithError(err).Error("signature error")
		return
	}

	data["sign"] = sign

	xml, err := req.toXml(data)
	if err != nil {
		req.logger.WithError(err).Error("change to xml error")
		return
	}

	req.logger.WithField("requestData", xml).WithField("url", url).Info("send request")

	response, err := req.requestFuYou(url, xml)
	if err != nil {
		req.logger.WithError(err).Error("request fuyou error")
		return
	}

	req.logger = req.logger.WithField("responseData", response)

	responseData, err = req.parseResponseToMap(response)
	if err != nil {
		req.logger.WithError(err).Error("parse response to map error")
		return
	}

	return
}

func (req *Request) toXml(prePayMap map[string]interface{}) (xmlContent string, err error) {

	xmlGenerator := transUtil.NewXmlGenerator("xml")
	keys := make([]string, 0, len(prePayMap))
	for key := range prePayMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {

		v := prePayMap[key]
		if v != "" {
			xmlGenerator.AddTag(key, fmt.Sprintf("%v", v))
		}
	}
	xml, err := xmlGenerator.ToXml()
	if err != nil {
		req.logger.WithError(err).Error("generate xml error")
		return
	}
	xmlContent = fmt.Sprintf("%s\n%s", constant.FuYouXMLHeader, xml)
	return
}

func (req *Request) requestFuYou(requestURL string, reqBody string) (res string, err error) {

	logger := req.logger.WithField("url", requestURL).WithField("requestData", reqBody)
	requestData, err := transUtil.Encode(reqBody, "gbk")
	if err != nil {
		logger.WithError(err).Error("transform encode error")
		return
	}

	data := url.Values{"req": {requestData}}

	logger = logger.WithField("body", data)

	logger.Infof("request data")

	var response *http.Response
	response, err = http.DefaultClient.PostForm(requestURL, data)
	if err != nil {
		logger.WithError(err).Error("request fuyou error")
		return
	}

	var body []byte
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		logger.WithError(err).Error("read body error")
		return
	}

	res, err = transUtil.Decode(string(body), "gbk")
	if err != nil {
		logger.WithError(err).Error("decode body error")
		return
	}

	logger.WithField("response data", res).Info("response")
	return
}

func (req *Request) parseResponseToMap(response string) (map[string]string, error) {

	return transUtil.Xml2mapWithRoot([]byte(response), "xml")
}
