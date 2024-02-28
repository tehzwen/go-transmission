package go_transmission

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	TORRENT_GET        = "torrent-get"
	TORRENT_ADD        = "torrent-add"
	TORRENT_STOP       = "torrent-stop"
	TORRENT_START      = "torrent-start"
	TORRENT_REMOVE     = "torrent-remove"
	STATUS_DOWNLOADING = 4
	STATUS_STOPPED     = 0
)

func httpRequest(req *http.Request) ([]byte, *http.Response, error) {
	hc := &http.Client{
		Timeout: 5 * time.Second,
	}
	response, err := hc.Do(req)
	if err != nil {
		return []byte{}, nil, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, nil, err
	}
	return body, response, nil
}

type TransmissionRequest struct {
	Method    string                `json:"method"`
	Arguments TransmissionArguments `json:"arguments"`
}

type TransmissionArguments struct {
	Filename string    `json:"filename,omitempty"`
	MetaInfo string    `json:"metainfo,omitempty"`
	Torrents []Torrent `json:"torrents,omitempty"`
	Fields   []string  `json:"fields,omitempty"`
	Ids      []int     `json:"ids,omitempty"`
}

type TransmissionResponse struct {
	Result    string                `json:"result"`
	Arguments TransmissionArguments `json:"arguments"`
}

type TransmissionClient struct {
	username string
	password string
	endpoint string
	port     int

	SessionId string
}

func (t *TransmissionClient) getUrl() string {
	return fmt.Sprintf(
		"http://%s:%s@%s:%d/transmission/rpc",
		t.username,
		t.password,
		t.endpoint,
		t.port,
	)
}

func (t *TransmissionClient) makeRequest(req TransmissionRequest) (*TransmissionResponse, error) {
	jsonPayload, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(
		"POST",
		t.getUrl(),
		bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	request.Header.Set("X-Transmission-Session-Id", t.SessionId)

	body, _, err := httpRequest(request)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s\n", string(body))
	var r TransmissionResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (t *TransmissionClient) Login() error {
	request, err := http.NewRequest(
		"GET",
		t.getUrl(),
		nil)
	if err != nil {
		return err
	}
	_, resp, err := httpRequest(request)
	if err != nil {
		return err
	}

	sessionId := resp.Header.Get("X-Transmission-Session-Id")
	t.SessionId = sessionId
	return nil
}

func (t *TransmissionClient) GetTorrent(id int) (*TransmissionResponse, error) {
	payload := TransmissionRequest{
		Method: TORRENT_GET,
		Arguments: TransmissionArguments{
			Fields: []string{"id", "name", "status", "percentDone"},
		},
	}
	return t.makeRequest(payload)
}

func (t *TransmissionClient) GetTorrents() (*TransmissionResponse, error) {
	payload := TransmissionRequest{
		Method: TORRENT_GET,
		Arguments: TransmissionArguments{
			Fields: []string{"id", "name", "status", "percentDone"},
		},
	}
	return t.makeRequest(payload)
}

func (t *TransmissionClient) AddTorrent(magnetLink string, opts ...func(*TransmissionRequest)) (*TransmissionResponse, error) {
	// TODO - validate input properly
	payload := TransmissionRequest{
		Method: TORRENT_ADD,
		Arguments: TransmissionArguments{
			Filename: magnetLink,
		},
	}

	for _, o := range opts {
		o(&payload)
	}

	return t.makeRequest(payload)
}

func (t *TransmissionClient) StopTorrent(id int) (*TransmissionResponse, error) {
	payload := TransmissionRequest{
		Method: TORRENT_STOP,
		Arguments: TransmissionArguments{
			Ids: []int{id},
		},
	}
	return t.makeRequest(payload)
}

func (t *TransmissionClient) StartTorrent(id int) (*TransmissionResponse, error) {
	payload := TransmissionRequest{
		Method: TORRENT_START,
		Arguments: TransmissionArguments{
			Ids: []int{id},
		},
	}
	return t.makeRequest(payload)
}

func (t *TransmissionClient) RemoveTorrent(id int) (*TransmissionResponse, error) {
	payload := TransmissionRequest{
		Method: TORRENT_REMOVE,
		Arguments: TransmissionArguments{
			Ids: []int{id},
		},
	}
	return t.makeRequest(payload)
}

func NewTransmissionClient(user, pw, endpoint string, port int) *TransmissionClient {
	tc := &TransmissionClient{
		endpoint: endpoint,
		username: user,
		password: pw,
		port:     port,
	}

	return tc
}
