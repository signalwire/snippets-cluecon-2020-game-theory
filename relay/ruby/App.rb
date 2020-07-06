require "signalwire"

class FireStarter < Signalwire::Relay::Consumer
  contexts ['FIRESTARTER']

  def on_incoming_call(call)
    call.answer

    winning_number = rand(1..9)
    puts "Winning number is: " + winning_number.to_s

    tts = [{ type: 'tts', params: { text: 'Please, guess a number between one and nine' } }]
    result = call.prompt(initial_timeout: 10, digits_max: 1, play: tts)

    if result.successful)
        guess = result.result # guess entered by user


    end

    call.play_tts 'Thank you for playing, Good Bye!'
    call.hangup
  end
end

FireStarter.new(project: '', token: '').run
