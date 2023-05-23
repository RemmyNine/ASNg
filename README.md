## ASNg
A useful script for gathering ASNs and prefixes of your target while you are conducting a wide recon process, translated from AliSkh3ykhi's script

### Key benefits over AliSkh3ykhi's ASNx script include::
1. it is not written by Ali (the most important factor)
2. it's Go based
3. It's entirely written by AI.

###### PLEASE NOTE THAT ALI SKHEYKHI IS THE SECOND PERSON WHO LOST HIS JOB BECAUSE OF AI (AFTER TOM IN TOME AND JERRY)

###### How to use?
1. clone this repo 
2. Then:
- Install following libraries:
  `github.com/miekg/dns`
  `github.com/spf13/cobra`
  `github.com/valyala/fasthttp`
3. use `go run ASNg.go` or compile it
4. Pass the arguments u want:
-h	Show this lovely help menu
-u	Single url or IP
-l	List of urls or IPs
-asn	Only show ASNs
-prefix	Only show prefixes 
