# picsdump

This project uses powerful `gopacket` library to sniff and decode packet in a `.cap` file. Every HTTP response packet with `jpeg` or `png` inside will be captured and saved.

## how to get a pcap file
   
	tcp dump -w /dir/to/xxx.cap
   
## how to run

	./picsdump -f filename

