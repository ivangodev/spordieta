package core

import (
	"bytes"
	"fmt"
	"github.com/ivangodev/spordieta/entity"
	"io/ioutil"
	"net/http"
)

type storageConnect struct {
	domain string
}

func newStrgConn(domain string) *storageConnect {
	return &storageConnect{domain}
}

func (c *storageConnect) get(url string) (ret string, err error) {
	var resp *http.Response
	resp, err = http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	ret = string(body)
	return
}

func (c *storageConnect) delete(url string) (err error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *storageConnect) put(url, content string) (err error) {
	b := bytes.NewBuffer([]byte(content))
	req, err := http.NewRequest("PUT", url, b)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (c *storageConnect) GetWinProofURL(userId entity.UserId, betId entity.BetId) (url string, err error) {
	winUrl := c.domain + fmt.Sprintf("/user/%v/bet/%v/winproof", userId, betId)
	url, err = c.get(winUrl)
	return
}

func (c *storageConnect) GetPayProofURL(userId entity.UserId, betId entity.BetId) (url string, err error) {
	payUrl := c.domain + fmt.Sprintf("/user/%v/bet/%v/payproof", userId, betId)
	url, err = c.get(payUrl)
	return
}

func (c *storageConnect) DeleteProofs(userId entity.UserId, betId entity.BetId) error {
	winUrl := c.domain + fmt.Sprintf("/user/%v/bet/%v/winproof", userId, betId)
	payUrl := c.domain + fmt.Sprintf("/user/%v/bet/%v/payproof", userId, betId)

	err := c.delete(winUrl)
	if err != nil {
		return err
	}
	c.delete(payUrl)
	if err != nil {
		return err
	}
	return nil
}

//Useful for tests:
func (c *storageConnect) uploadWinProof(userId entity.UserId, betId entity.BetId, s string) error {
	winUrl := c.domain + fmt.Sprintf("/user/%v/bet/%v/winproof", userId, betId)
	return c.put(winUrl, s)
}

func (c *storageConnect) uploadPayProof(userId entity.UserId, betId entity.BetId, s string) error {
	payUrl := c.domain + fmt.Sprintf("/user/%v/bet/%v/payproof", userId, betId)
	return c.put(payUrl, s)
}
