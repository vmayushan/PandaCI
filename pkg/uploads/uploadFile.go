package uploads

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func UploadFile(ctx context.Context, url string, file io.Reader, mime string) error {
	req, err := http.NewRequest("PUT", url, file)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", mime)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code from upload req: %d", resp.StatusCode)
	}

	return nil
}
