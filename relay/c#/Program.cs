using SignalWire.Relay;
using SignalWire.Relay.Calling;
using System;
using System.Collections.Generic;

namespace FireStarter
{
    class FireStarter : Consumer
    {
        protected override void Setup()
        {
            Project = "";
            Token = "";
            Contexts = new List<string> { "FIRESTARTER" };
        }

        private string RandomDigit()
        {
            Random rand = new Random();
            return rand.Next(1, 9).ToString();
        }

        protected override void OnIncomingCall(Call call)
        {
            AnswerResult resultAnswer = call.Answer();
            if (!resultAnswer.Successful) return;

            // Request a random digit between 1 - 9
            string winning_number = RandomDigit();            
            Console.WriteLine("Winning number is:" + winning_number);

            PromptResult guessResult = call.PromptTTS("Please, guess a number between one and nine.", 
                new CallCollect() {
                Digits = new CallCollect.DigitsParams
                {
                    Max = 1
                }
            }, "female");

            if (guessResult.Result == winning_number)
            {
                call.PlayTTS("Winner! Winner! Chicken Dinner! You guessed correctly.", "female");
            }
            else
            {
                call.PlayTTS("You do not pass go! You have guessed incorrectly.", "female");
            }

            call.PlayTTS("Thank you for playing, Good Bye!", "female");
            call.Hangup();
        }
    }
    class Program
    {
        static void Main(string[] args)
        {
            Console.WriteLine("Hello World!");

            FireStarter consumer = new FireStarter();

            consumer.Run();
        }
    }
}
