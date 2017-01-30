# tcp-flags-tools

Set of functions to create byte 13 of TCP header using command ligne 

The function recieves a command line where the following operators can be used:

* = => MATCH
* ! => NOT
* & => AND
* ' ' => OR

and where the following TCP flags could be set:
* 'C' as CWR
* 'E' as ECE
* 'U' as URG
* 'A' as ACK
* 'P' as PUSH
* 'R' as RST
* 'S' as SYN
* 'F' as FIN

The function returns two slices : one with TCP flags value and one as BGP Flowspec operator as describe in RFC 5575.

Example of commands :

* "=SA&=A" means traffic that have exact match of TCP flags SYN/ACK AND excat match of ACK,
* "SA" means partial match of SYN/ACK flags
* "=!U&=!A" means not match of URG and not match of ACK.



