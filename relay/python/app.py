import os
import requests
import pprint
import random
import json

# 1 (585) 653-1337 | FIRESTARTER
from signalwire.relay.consumer import Consumer

os.environ['SIGNALWIRE_PROJECT'] = ""
os.environ['SIGNALWIRE_TOKEN'] = ""
os.environ['RELAY_CONTEXT'] = "FIRESTARTER"

class FireStarterConsumer(Consumer):
    def setup(self):
        self.project = os.environ['SIGNALWIRE_PROJECT']
        self.token = os.environ['SIGNALWIRE_TOKEN']
        self.contexts = [os.environ['RELAY_CONTEXT']]

    #async def ready(self):
        # Consumer is successfully connected with Relay.
        # You can make calls or send messages here..

    async def on_incoming_call(self, call):
        result = await call.answer()
        if result.successful:
            print('Call answered..')

            # Generate a randome digit between 1 - 9
            winning_number = str(random.randint(1,9))
            print("Winning number is: " + winning_number)

            result = await call.prompt_tts(prompt_type='digits', text='Please, guess a number between one and 9', digits_max=1)
            if result.successful:
                guess = result.result # guess entered by the user
                await call.play_tts(text="You entered " + str(guess), gender='female')
                if(guess == winning_number):
                    await call.play_tts(text="Winner! Winner! Chicken Dinner! You guessed correctly.", gender='female')
                else:
                    await call.play_tts(text="You do not pass go! You have guessed incorrectly.", gender='female')

                await call.play_tts(text="Thank you for playing, Good Bye!", gender='female')
                await call.hangup()

# Run your consumer..
consumer = FireStarterConsumer()
consumer.run()

#def generate_random()
#    rand_number = str(random.randint(1,10))
#    return rand_number
