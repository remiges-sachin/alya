package config

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/remiges-tech/rigel"
	"github.com/remiges-tech/rigel/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// ConfigSource is an interface that represents a source from which application configuration can be loaded.
type Config interface {
	// LoadConfig reads config from the source and binds it to c.
	LoadConfig(c any) error
	// Check checks if the config source can be used. For example, a file config source would check
	// if the file exists. A Rigel config source would check if the Rigel server is available.
	Check() error
}

// Load first ensures that the config system valid and accessible. Then it loads the config into c.
func Load(cs Config, c any) error {
	if err := cs.Check(); err != nil {
		return err
	}
	return cs.LoadConfig(c)
}

// File

type File struct {
	ConfigFilePath string
}

func (f *File) Check() error {
	if f.ConfigFilePath == "" {
		return fmt.Errorf("configFilePath cannot be empty")
	}

	return nil
}

func newFile(configFilePath string) (*File, error) {
	file := &File{ConfigFilePath: configFilePath}

	if err := file.Check(); err != nil {
		return nil, err
	}

	return file, nil
}

func (f *File) LoadConfig(appConfig any) error {
	filePath := f.ConfigFilePath
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(appConfig)
}

// Rigel
type Rigel struct {
	Client        *rigel.Rigel
	SchemaName    string
	SchemaVersion int
	ConfigName    string
}

func (r *Rigel) LoadConfig(config any) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return r.Client.LoadConfig(ctx, r.SchemaName, r.SchemaVersion, r.ConfigName, config)
}

func NewRigelClient(etcdEndpoints string) (*rigel.Rigel, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdEndpoints},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create etcd client: %v", err)
		return nil, err
	}

	etcdStorage := &etcd.EtcdStorage{Client: cli}
	rigelClient := rigel.New(etcdStorage)

	return rigelClient, nil
}
