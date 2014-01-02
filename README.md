**What is `d8`?**

`d8` is a DNS crawler library written in Go. It is also comes with a simple DNS
crawling program that uses the library.  The crawler implements the logic that
automatically maps out the DNS infrastructures used by a set of particular
domains by recursively crawling that starts from the root domain name servers.
In specific, it take a domain as input, and gives back the cname redirection
chain, all the ip records, the non-registry name servers that supports the
domain resolving process, and all the records (A, CNAME, NS, SOA, TXT, MX) that
an Internet structure analytic might have interest.

**Is it a DNS client or server?**

It is neither. It implements a very simple DNS client that can parse very
limited types of DNS records, but it is not targeted to be a full client. It is
definitely not a DNS server.

**Any dependencies required?**

It only depends on Go language and its standard library.

**How to build?**

In short, type `make` under the root repo directory.

Long version. As the author, I am sorry that I do not follow the standard way
of using *remote import path* that most Go libraries and programs used.  I feel
`import "github.com/h8liu/d8"` is just too long and putting several these lines
at the head of each `.go` file is just so ugly. So Instead, this repository
serves as a standalone `GOPATH` directory. For convenience, it contains a
simple `makefile` that simply wraps some shortcuts. For example, run `make`
under the root of this repository will perform a `go install` with `GOPATH` set
to the current working directory. 

**Any documentation?**

Do not really have one now, since I don't see it will become popular in near
future.  If you need one, please tell me and I might consider prioritize that.
Anyway, you can always use `make doc` to browse the package APIs.

**Does it support IPv6?**

No. `d8` is IPv4 only.

**Does it cache records?**

Only NS records and their glued A records of registry name servers (like TLD
servers of `com`, `net`, `com.ru`, `org.cn`, etc.) are cached by default. All
other info are crawled fresh from the Internet.

**Library Structure**

The core library:

- `d8/domain` provides domain name parsing.
- `d8/client` provides a simple DNS async client.
- `d8/packet` provides DNS packet parsing (for crawling purposes).
- `d8/packet/consts` defines rdata type and class codes.
- `d8/packet/rdata` provides DNS records parsing (for crawling purposes).
- `d8/term` provides a recursive crawling cursor for executing crawling
  logics.
- `d8/tasks` implements several common crawling logics.

General purpose helpers: 

- `printer` - Provides a simple line printer that supports indenting. Used for
  printing logs.
- `subcmd` - Provides APIs for defining sub commands.

Binaries:

- `bin/d8` - Implements a utility program that wraps around `d8` library. It
  provides interactive crawling and batch crawling.

**TODO List**

- `d8/client` better truncated and invalid message parsing and handling
  (instead of just log and discard)
- `d8/tasks` mark temporarily unreachable name servers in ns caches
- `d8/term` retry fallback on service rejected (return code as refused)
- `bin/d8` implement `.recur` and `.dig`
- `bin/d8` provide output as SQL dumps
- implement crawling session for tracking changes
- implement crawling groups for crawling wildcard domains

**Example Run**

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
