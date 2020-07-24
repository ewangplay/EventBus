FROM index.alauda.cn/library/debian:jessie

RUN mkdir -p /opt/ewangplay/eventbus/bin \
        && mkdir -p /opt/ewangplay/eventbus/etc \
        && mkdir -p /opt/ewangplay/eventbus/data \
        && mkdir -p /opt/ewangplay/eventbus/log

ADD build/bin/ebnode /opt/ewangplay/eventbus/bin/ebnode

RUN chmod +x /opt/ewangplay/eventbus/bin/ebnode

VOLUME /opt/ewangplay/eventbus/etc
VOLUME /opt/ewangplay/eventbus/data
VOLUME /opt/ewangplay/eventbus/log

# nsqd tcp port
EXPOSE 4150

# nsqd http port
EXPOSE 4151

# pprof port
EXPOSE 6060

# eventbus http port
EXPOSE 8091

CMD ["/opt/ewangplay/eventbus/bin/ebnode", "--config=/opt/ewangplay/eventbus/etc/ebnode.yaml"]
