FROM iron/base

RUN mkdir /blog.git

EXPOSE 80 80

COPY blog /
RUN chmod +x /blog

VOLUME ["/blog.git"]

ENTRYPOINT ["/blog"]
CMD ["-path=/blog.git", "-http=:80"]
