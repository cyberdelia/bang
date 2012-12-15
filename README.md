# bang

Bang is a tool to send load to your web app.

## Installation

Download and install :

    $ git clone git://github.com/cyberdelia/bang.git
    $ git install

## Usage

To run 10 workers during 10 seconds :

    $ bang -url http://localhost:5000 -concurrency=10 -duration=10s
	Running 10 workers for at least 10s
	Starting to load the server
	Successful calls  		     68167
	Total time        		        10.00s
	Fastest           		         0.00s
	Slowest           		         0.01s
	Average           		         0.00s
	Standard deviation		         0.00s
	Mean rate         		      6815.36
	1-min rate        		      6813.85
	5-min rate        		      6813.49
	15-min rate       		      6813.63

## About

Bang is born out of the curiosity to build something like
[boom](https://github.com/tarekziade/boom) in Go.
