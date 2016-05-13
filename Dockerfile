FROM golang

EXPOSE 80 80

RUN mkdir /blog.git
RUN mkdir -p $GOPATH/src/github.com/ThatsMrTalbot/blog

COPY . $GOPATH/src/github.com/ThatsMrTalbot/blog
RUN go get github.com/ThatsMrTalbot/blog/cmd/blog/...
RUN go install github.com/ThatsMrTalbot/blog/cmd/blog

VOLUME ["/blog.git"]

ENTRYPOINT ["blog"]
CMD ["-path=/blog.git", "-http=:80"]
