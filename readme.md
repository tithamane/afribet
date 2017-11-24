## About
A toy application using golang to scrape and rust to use the processed data
to try and predict a certain number of numbers the system believes will be
drawn on that day. It's more about visualising the data and seeing if there
is some sort of pattern that can make prediction possible

Eg. Based on the numbers we have, what are the four numbers that will definitely
come on on this particular day.

There is no guarantee that it will work but it's a nice project to get me coding
in rust and golang. I could use one language (Go) but I choose both so that I get
the opportunity to play with them.

## TODO
- [ ] Write scraper to download results
- [ ] Make scraper check and download lastest results every "x" minutes
- [ ] Handle errors in scraper
- [ ] Write tests for scraper methods/functions
- [ ] Write rust server to process the numbers numbers so they can be visualised