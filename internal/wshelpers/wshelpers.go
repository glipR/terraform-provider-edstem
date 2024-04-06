package wshelpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"

	"terraform-provider-edstem/internal/client"
)

type TicketResponse struct {
	Ticket string `json:"ticket"`
}

type Message struct {
	Type string `json:"type"`
}

type MessageData struct {
	Type   string `json:"type"`
	Param1 string `json:"param1"`
	Param2 string `json:"param2"`
	Param3 string `json:"param3"`
}

type FSOPRequest struct {
	Type string      `json:"type"`
	Data MessageData `json:"data"`
}

type ListingReply struct {
	Type string      `json:"type"`
	Data ListingData `json:"data"`
}

type ListingData struct {
	Listing []ListingEntry `json:"listing"`
	Dir     string         `json:"dir"`
}

type ListingEntry struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type FileOpenCommand struct {
	Type string       `json:"type"`
	Data FileOpenData `json:"data"`
}

type FileOpenData struct {
	Path string `json:"path"`
	Soft bool   `json:"soft"`
}

type FileOTInitRespose struct {
	Type string         `json:"type"`
	Data FileOTInitData `json:"data"`
}

type FileOTInitData struct {
	FID      int    `json:"fid"`
	Revision int    `json:"rev"`
	Buffer   string `json:"buffer"`
}

type FileOTCommandCursor struct {
	Type string           `json:"type"`
	Data FileOTDataCursor `json:"data"`
}

type FileOTDataCursor struct {
	FID        int            `json:"fid"`
	Revision   int            `json:"rev"`
	CursorData FileCursorData `json:"cursor"`
}

type FileCursorData struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type FileOTWriteOp struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type FileOTWriteData struct {
	FID        int             `json:"fid"`
	Revision   int             `json:"rev"`
	Operations []FileOTWriteOp `json:"op"`
}

type FileOTWriteCommand struct {
	Type string          `json:"type"`
	Data FileOTWriteData `json:"data"`
}

func DeleteAllFiles(conn *websocket.Conn) error {
	var req FSOPRequest
	req.Type = "fsop"
	req.Data.Type = "list_folder"
	req.Data.Param1 = "/home"

	req_body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, req_body)
	if err != nil {
		fmt.Printf("Write error\n")
		return err
	}

	content, get_err := GetMessage(conn, "list_reply")
	if get_err != nil {
		return get_err
	}

	m_resp := &ListingReply{}
	err = json.Unmarshal(content, &m_resp)
	if err != nil {
		return err
	}

	for _, returned := range m_resp.Data.Listing {
		var req FSOPRequest
		req.Type = "fsop"
		req.Data.Type = "remove"
		req.Data.Param1 = fmt.Sprintf("/home/%s", returned.Name)
		req.Data.Param3 = returned.Type

		req_body, err := json.Marshal(req)
		if err != nil {
			return err
		}

		err = conn.WriteMessage(websocket.BinaryMessage, req_body)
		if err != nil {
			fmt.Printf("Write error\n")
			return err
		}
	}

	return nil
}

func CreateDir(conn *websocket.Conn, relative_path string) error {
	var req FSOPRequest
	req.Type = "fsop"
	req.Data.Type = "new_folder"
	req.Data.Param1 = strings.ReplaceAll(filepath.Join("/home", relative_path), "\\", "/")
	req.Data.Param2 = ""
	req.Data.Param3 = ""

	req_body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, req_body)
	if err != nil {
		fmt.Printf("Write error\n")
		return err
	}

	return nil
}

func WriteFileContents(conn *websocket.Conn, relative_path string, file_contents string) error {
	// This needs to:
	// 1: Create the file
	// 2: Get the FID
	// 3: Set the cursor position
	// 4: Insert the text
	var req FSOPRequest
	req.Type = "fsop"
	req.Data.Type = "new_file"
	req.Data.Param1 = strings.ReplaceAll(filepath.Join("/home", relative_path), "\\", "/")
	req.Data.Param2 = ""
	req.Data.Param3 = ""

	req_body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, req_body)
	if err != nil {
		fmt.Printf("Write error\n")
		return err
	}

	var new_req FileOpenCommand
	new_req.Type = "file_open"
	new_req.Data.Path = strings.ReplaceAll(filepath.Join("/home", relative_path), "\\", "/")
	new_req.Data.Soft = true

	req_body, err = json.Marshal(new_req)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, req_body)
	if err != nil {
		fmt.Printf("Write error\n")
		return err
	}

	content, get_err := GetMessage(conn, "file_ot_init")
	if get_err != nil {
		return get_err
	}

	ot_resp := &FileOTInitRespose{}
	err = json.Unmarshal(content, &ot_resp)
	if err != nil {
		return err
	}

	var ot_cursor FileOTCommandCursor
	ot_cursor.Type = "file_ot"
	ot_cursor.Data.FID = ot_resp.Data.FID
	ot_cursor.Data.Revision = ot_resp.Data.Revision
	ot_cursor.Data.CursorData.Start = 0
	ot_cursor.Data.CursorData.End = 0

	req_body, err = json.Marshal(ot_cursor)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, req_body)
	if err != nil {
		fmt.Printf("Write error\n")
		return err
	}

	var ot_write FileOTWriteCommand
	ot_write.Type = "file_ot"
	ot_write.Data.FID = ot_resp.Data.FID
	ot_write.Data.Revision = ot_resp.Data.Revision
	var ot_write_1 FileOTWriteOp
	ot_write_1.Type = "insert"
	ot_write_1.Value = file_contents
	ot_write.Data.Operations = append(ot_write.Data.Operations, ot_write_1)

	req_body, err = json.Marshal(ot_write)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, req_body)
	if err != nil {
		fmt.Printf("Write error\n")
		return err
	}

	return nil
}

func GetMessage(conn *websocket.Conn, message_type string) ([]byte, error) {
	cur_type := ""
	var mcontent []byte
	var err error

	for cur_type != message_type {
		_, mcontent, err = conn.ReadMessage()
		if err != nil {
			fmt.Printf("Read\n error")
			return nil, err
		}
		m_resp := &Message{}
		err = json.Unmarshal(mcontent, &m_resp)
		if err != nil {
			return nil, err
		}

		cur_type = m_resp.Type
	}

	return mcontent, nil
}

func UpdateChallengeRepo(conn *client.Client, challenge_id int, challenge_folder_path string, repo_name string) error {
	body, err := conn.HTTPRequest(fmt.Sprintf("challenges/%d/connect/%s", challenge_id, repo_name), "POST", bytes.Buffer{}, nil)
	if err != nil {
		return err
	}
	resp := &TicketResponse{}
	err = json.NewDecoder(body).Decode(resp)
	if err != nil {
		return err
	}

	c, _, ws_err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://sahara.au.edstem.org/connect?ticket=%s", resp.Ticket), nil)
	if ws_err != nil {
		return ws_err
	}
	defer c.Close()

	cur_type := "init"

	for cur_type != "client_join" {
		_, mcontent, merr := c.ReadMessage()
		if merr != nil {
			fmt.Printf("Read\n error")
			return merr
		}
		m_resp := &Message{}
		err = json.Unmarshal(mcontent, &m_resp)
		if err != nil {
			return err
		}

		cur_type = m_resp.Type
	}

	DeleteAllFiles(c)

	err = filepath.Walk(filepath.Join(challenge_folder_path, repo_name),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			var rel_path string
			rel_path, err = filepath.Rel(filepath.Join(challenge_folder_path, repo_name), path)
			if err != nil {
				return err
			}
			if rel_path == "." {
				return nil
			}
			if info.IsDir() {
				fmt.Printf("Making Dir %s\n", rel_path)
				err = CreateDir(c, rel_path)
				if err != nil {
					return err
				}
			} else {
				fmt.Printf("Writing File %s\n", rel_path)
				dat, read_err := os.ReadFile(path)
				if read_err != nil {
					return read_err
				}
				fmt.Printf("File contents %s\n", string(dat))
				read_err = WriteFileContents(c, rel_path, string(dat))
				if read_err != nil {
					return read_err
				}
			}
			return nil
		})
	return err
}

func ReadChallengeRepo(conn *client.Client, challenge_id int, challenge_folder_path string, repo_name string) error {
	body, err := conn.HTTPRequest(fmt.Sprintf("challenges/%d/connect/%s", challenge_id, repo_name), "POST", bytes.Buffer{}, nil)
	if err != nil {
		return err
	}
	resp := &TicketResponse{}
	err = json.NewDecoder(body).Decode(resp)
	if err != nil {
		return err
	}

	c, _, ws_err := websocket.DefaultDialer.Dial(fmt.Sprintf("wss://sahara.au.edstem.org/connect?ticket=%s", resp.Ticket), nil)
	if ws_err != nil {
		return ws_err
	}
	defer c.Close()

	GetMessage(c, "client_join")

	var req FSOPRequest
	req.Type = "fsop"
	req.Data.Type = "list_folder"
	req.Data.Param1 = "/home"

	req_body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = c.WriteMessage(websocket.BinaryMessage, req_body)
	if err != nil {
		fmt.Printf("Write error\n")
		return err
	}

	content, get_err := GetMessage(c, "list_reply")
	if get_err != nil {
		return get_err
	}

	m_resp := &ListingReply{}
	err = json.Unmarshal(content, &m_resp)
	if err != nil {
		return err
	}

	if len(m_resp.Data.Listing) != 0 {
		os.MkdirAll(path.Join(challenge_folder_path, repo_name), os.ModeDir)
	}

	for _, returned := range m_resp.Data.Listing {
		err = RecReadPath(c, fmt.Sprintf("/home/%s", returned.Name), path.Join(challenge_folder_path, repo_name, returned.Name), returned.Type != "file")
		if err != nil {
			return err
		}
	}

	return nil
}

func RecReadPath(conn *websocket.Conn, web_path string, local_path string, is_dir bool) error {
	if !is_dir {
		f, e := os.Create(local_path)
		if e != nil {
			return e
		}
		var new_req FileOpenCommand
		new_req.Type = "file_open"
		new_req.Data.Path = web_path
		new_req.Data.Soft = true

		req_body, err := json.Marshal(new_req)
		if err != nil {
			return err
		}

		err = conn.WriteMessage(websocket.BinaryMessage, req_body)
		if err != nil {
			fmt.Printf("Write error\n")
			return err
		}

		content, get_err := GetMessage(conn, "file_ot_init")
		if get_err != nil {
			return get_err
		}

		ot_resp := &FileOTInitRespose{}
		err = json.Unmarshal(content, &ot_resp)
		if err != nil {
			return err
		}
		f.WriteString(ot_resp.Data.Buffer)
	} else {
		err := os.MkdirAll(local_path, os.ModeDir)
		if err != nil {
			return err
		}
		var req FSOPRequest
		req.Type = "fsop"
		req.Data.Type = "list_folder"
		req.Data.Param1 = web_path

		req_body, err := json.Marshal(req)
		if err != nil {
			return err
		}

		err = conn.WriteMessage(websocket.BinaryMessage, req_body)
		if err != nil {
			fmt.Printf("Write error\n")
			return err
		}

		content, get_err := GetMessage(conn, "list_reply")
		if get_err != nil {
			return get_err
		}

		m_resp := &ListingReply{}
		err = json.Unmarshal(content, &m_resp)
		if err != nil {
			return err
		}

		for _, returned := range m_resp.Data.Listing {
			err = RecReadPath(conn, fmt.Sprintf("%s/%s", web_path, returned.Name), path.Join(local_path, returned.Name), returned.Type != "file")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
