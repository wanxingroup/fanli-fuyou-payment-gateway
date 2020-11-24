package fuyou

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/log"
)

func TestRequest_SendRequest(t *testing.T) {

	logger := log.GetLogger()
	logger.Level = logrus.DebugLevel

	tests := []struct {
		inputData    map[string]interface{}
		inputURL     string
		wantResponse map[string]string
		wantError    error
	}{
		{
			inputData: map[string]interface{}{
				"version": "1.0",
				"test":    "data",
			},
			inputURL: "http://dev-api.wanxingrowth.com/mock/131/wxPreCreate",
			wantResponse: map[string]string{
				"ins_cd":                    "08M0061071",
				"mchnt_cd":                  "0001000F2214654",
				"qr_code":                   "",
				"random_str":                "O4H5M1CXV1L24D0IS4NIIDMREYBN9I9U",
				"reserved_addn_inf":         "",
				"reserved_channel_order_id": "",
				"reserved_fy_order_no":      "",
				"reserved_fy_settle_dt":     "",
				"reserved_fy_trace_no":      "900337123423",
				"reserved_pay_info":         `{"appId":"wxfa089da95020ba1a","timeStamp":"1557649678","signType":"RSA","package":"prepay_id=wx121627586241613c7d311e003087612412","nonceStr":"683cedb4a5894ca9a681d4172ee6dab0","paySign":"hvMJreiPBOCrG45XppgPxmZnxL/mojcgv4aKKU0OGJNB5Cf/0e9CwhF8Uy4OOmjuwv4l/65xvCbom0UcNwJpQ9DQUc8lMRwYUwgKWsXgVcpDW0rw3yZLX7otQJQi8JGBagq3si55yIyca+1Pjg1sOXuRnh+/2yRHPj3wR2A+tqSM1Rl5Q0PrKmWgDocIU/5wIC4B/8uK7Kf+K/HtImaag/TTNsEKg/TEzzEshfhEVl5dYvMcoq1D4FkH/938eVmpUMqGI3DQMnBa1Oj6T12sjF9p/aPdV7S9OyVxPHUGNMqDoouXnACTd4HmkU3+G7/7oyYyozi2V+JDPX1lrO+Y1w=="}`,
				"reserved_transaction_id":   "",
				"result_code":               "000000",
				"result_msg":                "SUCCESS",
				"sdk_appid":                 "wxfa089da95020ba1a",
				"sdk_noncestr":              "683cedb4a5894ca9a681d4172ee6dab0",
				"sdk_package":               "prepay_id=wx121627586241613c7d311e003087612412",
				"sdk_partnerid":             "",
				"sdk_paysign":               "hvMJreiPBOCrG45XppgPxmZnxL/mojcgv4aKKU0OGJNB5Cf/0e9CwhF8Uy4OOmjuwv4l/65xvCbom0UcNwJpQ9DQUc8lMRwYUwgKWsXgVcpDW0rw3yZLX7otQJQi8JGBagq3si55yIyca+1Pjg1sOXuRnh+/2yRHPj3wR2A+tqSM1Rl5Q0PrKmWgDocIU/5wIC4B/8uK7Kf+K/HtImaag/TTNsEKg/TEzzEshfhEVl5dYvMcoq1D4FkH/938eVmpUMqGI3DQMnBa1Oj6T12sjF9p/aPdV7S9OyVxPHUGNMqDoouXnACTd4HmkU3+G7/7oyYyozi2V+JDPX1lrO+Y1w==",
				"sdk_signtype":              "RSA",
				"sdk_timestamp":             "1557649678",
				"session_id":                "wx121627586241613c7d311e003087612412",
				"sign":                      "KUWssvgIiL0AyYnVKigZY66xyHX25RWDKRMCukR2uRt9WA2MUIDu0Yp+PTGpF34XJj5Z575l77229vTNn1jzEJjupRWZDfGFziQk4OVneSEk6IEFz7d/Mpygzgpc8kRQlN1EtNy8eZaLwv0nVDyopVBuVQDQFkzNMjSN5q/Ba8s=",
				"sub_appid":                 "wxfa089da95020ba1a",
				"sub_mer_id":                "281826820",
				"sub_openid":                "ooIeqs5VwPJnDUYfLweOKcR5AxpE",
				"term_id":                   "",
			},

			wantError: nil,
		},
	}

	for _, test := range tests {

		request := NewRequest(logger)
		response, err := request.SendRequest(test.inputData, test.inputURL)
		assert.Equal(t, test.wantResponse, response, test)
		assert.Equal(t, test.wantError, err, test)
	}
}

func TestRequest_ToXML(t *testing.T) {

	logger := log.GetLogger()
	logger.Level = logrus.DebugLevel

	tests := []struct {
		input map[string]interface{}
		want  string
		err   error
	}{
		{
			input: map[string]interface{}{
				"key1": "value1",
				"key2": "value2",
			},
			want: `<?xml version="1.0" encoding="GBK" standalone="yes"?>
<xml>
  <key1>value1</key1>
  <key2>value2</key2>
</xml>
`,
		},
	}

	for _, test := range tests {

		request := NewRequest(logger)
		result, err := request.toXml(test.input)
		assert.Equal(t, test.want, result, test)
		assert.Equal(t, test.err, err, test)
	}
}

func TestRequest_RequestFuYou(t *testing.T) {

	logger := log.GetLogger()
	logger.Level = logrus.DebugLevel

	tests := []struct {
		inputData    string
		inputURL     string
		wantResponse string
		wantError    error
	}{
		{
			inputData: "test",
			inputURL:  "http://dev-api.wanxingrowth.com/mock/131/wxPreCreate",
			wantResponse: `<?xml version="1.0" encoding="UTF-8"?>
<xml>
    <ins_cd>08M0061071</ins_cd>
    <mchnt_cd>0001000F2214654</mchnt_cd>
    <qr_code></qr_code>
    <random_str>O4H5M1CXV1L24D0IS4NIIDMREYBN9I9U</random_str>
    <reserved_addn_inf></reserved_addn_inf>
    <reserved_channel_order_id></reserved_channel_order_id>
    <reserved_fy_order_no></reserved_fy_order_no>
    <reserved_fy_settle_dt></reserved_fy_settle_dt>
    <reserved_fy_trace_no>900337123423</reserved_fy_trace_no>
    <reserved_pay_info>{"appId":"wxfa089da95020ba1a","timeStamp":"1557649678","signType":"RSA","package":"prepay_id=wx121627586241613c7d311e003087612412","nonceStr":"683cedb4a5894ca9a681d4172ee6dab0","paySign":"hvMJreiPBOCrG45XppgPxmZnxL/mojcgv4aKKU0OGJNB5Cf/0e9CwhF8Uy4OOmjuwv4l/65xvCbom0UcNwJpQ9DQUc8lMRwYUwgKWsXgVcpDW0rw3yZLX7otQJQi8JGBagq3si55yIyca+1Pjg1sOXuRnh+/2yRHPj3wR2A+tqSM1Rl5Q0PrKmWgDocIU/5wIC4B/8uK7Kf+K/HtImaag/TTNsEKg/TEzzEshfhEVl5dYvMcoq1D4FkH/938eVmpUMqGI3DQMnBa1Oj6T12sjF9p/aPdV7S9OyVxPHUGNMqDoouXnACTd4HmkU3+G7/7oyYyozi2V+JDPX1lrO+Y1w=="}</reserved_pay_info>
    <reserved_transaction_id></reserved_transaction_id>
    <result_code>000000</result_code>
    <result_msg>SUCCESS</result_msg>
    <sdk_appid>wxfa089da95020ba1a</sdk_appid>
    <sdk_noncestr>683cedb4a5894ca9a681d4172ee6dab0</sdk_noncestr>
    <sdk_package>prepay_id=wx121627586241613c7d311e003087612412</sdk_package>
    <sdk_partnerid></sdk_partnerid>
    <sdk_paysign>hvMJreiPBOCrG45XppgPxmZnxL/mojcgv4aKKU0OGJNB5Cf/0e9CwhF8Uy4OOmjuwv4l/65xvCbom0UcNwJpQ9DQUc8lMRwYUwgKWsXgVcpDW0rw3yZLX7otQJQi8JGBagq3si55yIyca+1Pjg1sOXuRnh+/2yRHPj3wR2A+tqSM1Rl5Q0PrKmWgDocIU/5wIC4B/8uK7Kf+K/HtImaag/TTNsEKg/TEzzEshfhEVl5dYvMcoq1D4FkH/938eVmpUMqGI3DQMnBa1Oj6T12sjF9p/aPdV7S9OyVxPHUGNMqDoouXnACTd4HmkU3+G7/7oyYyozi2V+JDPX1lrO+Y1w==</sdk_paysign>
    <sdk_signtype>RSA</sdk_signtype>
    <sdk_timestamp>1557649678</sdk_timestamp>
    <session_id>wx121627586241613c7d311e003087612412</session_id>
    <sign>KUWssvgIiL0AyYnVKigZY66xyHX25RWDKRMCukR2uRt9WA2MUIDu0Yp+PTGpF34XJj5Z575l77229vTNn1jzEJjupRWZDfGFziQk4OVneSEk6IEFz7d/Mpygzgpc8kRQlN1EtNy8eZaLwv0nVDyopVBuVQDQFkzNMjSN5q/Ba8s=</sign>
    <sub_appid>wxfa089da95020ba1a</sub_appid>
    <sub_mer_id>281826820</sub_mer_id>
    <sub_openid>ooIeqs5VwPJnDUYfLweOKcR5AxpE</sub_openid>
    <term_id></term_id>
</xml>`,
		},
	}

	for _, test := range tests {

		request := NewRequest(logger)
		result, err := request.requestFuYou(test.inputURL, test.inputData)
		assert.Equal(t, test.wantResponse, result, test)
		assert.Equal(t, test.wantError, err, test)
	}
}

func TestRequest_ParseResponseToMap(t *testing.T) {

	logger := log.GetLogger()
	logger.Level = logrus.DebugLevel

	tests := []struct {
		input string
		want  map[string]string
		err   error
	}{
		{
			input: `<?xml version="1.0" encoding="GBK" standalone="yes"?>
<xml>
  <key1>value1</key1>
  <key2>value2</key2>
</xml>
`,
			want: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	for _, test := range tests {

		request := NewRequest(logger)
		result, err := request.parseResponseToMap(test.input)
		assert.Equal(t, test.want, result, test)
		assert.Equal(t, test.err, err, test)
	}
}
