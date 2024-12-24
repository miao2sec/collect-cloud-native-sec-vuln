package utils

import (
	"encoding/json"
	"fmt"
	"golang.org/x/xerrors"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	URL   = "https://kubernetes.io/docs/reference/issues-security/official-cve-feed/index.json"
	Retry = 5
)

type Kubernetes struct {
	KubernetesIo KubernetesIo `json:"_kubernetes_io,omitempty"`
	Authors      []*Author    `json:"authors,omitempty"`
	Description  string       `json:"description,omitempty"`
	FeedUrl      string       `json:"feed_url,omitempty"`
	HomePageUrl  string       `json:"home_page_url,omitempty"`
	Items        []*Itme      `json:"items,omitempty"`
	Title        string       `json:"title,omitempty"`
	Version      string       `json:"version,omitempty"`
}

type Author struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type Itme struct {
	KubernetesIo  KubernetesIo `json:"_kubernetes_io,omitempty"`
	ContentText   string       `json:"content_text,omitempty"`
	DatePublished time.Time    `json:"date_published,omitempty"`
	ExternalUrl   string       `json:"external_url,omitempty"`
	Id            string       `json:"id,omitempty"`
	Status        string       `json:"status,omitempty"`
	Summary       string       `json:"summary,omitempty"`
	Url           string       `json:"url,omitempty"`
}

type KubernetesIo struct {
	FeedRefreshJob string    `json:"feed_refresh_job,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	GoogleGroupUrl string    `json:"google_group_url,omitempty"`
	IssueNumber    int       `json:"issue_number,omitempty"`
}

func NewKubernetes() (*Kubernetes, error) {
	var (
		http       = HTTP{URL: URL, Method: Get, Retry: Retry}
		kubernetes Kubernetes
	)

	if err := http.Fetch(&kubernetes); err != nil {
		return nil, err

	}
	return &kubernetes, nil
}

func (k *Kubernetes) Save(cacheDir string) (err error) {
	var (
		dir  string
		data []byte
	)

	for _, item := range k.Items {
		dir = filepath.Join(cacheDir, "kubernetes",
			strconv.Itoa(item.DatePublished.Year()),
			fmt.Sprintf("%02d", int(item.DatePublished.Month())),
		)
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return xerrors.Errorf("failed to mkdir %s:%w", dir, err)
		}
		data, err = json.MarshalIndent(item, "", "\t")
		if err != nil {
			return xerrors.Errorf("failed to marshal %s:%w", item.Id, err)
		}
		err = os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.json", item.Id)), data, os.ModePerm)
		if err != nil {
			return xerrors.Errorf("failed to write %s:%w", item.Id, err)
		}
	}
	return nil
}
