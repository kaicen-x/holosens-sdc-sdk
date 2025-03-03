package digest

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

// ParseDigestWwwAuthenticate 解析HTTP响应头中的认证摘要信息
func ParseDigestWwwAuthenticate(header http.Header) (realm, nonce, algorithm string) {
	// 提取响应头
	wwwAuthenticate := header.Values("WWW-Authenticate")
	// 提取鉴权信息
	for _, item := range wwwAuthenticate {
		// 移除前置字符
		item = strings.TrimPrefix(item, "Digest")
		// 移除前后空格
		item = strings.TrimSpace(item)
		// 使用逗号分隔
		for _, e := range strings.Split(item, ",") {
			// 使用等号分割
			kv := strings.Split(e, "=")
			if len(kv) == 2 {
				switch strings.TrimSpace(kv[0]) {
				case "realm":
					realm = strings.Trim(strings.TrimSpace(kv[1]), "\"")
				case "nonce":
					nonce = strings.Trim(strings.TrimSpace(kv[1]), "\"")
				case "algorithm":
					algorithm = strings.Trim(strings.TrimSpace(kv[1]), "\"")
				}
			}
		}
	}

	// OK
	return
}

// MakeDigestAuthorization 构建HTTP请求认证摘要
func MakeDigestAuthorization(method, uri, realm, nonce, algorithm string, username, password string) string {
	// 声明加密数据
	a1 := fmt.Sprintf("%s:%s:%s", username, realm, password)
	a2 := fmt.Sprintf("%s:%s", method, uri)
	// 生成摘要response
	var response string
	// 区分算法
	switch strings.ToUpper(algorithm) {
	// 使用md5算法
	case "MD5":
		a1Md5 := md5.Sum([]byte(a1))
		a2Md5 := md5.Sum([]byte(a2))
		group := fmt.Sprintf("%s:%s:%s", hex.EncodeToString(a1Md5[:]), nonce, hex.EncodeToString(a2Md5[:]))
		resMd5 := md5.Sum([]byte(group))
		response = hex.EncodeToString(resMd5[:])

	// 使用sha256算法
	case "SHA-256":
		a1Sha256 := sha256.Sum256([]byte(a1))
		a2Sha256 := sha256.Sum256([]byte(a2))
		group := fmt.Sprintf("%s:%s:%s", hex.EncodeToString(a1Sha256[:]), nonce, hex.EncodeToString(a2Sha256[:]))
		resSha256 := sha256.Sum256([]byte(group))
		response = hex.EncodeToString(resSha256[:])

	// 不支持的算法
	default:
		return ""
	}
	// 填充摘要
	metaData := []string{
		fmt.Sprintf(`username="%s"`, username),
		fmt.Sprintf(`realm="%s"`, realm),
		fmt.Sprintf(`nonce="%s"`, nonce),
		fmt.Sprintf(`uri="%s"`, uri),
		fmt.Sprintf(`algorithm="%s"`, algorithm),
		fmt.Sprintf(`response="%s"`, response),
	}
	// 拼接认证摘要
	return fmt.Sprintf("Digest %s", strings.Join(metaData, ", "))
}
