package upstream

import (
	"io"
	"net/http"
	"strings"
)

type StoredResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func CopyHeaders(dst, src http.Header) {
	for key, values := range src {
		if isHopByHopHeader(key) {
			continue
		}
		dst.Del(key)
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func CloneHeaders(src http.Header) http.Header {
	dst := make(http.Header, len(src))
	CopyHeaders(dst, src)
	return dst
}

func ReadRequestBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}
	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func CaptureResponse(resp *http.Response) (*StoredResponse, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &StoredResponse{
		StatusCode: resp.StatusCode,
		Header:     CloneHeaders(resp.Header),
		Body:       body,
	}, nil
}

func CopyResponseBody(w http.ResponseWriter, body io.Reader, flush bool, tap io.Writer) error {
	buf := make([]byte, 32*1024)
	flusher, _ := w.(http.Flusher)

	for {
		n, err := body.Read(buf)
		if n > 0 {
			if tap != nil {
				if _, tapErr := tap.Write(buf[:n]); tapErr != nil {
					return tapErr
				}
			}
			if _, writeErr := w.Write(buf[:n]); writeErr != nil {
				return writeErr
			}
			if flush && flusher != nil {
				flusher.Flush()
			}
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func ShouldFailover(status int) bool {
	if status >= 500 {
		return true
	}
	switch status {
	case http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusRequestTimeout,
		http.StatusTooEarly,
		http.StatusTooManyRequests:
		return true
	default:
		return false
	}
}

func JoinPaths(basePath, reqPath string) string {
	basePath = strings.TrimSuffix(basePath, "/")
	reqPath = "/" + strings.TrimPrefix(reqPath, "/")
	if basePath == "" {
		return reqPath
	}
	return basePath + reqPath
}

func isHopByHopHeader(key string) bool {
	switch strings.ToLower(key) {
	case "connection", "keep-alive", "proxy-authenticate", "proxy-authorization", "te", "trailers", "transfer-encoding", "upgrade":
		return true
	default:
		return false
	}
}
