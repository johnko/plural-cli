package pluralfile

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
	"github.com/pluralsh/plural/pkg/config"
	"github.com/pluralsh/plural/pkg/api"
)

type Lockfile struct {
	Artifact    map[string]string
	Terraform   map[string]string
	Helm        map[string]string
	Recipe      map[string]string
	Integration map[string]string
	Crd         map[string]string
	Ird         map[string]string
	Tag         map[string]string
	Attrs       map[string]string
}

func lock() *Lockfile {
	return &Lockfile{
		Artifact:    map[string]string{},
		Terraform:   map[string]string{},
		Helm:        map[string]string{},
		Recipe:      map[string]string{},
		Integration: map[string]string{},
		Crd:         map[string]string{},
		Ird:         map[string]string{},
		Tag:         map[string]string{},
		Attrs:       map[string]string{},
	}
}

func (plrl *Pluralfile) Lock(path string) (*Lockfile, error) {
	client := api.NewClient()
	applyLock, err := client.AcquireLock(plrl.Repo)
	if err != nil {
		return lock(), nil
	}

	if applyLock.Lock == "" {
		return Lock(path), nil
	}

	lock := lock()
	yaml.Unmarshal([]byte(applyLock.Lock), lock)
	return lock, nil
}

func (plrl *Pluralfile) Flush(path string, lock *Lockfile) error {
	client := api.NewClient()
	io, err := yaml.Marshal(lock)
	if err != nil {
		return err
	}

	_, err = client.ReleaseLock(plrl.Repo, string(io))
	return err
}

func Lock(path string) *Lockfile {
	conf := config.Read()
	lock := lock()
	lockfile := lockPath(path, conf.LockProfile)
	content, err := ioutil.ReadFile(lockfile)
	if err != nil {
		return lock
	}

	yaml.Unmarshal(content, lock)
	return lock
}

func lockPath(path string, profile string) string {
	if profile == "" {
		return filepath.Join(filepath.Dir(path), "plural.lock")
	}

	return filepath.Join(filepath.Dir(path), fmt.Sprintf("plural.%s.lock", profile))
}

func (lock *Lockfile) Flush(path string) error {
	conf := config.Read()
	io, err := yaml.Marshal(lock)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(lockPath(path, conf.LockProfile), io, 0644)
}

func (lock *Lockfile) getSha(name ComponentName, key string) string {
	switch name {
	case HELM:
		sha, _ := lock.Helm[key]
		return sha
	case TERRAFORM:
		sha, _ := lock.Terraform[key]
		return sha
	case RECIPE:
		sha, _ := lock.Recipe[key]
		return sha
	case ARTIFACT:
		sha, _ := lock.Artifact[key]
		return sha
	case INTEGRATION:
		sha, _ := lock.Integration[key]
		return sha
	case CRD:
		sha, _ := lock.Crd[key]
		return sha
	case IRD:
		sha, _ := lock.Ird[key]
		return sha
	case TAG:
		sha, _ := lock.Tag[key]
		return sha
	case REPO_ATTRS:
		sha, _ := lock.Attrs[key]
		return sha
	default:
		return ""
	}
}

func (lock *Lockfile) addSha(name ComponentName, key string, sha string) {
	switch name {
	case HELM:
		lock.Helm[key] = sha
		return
	case TERRAFORM:
		lock.Terraform[key] = sha
	case RECIPE:
		lock.Recipe[key] = sha
	case ARTIFACT:
		lock.Artifact[key] = sha
	case INTEGRATION:
		lock.Integration[key] = sha
	case CRD:
		lock.Crd[key] = sha
	case IRD:
		lock.Ird[key] = sha
	case TAG:
		lock.Tag[key] = sha
	case REPO_ATTRS:
		lock.Attrs[key] = sha
	default:
		return
	}
}
