package fuyou

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/config"
	"dev-gitlab.wanxingrowth.com/fanli/fuyou-payment-gateway/pkg/utils/log"
)

func FYVerify(data map[string]interface{}, sign string) (err error) {

	// 测试环境不验证来源
	if config.Config.Environment == config.Test {
		return nil
	}

	b, err := ioutil.ReadFile(config.Config.FuYou.GetPublicPemPath())
	if err != nil {
		return
	}

	key, err := genPublicKey(b)
	if err != nil {
		return
	}

	s, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return
	}

	sum := md5.Sum([]byte(genFySignString(data)))
	return rsa.VerifyPKCS1v15(key, crypto.MD5, sum[:], s)
}

func Sign(m map[string]interface{}) (sign string, err error) {

	signString := genFySignString(m)

	log.GetLogger().WithField("signString", signString).Debug("get sign string")

	content, err := ioutil.ReadFile(config.Config.FuYou.GetPrivatePemPath())
	if err != nil {
		return
	}

	key, err := genPrivateKey(content)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(signString)), simplifiedchinese.GBK.NewEncoder()))
	if err != nil {
		return
	}

	sum := md5.Sum(data)
	signedByte, err := rsa.SignPKCS1v15(nil, key, crypto.MD5, sum[:])
	if err != nil {
		return
	}

	sign = base64.StdEncoding.EncodeToString(signedByte)
	return
}

func genPrivateKey(pemResource []byte) (privateKey *rsa.PrivateKey, err error) {

	block, _ := pem.Decode(pemResource)
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	return
}

func genPublicKey(pemResource []byte) (publicKey *rsa.PublicKey, err error) {

	block, _ := pem.Decode(pemResource)
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}

	publicKey = key.(*rsa.PublicKey)
	return
}

// 1.将接口中每一个字段(sign 及reserved 开头字段除外),以字典顺序排
// 序之后,按照key1=value1&key2=value2.....的顺序,进行拼接。
// 2.对得到的字符串进行RSA 签名/验签
// 注：sign 及reserved 开头字段除外的其他非必填字段也需要参与验签。
func genFySignString(m map[string]interface{}) string {
	keys := make([]string, len(m))
	i := -1
	for k := range m {
		i++
		if strings.EqualFold(strings.ToLower(k), "sign") || strings.HasPrefix(strings.ToLower(k), "reserved") {
			continue
		}
		keys[i] = k
	}
	sort.Strings(keys)
	var s []string
	for i := 0; i < len(keys); i++ {
		if keys[i] != "" {
			s = append(s, strings.ToLower(keys[i])+"="+fmt.Sprintf("%v", m[keys[i]]))
		}
	}
	s1 := strings.Join(s, "&")
	return s1
}
