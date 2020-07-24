FROM index.alauda.cn/library/debian:jessie

RUN mkdir -p /opt/ewangplay/eventbus/bin \
        && mkdir -p /opt/ewangplay/eventbus/etc \
        && mkdir -p /opt/ewangplay/eventbus/log

ADD build/bin/ebadmin /opt/ewangplay/eventbus/bin/ebadmin

RUN chmod +x /opt/ewangplay/eventbus/bin/ebadmin

VOLUME /opt/ewangplay/eventbus/etc
VOLUME /opt/ewangplay/eventbus/log

# ebadmin http port
EXPOSE 4171

CMD ["/opt/ewangplay/eventbus/bin/ebadmin", "--config=/opt/ewangplay/eventbus/etc/ebadmin.cfg"]
