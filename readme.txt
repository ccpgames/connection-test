to use: extract to a folder with write permissions, then double-click the executable

ping test:
*this test will run the OS ping tool against TQ and collect the results
*if TQ is under DDOS mitigation, it may time out despite a functional connection.

tcp test:
*this test will attempt to perform a tcp connection to TQ, functionally the same as running "telnet 87.237.38.200 26000"
*it will return either "connection successful" or an error

port forwarding test:
*this test requests that an external service probe ports 26000 and 3724 and report whether they are forwarded or not
*this test will usually return false, as it is not neccisary to forward these ports to access TQ

http test:
*this test makes requests to http endpoints that are commonly used by the launcher/repair tool
*it will report either a successful connection or an error
*a successful connection may still be an error, if the endpoint is returning bad data. Launching this tool with the flag -keepweb=true will instruct it to save copies of the data it recieves to text files that the user may upload for examination

Note for mac client: doubleclicking the file may cause results.txt to be created in the users home directory ("cd ~" in terminal) rather than at the exec location