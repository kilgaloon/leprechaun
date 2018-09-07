
# Leprechaun

  
[![Go Report Card](https://goreportcard.com/badge/github.com/Kilgaloon/Leprechaun)](https://goreportcard.com/report/github.com/Kilgaloon/Leprechaun)

[![Build Status](https://travis-ci.com/Kilgaloon/Leprechaun.svg?branch=master)](https://travis-ci.com/Kilgaloon/Leprechaun)
  

Current Version: **0.6.0**  <br  />
Current Release: **Calimero**

**Leprechaun** is tool where you can schedule your recurring tasks to be performed over and over.

In **Leprechaun** tasks are **recipes**, lets observe simple recipe file which is written using **YAML** syntax.

File is located in recipes directory which can be specified in client.ini configurational file, we will talk about this file a bit later.

By definition there are 2 types of recipes, the ones that can be scheduled and the others that can be hooked, they are similiar regarding steps but have some difference in definition 

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


Steps also support variables which syntax is `$variable`. At this moment we can talk about specific section in *client.ini* file where you can define all variables to be used in your steps.

  
Dedicated section is named `[variables]`, anything defined in this section can and will be replaced in steps if variable defined cooresponds to syntax, sooo for example we can defined something like this

	[variables]
	testFile = ./test.txt

Also all environment variables are available in steps ex: `$LOGNAME`

and in our steps it will be available as `$testFile`. We can now rewrite our job file and it will look like something like this:

	name: job1 // name of recipe
	definition: schedule
	schedule:
		min: 0 // every min
		hour: 0 // every hour
		day: 0 // every day
	steps: // steps are done from first to last
		- touch $testFile
		- echo "Is this working?" > $testFile
		- mv $testFile ./imwondering.txt

Usage is very straightforward, you just need to start client and it will run recipes you defined previously.
  

## Starting/Stopping service

To start leprechaun just simply run it in background like this : `leprechaun &`
and you can stop leprechaun simple as this: `leprechaun stop`