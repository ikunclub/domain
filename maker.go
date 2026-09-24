package main

import (
	"bufio"
	"fmt"
	"net/netip"
	"os"
	"strings"

	"github.com/v2fly/v2ray-core/v4/app/router"
	"google.golang.org/protobuf/proto"
)

func main() {
	buildGeoIP()
	buildGeoSite()
}

func buildGeoIP() {
	file, err := os.Open("ips_iptv.txt")
	if err != nil {
		fmt.Println("打开 ips_iptv.txt 失败:", err)
		return
	}
	defer file.Close()

	var cidrs []*router.CIDR
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.Contains(line, "/") {
			if strings.Contains(line, ":") {
				line += "/128"
			} else {
				line += "/32"
			}
		}
		prefix, err := netip.ParsePrefix(line)
		if err != nil {
			fmt.Printf("跳过无效 IP/CIDR: %s\n", line)
			continue
		}
		addr := prefix.Addr()
		cidrs = append(cidrs, &router.CIDR{
			Ip:     addr.AsSlice(),
			Prefix: uint32(prefix.Bits()),
		})
	}

	geoipList := &router.GeoIPList{
		Entry: []*router.GeoIP{
			{
				CountryCode: "IPTV",
				Cidr:        cidrs,
			},
		},
	}

	data, err := proto.Marshal(geoipList)
	if err != nil {
		panic(err)
	}

	_ = os.WriteFile("geoip-iptv.dat", data, 0644)
	_ = os.WriteFile("geoip.dat", data, 0644)
	fmt.Printf("GeoIP 编译完成：写入 %d 个 CIDR 条目\n", len(cidrs))
}

func buildGeoSite() {
	file, err := os.Open("domains_iptv.txt")
	if err != nil {
		fmt.Println("打开 domains_iptv.txt 失败:", err)
		return
	}
	defer file.Close()

	var domains []*router.Domain
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		domains = append(domains, &router.Domain{
			Type:  router.Domain_Domain,
			Value: line,
		})
	}

	geositeList := &router.GeoSiteList{
		Entry: []*router.GeoSite{
			{
				CountryCode: "IPTV",
				Domain:      domains,
			},
		},
	}

	data, err := proto.Marshal(geositeList)
	if err != nil {
		panic(err)
	}

	_ = os.WriteFile("geosite-iptv.dat", data, 0644)
	_ = os.WriteFile("geosite.dat", data, 0644)
	fmt.Printf("GeoSite 编译完成：写入 %d 个域名条目\n", len(domains))
}
