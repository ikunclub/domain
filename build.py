import urllib.request
import re
import ipaddress

URL = "https://raw.githubusercontent.com/marcuccilli/gary/refs/heads/main/iptv_passwall.txt"

def is_ip_or_cidr(val):
    val = val.strip()
    try:
        ipaddress.ip_network(val, strict=False)
        return True
    except ValueError:
        return False

def clean_domain(val):
    val = val.strip()
    val = val.split('#')[0].strip()
    val = re.sub(r'^https?://', '', val)
    val = val.split('/')[0].split(':')[0]
    return val

def main():
    req = urllib.request.Request(URL, headers={'User-Agent': 'Mozilla/5.0'})
    with urllib.request.urlopen(req) as resp:
        content = resp.read().decode('utf-8')

    domains = set()
    ips = set()

    for line in content.splitlines():
        line = line.strip()
        if not line or line.startswith('#'):
            continue

        if line.startswith("domain:") or line.startswith("full:"):
            line = line.split(":", 1)[1]
        elif line.startswith("ip:"):
            line = line.split(":", 1)[1]

        if is_ip_or_cidr(line):
            ips.add(line)
        else:
            domain = clean_domain(line)
            if domain:
                domains.add(domain)

    with open("domains_iptv.txt", "w", encoding="utf-8") as f:
        for d in sorted(domains):
            f.write(f"{d}\n")

    with open("ips_iptv.txt", "w", encoding="utf-8") as f:
        for ip in sorted(ips):
            f.write(f"{ip}\n")

    print(f"解析完成：提取到 {len(domains)} 个域名，{len(ips)} 个 IP/网段。")

if __name__ == "__main__":
    main()