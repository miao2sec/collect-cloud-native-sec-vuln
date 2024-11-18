package component

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v66/github"
	"golang.org/x/xerrors"
	"os"
	"path/filepath"
	"strconv"
)

type Component struct {
	Owner      string                     `yaml:"owner,omitempty"`
	Repo       string                     `yaml:"repo,omitempty"`
	Advisories []*github.SecurityAdvisory `yaml:"advisories,omitempty"`
}

func (c *Component) Save(dir string) error {
	for _, advisory := range c.Advisories {
		dir := filepath.Join(dir, c.Repo,
			strconv.Itoa(advisory.PublishedAt.Year()),
			fmt.Sprintf("%02d", int(advisory.PublishedAt.Month())),
		)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return xerrors.Errorf("failed to mkdir %s:%w", dir, err)
		}

		data, err := json.MarshalIndent(advisory, "", "\t")
		if err != nil {
			return xerrors.Errorf("failed to marshal %s:%w", *advisory.GHSAID, err)
		}
		var vulnName string
		if advisory.CVEID != nil {
			vulnName = *advisory.CVEID
		} else {
			vulnName = *advisory.GHSAID
		}
		err = os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s.json", vulnName)), data, os.ModePerm)
		if err != nil {
			return xerrors.Errorf("failed to write %s:%w", *advisory.GHSAID, err)
		}
	}
	return nil
}
