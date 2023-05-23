package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net"
    "os"
    "strings"
    "regexp"
    "github.com/miekg/dns"
    "github.com/spf13/cobra"
    "github.com/valyala/fasthttp"
)

type BGPViewData struct {
    Data struct {
        Prefixes []struct {
            Prefix string `json:"prefix"`
            ASN    struct {
                ASN int `json:"asn"`
            } `json:"asn"`
        } `json:"prefixes"`
    } `json:"data"`
}

func isDomain(input string) bool {
    domainPattern := `^([a-zA-Z0-9-]+\.){1,}[a-zA-Z]{2,}$`
    match, _ := regexp.MatchString(domainPattern, input)
    return match
}

func isIPAddress(input string) bool {
    return net.ParseIP(input) != nil
}

func resolveDomain(domain string) ([]string, error) {
    var IPs []string
    m := new(dns.Msg)
    m.SetQuestion(domain+".", dns.TypeA)
    r, err := dns.Exchange(m, "8.8.8.8:53")
    if err != nil {
        return nil, err
    }

    for _, answer := range r.Answer {
        if a, ok := answer.(*dns.A); ok {
            IPs = append(IPs, a.A.String())
        }
    }
    return IPs, nil
}

func fetchBGPViewData(ip string) (*BGPViewData, error) {
    url := fmt.Sprintf("https://api.bgpview.io/ip/%s", ip)
    statusCode, resp, err := fasthttp.Get(nil, url)
    if err != nil {
        return nil, err
    }

    if statusCode != fasthttp.StatusOK {
        return nil, fmt.Errorf("HTTP request failed with status code: %d", statusCode)
    }

    var data BGPViewData
    err = json.Unmarshal(resp, &data)
    if err != nil {
        return nil, err
    }

    return &data, nil
}

func processBGPViewIP(ips []string, asn bool, prefix bool) {
    for _, ip := range ips {
        fmt.Printf("[+] %s:\n", ip)
        data, err := fetchBGPViewData(ip)
        if err != nil {
            fmt.Println(err)
            continue
        }
        for _, prefixData := range data.Data.Prefixes {
            if asn {
                fmt.Println(prefixData.ASN.ASN)
            } else if prefix {
                fmt.Println(prefixData.Prefix)
            } else {
                fmt.Printf("prefix: %s, ASN: %d\n", prefixData.Prefix, prefixData.ASN.ASN)
            }
        }
    }
}

func main() {
    var url, list string
    var asn, prefix bool

    cmd := &cobra.Command{
        Use:   "bgpview",
        Short: "Process IPs or subdomains using BGPView API",
        Run: func(cmd *cobra.Command, args []string) {
    var urls []string
    if url != "" {
        urls = append(urls, url)
    } else if list != "" {
        content, err := ioutil.ReadFile(list)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        urls = strings.Split(strings.TrimSpace(string(content)), "\n")
    }

    for _, u := range urls {
        u := strings.TrimSpace(u) // Trim any leading/trailing whitespace
        if u == "" {
            continue // Skip empty lines
        }

        if isDomain(u) {
            IPs, err := resolveDomain(u)
            if err != nil {
                fmt.Println(err)
                continue
            }
            processBGPViewIP(IPs, asn, prefix)
        } else if isIPAddress(u) {
            IPs := []string{u}
            processBGPViewIP(IPs, asn, prefix)
        } else {
            fmt.Printf("[!] The '%s' is not my type :)\n", u)
            continue
        }
    }
},
    }

    cmd.Flags().StringVarP(&url, "url", "u", "", "Single URL to check")
    cmd.Flags().StringVarP(&list, "list", "l", "", "List of URLs to check")
    cmd.Flags().BoolVarP(&asn, "asn", "a", false, "Only print ASNs of target IP")
    cmd.Flags().BoolVarP(&prefix, "prefix", "p", false, "Only print prefixes of target IP")

    if err := cmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
