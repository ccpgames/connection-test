to use: extract to a folder with write permissions, then double-click the executable

ping test:
*this test will run the OS ping tool against TQ and google DNS and collect the results
*a failure to ping 8.8.8.8 may indicate a general failure of the users internet (8.8.8.8 is extremely reliable)
*if TQ is under DDOS mitigation, the tq ping may time out despite a functional connection.

unlimited ping mode:
*this will, after conducting all other tests, ping tranquility until stopped
*this may be useful for testing connection stability
*if TQ is under DDOS mitigation, the tq ping may time out despite a functional connection.
*the user can stop the ping at any time by closing the window or pressing ctrl-c
*the results file can be copied elsewhere and submitted for examination while the unlimited ping is running

tcp test:
*this test will attempt to perform a tcp connection to TQ on ports 26000 and 3724
*it will return either "connection successful" or an error

http test:
*this test makes requests to http endpoints that are commonly used by the launcher/repair tool
*it will report either a successful connection or an error
*a successful connection may still be an error, if the endpoint is returning bad data. On launch, the user can instruct the tool to save the data it receives, for upload and examination

Note for mac client: double-clicking the file may cause results.txt to be created in the users home directory ("cd ~" in terminal, favourites->all my files in Finder) rather than at the exec location