# SVC

Stochastic Vehicular Constabulary
Snarky Vampires Cavorting
Soggy View Controller
Superb Vanishing Carrots*

Or, perhaps, a service framework of some kind.

## (But seriously....)

At [Anki](https://anki.com/) we built a lot of REST and GRPC services
in Go, and had built up a clean and fully featured framework for
building cloud service projects. It covered handling the basics of
making a well behaved process with a defined startup and graceful
shutdown lifecycle, composition of shared components, standard
handlers and interceptors, integrated with command line handling and
standardized logging and metrics and multiple authentication systems.

It was my hope that we'd be able to remove the product-specific bits
that inevitably crept in and open-source it one day, but that hope
died when Anki abruptly shut down.

Now that I'm working on some personal projects again in Go, I find
myself needing that again, so this is my attempt to recreate part of
it from scratch, but better this time. 

* Apologies to [npm](https://www.npmjs.com/)

