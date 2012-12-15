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
	Successful calls  		     67040
	Total time        		        10.00s
	Fastest           		         0.00s
	Slowest           		         0.01s
	Mean              		         0.00s
	Standard deviation		         0.00s
	Median            		         0.00s
	75th percentile   		         0.00s
	95th percentile   		         0.00s
	99th percentile   		         0.00s
	99.9th percentile 		         0.00s
	Mean rate         		      6702.37
	1-min rate        		      6757.18
	5-min rate        		      6765.09
	15-min rate       		      6766.49

## Options

Bang tries to be helpful and has plenty of options : 

	$ bang -h
	Usage of bang:
	  -body="": Request body
	  -concurrency=10: Concurrency
	  -content-type="text/plain": Content-Type
	  -duration="10s": Duration
	  -method="GET": HTTP method
	  -url="": URL to hit


## About

Bang is born out of the curiosity to build something like
[boom](https://github.com/tarekziade/boom) in Go.
