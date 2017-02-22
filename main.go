package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
)

var (
	tpl_id      = "your yunpian sms template id"
	url_tpl_sms = "https://sms.yunpian.com/v2/sms/tpl_single_send.json"
	apikey      = "your yunpian api key"
	mobile      = "your phone number"

	address              = "your monitor address"
	etherscanApiEndpoint = "https://api.etherscan.io/api"
	oneHundredEtherInWei = common.String2Big("100000000000000000000")
)

type accountBalance struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func getBalance(address string, target interface{}) error {
	baseUrl, err := url.Parse(etherscanApiEndpoint)
	if err != nil {
		return err
	}

	params := url.Values{}
	params.Add("module", "account")
	params.Add("action", "balance")
	params.Add("address", address)

	baseUrl.RawQuery = params.Encode()

	resp, err := http.Get(baseUrl.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func sendNotification(content string) error {
	tpl_value := url.Values{"#content#": {fmt.Sprintf("%v", content)}}.Encode()

	data := url.Values{
		"apikey":    {apikey},
		"mobile":    {mobile},
		"tpl_id":    {tpl_id},
		"tpl_value": {tpl_value},
	}

	resp, err := http.PostForm(url_tpl_sms, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(string(body))

	return nil
}

func main() {
	balance := &accountBalance{}
	err := getBalance(address, balance)
	if err != nil {
		panic(err)
	}

	etherInWei := common.String2Big(balance.Result)
	if etherInWei.Cmp(oneHundredEtherInWei) > 0 {
		return
	}

	err := sendNotification(common.CurrencyToString(etherInWei))
	if err != nil {
		panic(err)
	}
}
