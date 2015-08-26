require "dogapi"
require "json"
require "sinatra/base"

class App < Sinatra::Base
  configure do
    set :dog, Dogapi::Client.new(ENV["DD_API_KEY"])
  end

  def dog
    settings.dog
  end

  get "/" do
    "https://github.com/dtan4/sendgrid2datadog"
  end

  post "/webhook" do
    events = JSON.parse(request.body.read)

    events.group_by { |event| event["event"] }.each do |type, evs|
      points = evs.map { |ev| [Time.at(ev["timestamp"]), 1] }
      dog.emit_points("sendgrid.event.#{type}", points, type: "counter")
    end

    "Events was sent to Datadog"
  end
end
