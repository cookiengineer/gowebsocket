package gowebsocket

import "bufio"
import "crypto/tls"
import "fmt"
import "net"
import net_http "net/http"
import net_url "net/url"
import "strings"

type Client struct {
	Socket *WebSocket
	URL    *net_url.URL
}

func NewClient(raw_url string) (*Client, error) {

	url, err1 := net_url.Parse(raw_url)

	if err1 == nil {

		scheme := strings.ToLower(url.Scheme)

		if scheme == "ws" || scheme == "wss" {

			host := url.Host
			path := url.Path

			if strings.Contains(host, ":") == false {

				if scheme == "ws" {
					host = host + ":80"
				} else if scheme == "wss" {
					host = host + ":443"
				}

			}

			if path == "" {
				path = "/"
			}

			if url.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, url.RawQuery)
			}

			nonce_key, err2 := generateNonceKey()

			if err2 == nil {

				var connection net.Conn = nil
				var http_url string = ""

				if scheme == "wss" {

					tmp, err := tls.Dial("tcp", host, &tls.Config{})

					if err == nil {

						request_url = fmt.Sprintf("https://%s%s", host, path)
						connection  = &tmp

					} else {
						return nil, fmt.Errorf("websocket: tls dial failed: %s", err)
					}

				} else if scheme == "ws" {

					tmp, err := net.Dial("tcp", host)

					if err == nil {

						request_url = fmt.Sprintf("http://%s%s", host, path)
						connection  = &tmp

					} else {
						return nil, fmt.Errorf("websocket: tcp dial failed: %s", err)
					}

				}

				if request_url != "" && connection != nil {

					request, err3 := net_http.NewRequest(net_http.MethodGet, request_url, nil)

					if err3 == nil {

						request.Header.Set("Upgrade", "websocket")
						request.Header.Set("Connection", "Upgrade")
						request.Header.Set("Sec-WebSocket-Key", nonce_key)
						request.Header.Set("Sec-WebSocket-Version", "13")

						err4 := request.Write(*connection)

						if err4 == nil {

							pipe           := &pipe_connection{Conn: connection}
							reader         := bufio.NewReader(pipe)
							response, err5 := net_http.ReadResponse(reader, request)

							if err5 == nil {

								err6 := verifyUpgradeResponse(nonce_key, response)

								if err6 == nil {

									remaining := reader.Buffered()

									if remaining > 0 {
										// Recover remaining WebSocket frames
										pipe.buffer = make([]byte, remaining)
										reader.Read(pipe.buffer)
									}

									websocket := NewWebSocket(net.Conn(pipe), nil)
									client    := &Client{
										Socket: websocket,
										URL:    url,
									}

									return client, nil

								} else {

									defer connection.Close()
									return nil, fmt.Errorf("websocket: upgrade failed: %s", err6)

								}

							} else {

								defer connection.Close()
								return nil, fmt.Errorf("websocket: failed to read response: %s", err5)

							}

						} else {

							defer connection.Close()
							return nil, fmt.Errorf("websocket: failed to write request: %s", err4)

						}

					} else {

						defer connection.Close()
						return nil, fmt.Errorf("websocket: failed to create HTTP request: %s", err3)

					}

				} else {
					return nil, fmt.Errorf("websocket: net dial failed")
				}

			} else {
				return nil, fmt.Errorf("websocket: failed to generate Nonce key: %s", err2)
			}

		} else {
			return nil, fmt.Errorf("websocket: unsupported URL scheme: %s", url.Scheme)
		}

	} else {
		return nil, fmt.Errorf("websocket: invalid URL: %s", err1)
	}

}
