# Leprechaun

Current Version: **0.1.1** <br />
Current Release: **Calimero**

**Leprechaun** is tool where you can schedule your recurring tasks to be performed over and over. 

In **Leprechaun** tasks are **recipes**, lets observe simple recipe file which is written using **YAML** syntax.

File is located in recipes directory which can be specified in client.ini configurational file, we will talk about this file a bit later.

    name: job1 // name of recipe
    startin: 1 // when client is started it will start in X minutes
    workevery: 10 // task will start every X minutes
    steps: // steps are done from first to last
    	- touch ./test.txt
    	- echo "Is this working?" > ./test.txt
    	- mv ./test.txt ./imwondering.txt

Steps also support variables which syntax is `$variable`. At this moment we can talk about specific section in *client.ini* file where you can define all variables to be used in your steps.

Dedicated section is named `[variables]`, anything defined in this section can and will be replaced in steps if variable defined cooresponds to syntax, sooo for example we can defined something like this

    [variables]
    testFile = ./test.txt
   
   and in our steps it will be available as `$testFile`. We can now rewrite our job file and it will look like something like this:

    name: job1 // name of recipe
    startin: 1 // when client is started it will start in X minutes
    workevery: 10 // task will start every X minutes
    steps: // steps are done from first to last
	    - touch $testFile
        - echo "Is this working?" > $testFile
        - mv $testFile ./imwondering.txt
   
Usage is very straightforward, you just need to start client and it will run recipes you defined previously.