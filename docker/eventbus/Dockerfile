FROM index.alauda.cn/library/debian:jessie

RUN mkdir -p /opt/ewangplay/eventbus/bin \
        && mkdir -p /opt/ewangplay/eventbus/etc \
        && mkdir -p /opt/ewangplay/eventbus/data \
        && mkdir -p /opt/ewangplay/eventbus/log

ADD build/bin/eventbus /opt/ewangplay/eventbus/bin/eventbus

RUN chmod +x /opt/ewangplay/eventbus/bin/eventbus

VOLUME /opt/ewangplay/eventbus/etc
VOLUME /opt/ewangplay/eventbus/data
VOLUME /opt/ewangplay/eventbus/log

# pprof port
EXPOSE 6060

# http port
EXPOSE 8091

CMD ["/opt/ewangplay/eventbus/bin/eventbus", "--config=/opt/ewangplay/eventbus/etc/eventbus.yaml"]
