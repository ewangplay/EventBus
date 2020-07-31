FROM index.alauda.cn/library/debian:jessie

RUN mkdir -p /opt/eventbus/bin \
        && mkdir -p /opt/eventbus/etc \
        && mkdir -p /opt/eventbus/log

ADD build/bin/eventbus /opt/eventbus/bin/eventbus

RUN chmod +x /opt/eventbus/bin/eventbus

VOLUME /opt/eventbus/etc
VOLUME /opt/eventbus/log

# pprof port
EXPOSE 6060

# http port
EXPOSE 8091

CMD ["/opt/eventbus/bin/eventbus", "--config=/opt/eventbus/etc/eventbus.yaml"]
