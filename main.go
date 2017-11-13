package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	domainPtr := flag.String("domain", os.Getenv("CERTBOT_DOMAIN"), "The domain being authenticated")
	validationPtr := flag.String("validation", os.Getenv("CERTBOT_VALIDATION"), "The validation string")
	configFilePtr := flag.String("config", "config.json", "The config file location")
	debugPtr := flag.Bool("debug", false, "Specify to include debug logging")

	flag.Parse()

	domain := *domainPtr
	validation := *validationPtr

	if domain == "" || validation == "" {
		flag.Usage()

		os.Exit(1)
	}

	configFile, err := ioutil.ReadFile(*configFilePtr)
	if err != nil {
		log.Fatal(err)
	}

	config := &config{}

	if err := json.Unmarshal(configFile, config); err != nil {
		log.Fatal(err)
	}

	domainCfg, ok := config.Domains[domain]
	if !ok {
		log.Fatalf("domain %s is not in config file", domain)
	}

	if domainCfg.AcmeDNSURL == "" {
		domainCfg.AcmeDNSURL = config.AcmeDNSURL
	}

	requestUrl, err := resolveURL(domainCfg.AcmeDNSURL, "update")
	if err != nil {
		log.Fatal(err)
	}

	body := &updateBody{
		Subdomain: domainCfg.Subdomain,
		Txt:       validation,
	}

	data, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, requestUrl.String(), bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("X-Api-User", domainCfg.Username)
	req.Header.Add("X-Api-Key", domainCfg.Password)

	if *debugPtr {
		fmt.Printf("Executing %+v\n", req)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		log.Fatalf("response status code was not 200: %d: %s", res.StatusCode, string(body))
	}

	if *debugPtr {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Response received, unable to read body: %v\n", err)
		} else {
			fmt.Printf("Response received: %s\n", string(body))
		}
	}

	propagationDuration := domainCfg.PropagationDuration
	if propagationDuration == 0 {
		propagationDuration = config.PropagationDuration
	}

	if *debugPtr {
		fmt.Printf("Waiting for %v for DNS to propagate\n", time.Duration(propagationDuration))
	}

	time.Sleep(time.Duration(propagationDuration))

	if *debugPtr {
		fmt.Println("Done")
	}
}

func resolveURL(base, path string) (*url.URL, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return nil, err
	}

	pathURL, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	return baseURL.ResolveReference(pathURL), nil
}
