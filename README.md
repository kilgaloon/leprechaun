# Leprechaun

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/b5252a28362743fa8bd94f56e5637c1c)](https://www.codacy.com/gh/kilgaloon/leprechaun/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=kilgaloon/leprechaun&amp;utm_campaign=Badge_Grade)
[![Go Report Card](https://goreportcard.com/badge/github.com/Kilgaloon/Leprechaun)](https://goreportcard.com/report/github.com/Kilgaloon/Leprechaun) [![Build Status](https://travis-ci.com/kilgaloon/leprechaun.svg?branch=master)](https://travis-ci.com/kilgaloon/leprechaun) [![codecov](https://codecov.io/gh/Kilgaloon/Leprechaun/branch/master/graph/badge.svg)](https://codecov.io/gh/Kilgaloon/Leprechaun)

**Leprechaun** is tool where you can schedule your recurring tasks to be performed over and over.

In **Leprechaun** tasks are **recipes**, lets observe simple recipe file which is written using **YAML** syntax.
  
File is located in recipes directory which can be specified in `configs.ini` configurational file. For all possible settings take a look [here](https://github.com/Kilgaloon/Leprechaun/blob/master/dist/configs/config.ini)

By definition there are 3 types of recipes, the ones that can be scheduled, the others that can be hooked and last ones that use cron pattern for scheduling jobs, they are similiar regarding steps but have some difference in definition

First we will talk about scheduled recipes and they are defined like this:

	name: job1 // name of recipe
	definition: schedule // definition of which type is recipe
	schedule:
		min: 0 // every min
		hour: 0 // every hour
		day: 0 // every day
	steps: // steps are done from first to last
		- touch ./test.txt
		- echo "Is this working?" > ./test.txt
		- mv ./test.txt ./imwondering.txt

If we set something like this

	schedule:
		min: 10 // every min
		hour: 2 // every hour
		day: 2 // every day

Task will run every 2 days 2 hours and 10 mins, if we put just days to 0 then it will run every 2 hours and 10 mins

	name: job2 // name of recipe
	definition: hook // definition of which type is recipe
	id: 45DE2239F // id which we use to find recipe
	steps:
		- echo "Hooked!" > ./hook.txt

Hooked recipe can be run by sending request to `{host}:{port}/hook?id={id_of_recipe}` on which Leprechaun server is listening, for example `localhost:11400/hook?id=45DE2239F`.

Recipes that use cron pattern to schedule tasks are used like this:

	name: job3 // name of recipe
	definition: cron // definition of which type is recipe
	pattern: * * * * *
	steps: // steps are done from first to last
		- touch ./test.txt
		- echo "Is this working?" > ./test.txt
		- mv ./test.txt ./imwondering.txt

Steps also support variables which syntax is `$variable`, and those are environment variables ex: `$LOGNAME` and in our steps it will be available as `$LOGNAME`. We can now rewrite our job file and it will look like something like this:

	name: job1 // name of recipe
	definition: schedule
	schedule:
		min: 0 // every min
		hour: 0 // every hour
		day: 0 // every day
	steps: // steps are done from first to last
		- echo "Is this working?" > $LOGNAME

Usage is very straightforward, you just need to start client and it will run recipes you defined previously.

Steps also can be defined as `sync/async` tasks which is defined by `->`, but keep in mind that steps in recipes are performed by linear path because one that is not async can block other from performing, lets take this one as example

	steps:
	- -> ping google.com
	- echo "I will not wait above task to perform, he is async so i will start immidiatelly"

but in this case for example first task will block performing on any task and all others will hang waiting it to finish
  
	steps:
	- ping google.com
	- -> echo "I need to wait above step to finish, then i can do my stuff"

## Step Pipe

Output from one step can be passed to input of next step:

	name: job1 // name of recipe
	definition: schedule
	schedule:
		min: 0 // every min
		hour: 0 // every hour
		day: 0 // every day
	steps: // steps are done from first to last
		- echo "Pipe this to next step" }>
		- cat > piped.txt

As you see, first step is using syntax `}>` at the end, which tells that this command output will be passed to next command input, you can chain like this how much you want.

## Step Failure

Since steps are executed linear workers doesn't care if some of the commands fail, they continue with execution, but you get notifications if you did setup those configurations. If you want that workers stop execution of next steps if some command failes you can specifify it with `!` like in example:

	name: job1 // name of recipe
	definition: schedule
	schedule:
		min: 0 // every min
		hour: 0 // every hour
		day: 0 // every day
	steps: // steps are done from first to last
		- ! echo "Pipe this to next step" }>
		- cat > piped.txt
		
If first step fails, recipe will fail and all other steps wont be executed

## Remote step execution

Steps can be handled by your local machine using regular syntax, if there is any need that you want specific step to be
executed by some remote machine you can spoecify that in step provided in example under, syntax is `rmt:some_host`, leprechaun will try to communicate with remote service that is configured on provided host and will run this command at that host.

	name: job1 // name of recipe
	definition: schedule
	schedule:
		min: 0 // every min
		hour: 0 // every hour
		day: 0 // every day
	steps: // steps are done from first to last
		- rmt:some_host echo "Pipe this to next step"

 Note that also as regular step this step also can pipe output to next step, so something like this is possible also:

	steps: // steps are done from first to last
		- rmt:some_host echo "Pipe this to next step" }>
		- rmt:some_other_host grep -a "Pipe" }>
		- cat > stored.txt

## Installation

Go to `leprechaun` directory and run `make install`, you will need sudo privileges for this. This will install scheduler, cron, and webhook services.

To install remote service run `make install-remote-service`, this will create `leprechaunrmt` binary.

## Build

Go to `leprechaun` directory and run `make build`. This will build scheduler, cron, and webhook services.

To build remote service run `make build-remote-service`, this will create `leprechaunrmt` binary.


## Starting/Stopping services

To start leprechaun just simply run it in background like this : `leprechaun &`

For more available commands run `leprechaun --help`

# Lepretools

For cli tools take a look [here](https://github.com/Kilgaloon/Leprechaun/blob/master/cmd/lepretools/README.md)

# Testing

To run tests with covarage `make test`, to run tests and generate reports run `make test-with-report` files will be generated in `coverprofile` dir. To test specific package run `make test-package package=[name]`
