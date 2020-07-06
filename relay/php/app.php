<?php

require dirname(__FILE__) . '/vendor/autoload.php';

use Generator as Coroutine;
use SignalWire\Relay\Consumer;

//$_ENV['SIGNALWIRE_PROJECT'] = "";
//$_ENV['SIGNALWIRE_TOKEN'] = "";
//$_ENV['RELAY_CONTEXT'] = "FIRESTARTER";


class FireStarterConsumer extends Consumer {
  public $project = 'd0e30b8d-58c2-481b-9188-be2ad8ab706e'; //constant($_ENV["SIGNALWIRE_PROJECT"]);
  public $token = 'PTb7f247def90655b794ae5e635e0934b13a481cef18dba35b'; //constant($_ENV["SIGNALWIRE_TOKEN"]);
  public $contexts = ['FIRESTARTER']; //[constant($_ENV["RELAY_CONTEXT"])];

//  public function __construct() {
//    self::project = $_ENV['SIGNALWIRE_PROJECT'];
//    self::token = $_ENV['SIGNALWIRE_TOKEN'];
//    self::contexts = [$_ENV['RELAY_CONTEXT']];
//  }

  public function ready(): Coroutine {
    yield;
    // Consumer is successfully connected with Relay.
    // You can make calls or send messages here..
  }

  public function randomDigit() {
    return strval( rand(1,9) );
  }

  public function onIncomingCall($call): Coroutine {
    $result = yield $call->answer();

    $winning_number = $this->randomDigit();

    if ($result->isSuccessful()) {

      $promptResult = yield $call->promptTTS([ 'type' => 'digits', 'text'=>'Please, guess a number between one and nine', 'digits_max' => 1 ]);
      if($promptResult->isSuccessful()){

        $guess = $promptResult->getResult(); // guess entered by the user
        print($guess);

        yield $call->playTTS(['text'=>"You entered {$guess}", 'gender'=>'female' ]);

        if($guess == $winning_number) {
          yield $call->playTTS(['text'=>'Winner! Winner! Chicken Dinner! You guessed correctly.', 'gender'=>'female' ]);
        }else{
          yield $call->playTTS(['text'=>'You do not pass go! You have guessed incorrectly.', 'gender'=>'female' ]);
        }

        yield $call->playTTS(['text'=>'Thank you for playing, Good Bye!', 'gender'=>'female' ]);
        yield $call->hangup();
      }
    }
  }
}

$consumer = new FireStarterConsumer();
$consumer->run();
