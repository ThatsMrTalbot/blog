FROM iron/base

RUN mkdir /blog.git

EXPOSE 80 80

COPY blog /
RUN chmod +x /blog
RUN apk --no-cache --no-progress add ca-certificates bash git linux-pam s6 curl openssh socat

VOLUME ["/blog.git"]

ENTRYPOINT ["/blog"]
CMD ["-path=/blog.git", "-http=:80"]
