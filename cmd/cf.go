package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var cfCmd = &cobra.Command{
	Use:   "cf",
	Short: "cloudflare DDNS update",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		req, err := http.NewRequest("GET",
			fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?type=%s&name=%s&page=1&per_page=20&order=type&direction=desc&match=all",
				zone, dnsType, domain), nil)
		if err != nil {
			log.Fatal("failed", err)
		}
		req.Header.Set("X-Auth-Email", "trungtvq.work@gmail.com")
		req.Header.Set("X-Auth-Key", "2410c88d19996ecb35eacb75544c85f2c7f05")
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal("failed", err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("failed", err)
		}
		fmt.Println("body", string(body))

		var dnsRespValue dnsResp
		if err = json.Unmarshal(body, &dnsRespValue); err != nil {
			log.Fatal("failed", err)
		}
		if len(dnsRespValue.Result) < 1 {
			createNew()
		} else {
			update(dnsRespValue.Result[0].Id)
		}
		os.Exit(0)
	},
}

func update(dnsID string) {
	fmt.Println("update")
	type Payload struct {
		Type    string `json:"type"`
		Name    string `json:"name"`
		Content string `json:"content"`
		TTL     int    `json:"ttl"`
		Proxied bool   `json:"proxied"`
	}

	data := Payload{
		Type:    dnsType,
		Name:    domain,
		Content: value,
		TTL:     120,
		Proxied: proxy,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal("failed", err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("PUT", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zone, dnsID), body)
	if err != nil {
		log.Fatal("failed", err)
	}
	req.Header.Set("X-Auth-Email", "trungtvq.work@gmail.com")
	req.Header.Set("X-Auth-Key", "2408c88d19996ecb35eacb75544c85f2c7f05")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("failed", err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("failed", err)
	}
	resp.Body.Close()
	fmt.Println(string(respBody))
}
func createNew() {
	type Payload struct {
		Type     string `json:"type"`
		Name     string `json:"name"`
		Content  string `json:"content"`
		TTL      int    `json:"ttl"`
		Priority int    `json:"priority"`
		Proxied  bool   `json:"proxied"`
	}

	data := Payload{
		Type:     dnsType,
		Name:     domain,
		Content:  value,
		TTL:      120,
		Priority: 10,
		Proxied:  proxy,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal("failed", err)
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zone), body)
	if err != nil {
		log.Fatal("failed", err)
	}
	req.Header.Set("X-Auth-Email", "trungtvq.work@gmail.com")
	req.Header.Set("X-Auth-Key", "2410c88d19996ecb35eacb75544c85f2c7f05")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("failed", err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("failed", err)
	}
	resp.Body.Close()
	fmt.Println(string(respBody))
}

type dnsResp struct {
	Result []struct {
		Id string `json:"id"`
	} `json:"result"`
}

var (
	domain  string
	dnsType string
	value   string
	zone    string
	proxy   bool
)

func init() {
	cfCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "", "domain name to update")
	cfCmd.PersistentFlags().StringVarP(&dnsType, "dnsType", "t", "", "dns type to update")
	cfCmd.PersistentFlags().StringVarP(&value, "value", "v", "", "dns value to update")
	cfCmd.PersistentFlags().BoolVarP(&proxy, "proxy", "p", false, "use proxy")
	cfCmd.PersistentFlags().StringVarP(&zone, "zone", "z", "8ecbe52e489a1d34bf7933b5ea9623c0", "zone to update (default zone of tool2.xyz)")
	rootCmd.AddCommand(cfCmd)
}
