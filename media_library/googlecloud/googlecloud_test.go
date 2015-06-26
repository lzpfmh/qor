package googlecloud

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/qor/qor/test/utils"
)

var (
	db       = utils.TestDB()
	testFile *os.File
)

type Product struct {
	gorm.Model
	Image Storage
}

func init() {
	DefaultBucket = os.Getenv("GOOGLE_CLOUD_BUCKET")
	if DefaultBucket == "" {
		println("please specify default bucket by GOOGLE_CLOUD_BUCKET=name go test")
		os.Exit(1)
	}

	db.AutoMigrate(&Product{})

	var err error
	if testFile, err = ioutil.TempFile("", "google-cloud"); err != nil {
		panic(err)
	}
	testFile.WriteString("i am a test file")
	if err = testFile.Close(); err != nil {
		panic(err)
	}
	if testFile, err = os.Open(testFile.Name()); err != nil {
		panic(err)
	}
}

func TestStoreAndRetrieve(t *testing.T) {
	var product Product

	if err := product.Image.Scan(testFile); err != nil {
		t.Error(err)
	}

	if err := db.Save(&product).Error; err != nil {
		t.Error(err)
	}

	resp, err := http.Get(product.Image.URL())
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("Status code is not 200, is %+v", resp.Status)
	}
}
