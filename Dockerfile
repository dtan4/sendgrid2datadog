FROM ruby:2.3.3-alpine
MAINTAINER Daisuke Fujita <dtanshi45@gmail.com> (@dtan4)

WORKDIR /app
COPY Gemfile /app/
COPY Gemfile.lock /app/

RUN apk add --no-cache --update --virtual=build-deps \
      g++ \
      make \
    && bundle install --without test development --system -j4 \
    && apk del build-deps

COPY . /app

EXPOSE 9292

CMD ["bundle", "exec", "rackup", "-E", "production", "-p", "9292"]
