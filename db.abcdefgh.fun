$ORIGIN abcdefgh.fun.
@	3600 IN	SOA dns.abcdefgh.fun. dns.abcdefgh.fun. (
				2017042745 ; serial
				7200       ; refresh (2 hours)
				3600       ; retry (1 hour)
				1209600    ; expire (2 weeks)
				3600       ; minimum (1 hour)
				)

	3600 IN NS dns.abcdefgh.fun.

www     IN A     36.111.184.219
static     IN CNAME     www.abcdefgh-fun.fogcdn.top
dns	IN A	36.111.166.37
ns1	IN A	36.111.166.37
ns2	IN A	36.111.166.37
*	IN A	36.111.184.219