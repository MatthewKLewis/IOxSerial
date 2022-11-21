# IOx app needs base-rootfs. Select either 32 or 64 bit version by uncommenting out one of the below link.

# Use below 2 lines for 64 bit Access Points
FROM devhub-docker.cisco.com/iox-docker/ir1101/base-rootfs:latest
RUN opkg update && opkg install python3 && opkg install python3-pip

#Acquire necessary python packages
RUN pip3 install requests 
RUN pip3 install pyserial==2.7

COPY main.py /usr/bin/main.py
RUN chmod 777 /usr/bin/main.py
