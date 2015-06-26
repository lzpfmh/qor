package googlecloud

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/qor/qor/media_library"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

var (
	DefaultBucket   string
	DefaultEndpoint = "https://storage.googleapis.com"
)

type Storage struct {
	media_library.Base
}

var service *storage.Service

func init() {
	client, err := google.DefaultClient(context.Background(), storage.DevstorageFullControlScope)
	if err != nil {
		panic(err)
	}
	if service, err = storage.New(client); err != nil {
		panic(err)
	}
}

func getBucket(option *media_library.Option) string {
	if bucket := option.Get("bucket"); bucket != "" {
		return bucket
	}

	return DefaultBucket
}

func getEndpoint(option *media_library.Option) string {
	endpoint := option.Get("endpoint")
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}

	return endpoint + "/" + getBucket(option)
}

func (s Storage) GetURLTemplate(option *media_library.Option) (path string) {
	if path = option.Get("URL"); path == "" {
		path = "/{{class}}/{{primary_key}}/{{column}}/{{filename_with_hash}}"
	}

	return getEndpoint(option) + path
}

func (s Storage) Store(url string, option *media_library.Option, reader io.Reader) error {
	name := strings.Replace(url, getEndpoint(option), "", -1)
	object := &storage.Object{Name: name}

	bucket := getBucket(option)

	insert := service.Objects.Insert(bucket, object)
	insert = insert.PredefinedAcl("publicRead")

	_, err := insert.Media(reader).Do()
	return err
}

func (s Storage) Retrieve(url string) (*os.File, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	file, err := ioutil.TempFile("", "google-cloud")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(file, response.Body)
	return file, err
}
