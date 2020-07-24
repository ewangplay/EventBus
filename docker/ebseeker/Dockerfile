FROM index.alauda.cn/library/debian:jessie

RUN mkdir -p /opt/ewangplay/eventbus/bin \
        && mkdir -p /opt/ewangplay/eventbus/etc \
        && mkdir -p /opt/ewangplay/eventbus/log

ADD build/bin/ebseeker /opt/ewangplay/eventbus/bin/ebseeker

RUN chmod +x /opt/ewangplay/eventbus/bin/ebseeker

VOLUME /opt/ewangplay/eventbus/etc
VOLUME /opt/ewangplay/eventbus/log

# ebseeker tcp port
EXPOSE 4160

# ebseeker http port
EXPOSE 4161

CMD ["/opt/ewangplay/eventbus/bin/ebseeker", "--config=/opt/ewangplay/eventbus/etc/ebseeker.cfg"]
