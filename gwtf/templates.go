package gwtf

import (
	"bufio"
	"encoding/base64"
	"io/ioutil"
	"math"
	"net/http"
	"os"
)

func splitByWidthMake(str string, size int) []string {
	strLength := len(str)
	splitedLength := int(math.Ceil(float64(strLength) / float64(size)))
	splited := make([]string, splitedLength)
	var start, stop int
	for i := 0; i < splitedLength; i += 1 {
		start = i * size
		stop = start + size
		if stop > strLength {
			stop = strLength
		}
		splited[i] = str[start:stop]
	}
	return splited
}

func fileAsImageData(path string) (string, error) {
	f, _ := os.Open(path)

	defer f.Close()

	// Read entire JPG into byte slice.
	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(content)

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)

	return "data:" + contentType + ";base64, " + encoded, nil
}

func fileAsBase64(path string) (string, error) {
	f, _ := os.Open(path)

	defer f.Close()

	// Read entire JPG into byte slice.
	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)

	return encoded, nil
}

//UploadFile reads a file, base64 encodes it and chunk upload to /storage/upload
func (f *GoWithTheFlow) UploadFile(filename string, accountName string) error {
	content, err := fileAsBase64(filename)
	if err != nil {
		return err
	}

	return f.UploadString(content, accountName)
}

//UploadImageAsDataUrl will upload a image file from the filesystem into /storage/upload of the given account
func (f *GoWithTheFlow) UploadImageAsDataUrl(filename string, accountName string) error {
	content, err := fileAsImageData(filename)
	if err != nil {
		return err
	}

	return f.UploadString(content, accountName)
}

//UploadStrgin will upload the given string data in 1mb chunkts to /storage/upload of the given account
func (f *GoWithTheFlow) UploadString(content string, accountName string) error {

	//unload previous content if any.
	_, err := f.Transaction(`
transaction(part: String) {
    prepare(signer: AuthAccount) {
        let path = /storage/upload
        let existing = signer.load<String>(from: path) ?? ""
				log(existing)
    }
}

  `).SignProposeAndPayAs(accountName).RunE()
	if err != nil {
		return err
	}

	parts := splitByWidthMake(content, 1_000_000)
	for _, part := range parts {
		_, err := f.Transaction(`
transaction(part: String) {
    prepare(signer: AuthAccount) {
        let path = /storage/upload
        let existing = signer.load<String>(from: path) ?? ""
        signer.save(existing.concat(part), to: path)
    }
}
    `).SignProposeAndPayAs(accountName).StringArgument(part).SignProposeAndPayAsService().RunE()
		if err != nil {
			return err
		}
	}
	return nil
}
