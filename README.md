_What is `d8`?_

`d8` is a DNS crawler library written in Go. It is also a DNS crawling utility.
The crawler is for mapping out and tracking DNS infrastructures used by a set of 
particular domains. In specific, it take a domain as input, and gives back the cname
redirection chain, all the ip records, the non-registry name servers that supports
the domain resolving process, and all the records (A, CNAME, NS, SOA, TXT, MX) that
an Internet structure analytic might have interest.

_Is it a DNS client or a DNS server?_

It is neither. It implements a simple DNS client that can parse several types of DNS
records, but it is not targeted to be a full DNS client.

_Does it support IPv6?_

No. `d8` is IPv4 only.

_Example Run_

	$ make
	$ ./d8 www.yahoo.com
	// www.yahoo.com
	cnames {
	    www.yahoo.com -> fd-fp3.wg1.b.yahoo.com
	    fd-fp3.wg1.b.yahoo.com -> ds-fp3.wg1.b.yahoo.com
	    ds-fp3.wg1.b.yahoo.com -> ds-any-fp3-lfb.wa1.b.yahoo.com
	    ds-any-fp3-lfb.wa1.b.yahoo.com -> ds-any-fp3-real.wa1.b.yahoo.com
	}
	ips {
	    98.138.253.109(ds-any-fp3-real.wa1.b.yahoo.com)
	    98.138.252.30(ds-any-fp3-real.wa1.b.yahoo.com)
	    206.190.36.105(ds-any-fp3-real.wa1.b.yahoo.com)
	    206.190.36.45(ds-any-fp3-real.wa1.b.yahoo.com)
	}
	servers {
	    yahoo.com ns ns1.yahoo.com(68.180.131.16)
	    yahoo.com ns ns5.yahoo.com(119.160.247.124)
	    yahoo.com ns ns2.yahoo.com(68.142.255.16)
	    yahoo.com ns ns3.yahoo.com(203.84.221.53)
	    yahoo.com ns ns4.yahoo.com(98.138.11.157)
	    wg1.b.yahoo.com ns yf1.yahoo.com(68.142.254.15)
	    wg1.b.yahoo.com ns yf2.yahoo.com(68.180.130.15)
	}
	records {
	    www.yahoo.com cname fd-fp3.wg1.b.yahoo.com
	    fd-fp3.wg1.b.yahoo.com cname ds-fp3.wg1.b.yahoo.com
	    ds-fp3.wg1.b.yahoo.com cname ds-any-fp3-lfb.wa1.b.yahoo.com
	    ds-any-fp3-lfb.wa1.b.yahoo.com cname ds-any-fp3-real.wa1.b.yahoo.com
	    ds-any-fp3-real.wa1.b.yahoo.com a 98.138.253.109
	    ds-any-fp3-real.wa1.b.yahoo.com a 98.138.252.30
	    ds-any-fp3-real.wa1.b.yahoo.com a 206.190.36.105
	    ds-any-fp3-real.wa1.b.yahoo.com a 206.190.36.45
	    yahoo.com ns ns1.yahoo.com
	    yahoo.com ns ns5.yahoo.com
	    yahoo.com ns ns2.yahoo.com
	    yahoo.com ns ns3.yahoo.com
	    yahoo.com ns ns4.yahoo.com
	    ns1.yahoo.com a 68.180.131.16
	    ns5.yahoo.com a 119.160.247.124
	    ns2.yahoo.com a 68.142.255.16
	    ns3.yahoo.com a 203.84.221.53
	    ns4.yahoo.com a 98.138.11.157
	    wg1.b.yahoo.com ns yf3.a1.b.yahoo.net
	    wg1.b.yahoo.com ns yf1.yahoo.com
	    wg1.b.yahoo.com ns yf2.yahoo.com
	    wg1.b.yahoo.com ns yf4.a1.b.yahoo.net
	    yf1.yahoo.com a 68.142.254.15
	    yf2.yahoo.com a 68.180.130.15
	    yahoo.com ns ns8.yahoo.com
	    yahoo.com ns ns6.yahoo.com
	    yahoo.com mx mta6.am0.yahoodns.net/1
	    yahoo.com mx mta5.am0.yahoodns.net/1
	    yahoo.com mx mta7.am0.yahoodns.net/1
	    yahoo.com soa ns1.yahoo.com/hostmaster.yahoo-inc.com serial=2013123014 refresh=3600 retry=300 exp=1814400 min=600
	    wg1.b.yahoo.com soa yf1.yahoo.com/hostmaster.yahoo-inc.com serial=1388446812 refresh=30 retry=30 exp=86400 min=300
	}
