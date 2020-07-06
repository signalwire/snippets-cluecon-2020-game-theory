const { RelayConsumer } = require('@signalwire/node');

const consumer = new RelayConsumer({
  project: '',
  token: '',
  contexts: ['FIRESTARTER'],
  ready: async ({ client }) => {
    // Consumer is successfully connected with Relay.
    // You can make calls or send messages here..
  },
  onIncomingCall: async (call) => {
    const { successful } = await call.answer()
    if (!successful) { return }

    console.log("Call answered.");

    // Generate a random number
    var winning_number = Math.floor(Math.random() * 9) + 1;
    console.log("Winning number is: " + winning_number);

    const collect = {
     type: 'digits',
     digits_max: 1,
     initial_timeout: 10,
     text: 'Please, guess a number between one and nine.'
    };

    var result = await call.promptTTS(collect);

    console.log(result);

    if(result.successful){
        var guess = result.result;
        //await call.playTTS({ text: 'You entered ' + guess })

        if(guess == winning_number){
            await call.playTTS({ text: "Winner! Winner! Chicken Dinner! You guessed correctly." });
        }else{
            await call.playTTS({ text: "You do not pass go! You have guessed incorrectly." });
        }
    }

    await call.playTTS({ text: "Thank you for playing, Good Bye!." });
    await call.hangup();

  }
});

consumer.run();
