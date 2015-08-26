require "json"
require "sinatra/base"
require "statsd"

class App < Sinatra::Base
  configure do
    set :statsd, Statsd.new("localhost", 8125)
  end

  def statsd
    settings.statsd
  end

  get "/" do
    "https://github.com/dtan4/sendgrid2datadog"
  end

  post "/webhook" do
    events = JSON.parse(request.body.read)

    statsd.batch do |s|
      events.each do |event|
        s.increment("sendgrid.event.#{event['type']}")
      end
    end

    "Events was sent to Datadog"
  end
end
