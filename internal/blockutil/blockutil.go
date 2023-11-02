package blockutil

import (
	"encoding/xml"
	"io"
	"net/http"

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

func ReadBlockByKey(url string, key ObjectKey) (*block.Block, error) {
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
