#!/bin/bash
 
# Initial status message
echo -n "Pinging IPs"
 
# Create a temporary file for storing ARP results
arp_tmp_file=$(mktemp)
 
# Ping each IP in the 192.168.1.2 to 192.168.1.100 range
for i in {2..100}; do
    ping -c 1 -W 1 192.168.1.$i > /dev/null 2>&1
    echo -n "."
done
echo ""  # Move to a new line after the dots
 
# Store ARP cache in the temp file
arp -a | grep -v '<incomplete>' > $arp_tmp_file
 
# Check each IP with a MAC for port 9999
cat $arp_tmp_file | awk '{print $2, $4}' | tr -d '()' | while read ip mac; do
    # Use nc to check port 9999 with a timeout of 1 second
    nc -z -w 1 $ip 9999
    if [ $? -eq 0 ]; then
        echo "$ip has MAC $mac with port 9999 open"
    fi
done
 
# Clean up
rm -f $arp_tmp_file
 
echo "Data collection complete!"