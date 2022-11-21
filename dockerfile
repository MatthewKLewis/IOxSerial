
# IOx app needs base-rootfs. Select either 32 or 64 bit version by uncommenting out one of the below link.

# Use below 2 lines for 32 bit Access Points
# FROM devhub-docker.cisco.com/iox-docker/ap3k/base-rootfs:latest
# RUN opkg update && opkg install python

# Use below 2 lines for 64 bit Access Points
FROM devhub-docker.cisco.com/iox-docker/ir1101/base-rootfs:latest
RUN opkg update && opkg install python

#Acquire necessary python packages
RUN pip install requests && pip install pyserial

COPY main.py /usr/bin/main.py
RUN chmod 777 /usr/bin/main.py