# Note you cannot run golang binaries on Alpine directly
FROM            debian:buster-slim

MAINTAINER      chris@shokunin.co

COPY            micro-leaderboard /micro-leaderboard

WORKDIR		/
ENV		GIN_MODE=release

EXPOSE          8080

ENTRYPOINT      [ "/micro-leaderboard" ]
