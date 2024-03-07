# HTTP-loader

This app use for load testing sites.

App have two versions of loader:
1) v1 wrote with recursion.
2) v2 wrote with workers.

App have two regime of work:
1) With proxy. App stop them self after use all proxies. Proxy skip after app get err.
2) Without proxy. You should stop app with ctrl+c. Only v2.

App can load with two layer HTTP and TCP.

In project must be created file with name socks5_proxies.txt with proxies inside. Or use -f and chose another file with
socks5 proxies in format on field on ip:port.

You can use -f or --help for get help with args.

Default values:

--maxProxyRead = [default: 100]

--rps  = [default: 100]

--version, -v = [default: 2]

--useproxy, -p = [default: false]

--file = [default: socks5_proxies.txt]

For example:
cmd.exe --targetHTTP "https://example.com/" --targetTCP "127.0.0.1:80" -r 5000 


List of site with proxies:
https://proxyscrape.com/free-proxy-list  
http://free-proxy.cz/en/proxylist/country/all/socks5/ping/all