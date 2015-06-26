package rackspace

import (
	"fmt"
	"io"

	"github.com/qor/qor/media_library"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/rackspace"
	"github.com/rackspace/gophercloud/rackspace/objectstorage/v1/objects"
)

type Rackspace struct {
	media_library.Base
}

func (f Rackspace) Store(url string, option *media_library.Option, reader io.Reader) error {
	ao := gophercloud.AuthOptions{
		Username: "bodhiphilpot",
		APIKey:   "e9b2027a267e48628b9dae0b66720824",
	}
	provider, _ := rackspace.AuthenticatedClient(ao)

	serviceClient, errE := rackspace.NewObjectStorageV1(provider, gophercloud.EndpointOpts{
		Region: "HKG",
	})

	if errE != nil {
		panic(errE)
	}

	// readSeeker := reader.(multipart.File)
	readSeeker := reader.(io.ReadSeeker)
	err := objects.Create(
		serviceClient,
		"qor",
		"test1.png",
		readSeeker,
		nil,
	)
	fmt.Println(err)

	return errE
}

// func (f Rackspace) Retrieve(url string) (*os.File, error) {
// 	// if fullpath, err := f.GetFullPath(url, nil); err == nil {
// 	// 	return os.Open(fullpath)
// 	// } else {
// 	// 	return nil, os.ErrNotExist
// 	// }
// 	return nil, nil
// }
