**Author: Maxwell You**

**Assignment: Anonymous wget**

## Setup
Place these files and directories under a `$GOPATH/src/proj3` directory to ensure correct behavior.

## Logistics
My program will look for the default "chaingang.txt" in the proj3 directory. Otherwise it will look for the chainfile at the location specified.

awget assumes the arguments will be provided in the order:
    
    awget URL [chainfile]`

If no port is specified for ss, an open port is chosen

The downloaded file will take the name from the URL if there is one, otherwise, the name will be "index.html"
I parsed the filename from the URL assuming there was a dot (.) in the path portion of the URL

## Functionality
awget initiates a request for a file through a URL.
It chooses a random stepping stone in the chainfile and forwards it there.
Stepping stones will strip themselves from the stepping stone list and check if the resulting list is empty.
If it is, the stepping stone will issue a GET request with the URL. The file will be sent back to the previous stepping stones until it reaches the client that initiated awget.
If it is not, the stepping stone will send the request to a random stepping stone.

I have tested this program with the some websites as links (TOR, google, etc), a 174MB file from Imgur, and a few books > 1MB from Project Gutenberg.
I was able to retrieve the files and inspect their contents to find that they downloaded correctly.
