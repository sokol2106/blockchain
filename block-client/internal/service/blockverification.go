package service

import (
	"encoding/json"
	"fmt"
	"github.com/ivan/blockchain/block-client/internal/model"
	"io"
	"net/http"
)

type Blockverification struct {
	block   model.VerificationBlock
	address string
}

func NewBlockVerification(a string) *Blockverification {
	return &Blockverification{
		block:   model.VerificationBlock{},
		address: a,
	}
}

func (b *Blockverification) RequestVerificationData() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", b.address+"/api/blockchain/block/verify", nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error bad status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	data := model.VerificationBlock{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	b.block = data
	return nil
}

func (b *Blockverification) VerifyData() {

}
