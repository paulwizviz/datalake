// Package blockutil contains operations to extract raw protobuf from S3
// and local cache
package blockutil

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/paulwizviz/datalake/internal/block"
	"google.golang.org/protobuf/proto"
)

type s3BucketList struct {
	XMLName     xml.Name   `xml:"ListBucketResult"`
	Name        string     `xml:"Name"`
	Prefix      string     `xml:"Prefix"`
	Marker      string     `xml:"Marker"`
	MaxKeys     int        `xml:"MaxKeys"`
	IsTruncated bool       `xml:"IsTruncated"`
	Contents    []s3Object `xml:"Contents"`
}

type s3Object struct {
	Key          string `xml:"Key"`
	LastModified string `xml:"LastModified"`
	ETag         string `xml:"ETag"`
	Size         int    `xml:"Size"`
	Owner        struct {
		ID          string `xml:"ID"`
		DisplayName string `xml:"DisplayName"`
	} `xml:"Owner"`
	StorageClass string `xml:"StorageClass"`
}

var S3URL = "https://s3.us-east-1.amazonaws.com/public.blocks.datalake"

type ObjectKey string

// ReadS3ListURL read S3 xml list
func ReadS3ListURL(url string) ([]ObjectKey, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	list, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return extractListKeys(list)
}

func extractListKeys(list []byte) ([]ObjectKey, error) {
	var lbr s3BucketList
	err := xml.Unmarshal(list, &lbr)
	if err != nil {
		return nil, err
	}

	var keys []ObjectKey
	for _, c := range lbr.Contents {
		keys = append(keys, ObjectKey(c.Key))
	}

	return keys, nil
}

// ReadObjectByKey extract object after download from S3 container
func ReadObjectByKey(key ObjectKey) ([]byte, error) {
	objURL := S3URL + "/" + string(key)
	resp, err := http.Get(objURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// GetBlockNumber extract block number from raw protobuf
func GetBlockNumber(b []byte) string {
	var blk block.Block
	proto.Unmarshal(b, &blk)
	return blk.Number
}

func ReadBlockByHash(url string, key ObjectKey) (*block.Block, error) {
	objectURL := url + "/" + string(key) + ".datalake"

	resp, err := http.Get(objectURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var blk block.Block
	err = proto.Unmarshal(body, &blk)
	if err != nil {
		return nil, err
	}

	return &blk, nil
}

// ReadBlockByHashC extract block from a file stored in cache after download from S3
func ReadBlockByHashC(cache string, hash string) (*block.Block, error) {
	content, err := os.ReadFile(fmt.Sprintf("%s/%s.datalake.pb", cache, hash))
	if err != nil {
		return nil, err
	}
	var blk block.Block
	err = proto.Unmarshal(content, &blk)
	if err != nil {
		return nil, err
	}
	return &blk, nil
}

// ReadBlockByNumber extract a file named using block number from cache
func ReadBlockByNumber(cache string, num string) (*block.Block, error) {
	content, err := os.ReadFile(fmt.Sprintf("%s/%s.pb", cache, num))
	if err != nil {
		return nil, err
	}
	var blk block.Block
	err = proto.Unmarshal(content, &blk)
	if err != nil {
		return nil, err
	}
	return &blk, nil
}
