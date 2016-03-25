# go-readthedocs

A project that'll crawl an org and determine whether the READMEs meet the organizational standard.

What I've determined so far:
* will search for all repositories that follow a certain naming convention
* create a scorecard of whether docs are pretty

What does it mean to have "pretty" docs? Still an open question, but will likely:
* First ensure a README.md exists in the root directory and has content in it
* Next, does it have each section we've determined is needed?
* Do the sections have some reasonable length of characters in them?
* Do they have the keywords I'd expect (a valid version of Go mentioned, supported OSs, others I can put into a grid of expected results)?

Now.. let's go build this.

Docs I read:
https://godoc.org/github.com/google/go-github/github#Organization
http://blog.fmpwizard.com/blog/go_making_a_program_70_faster_by_avoiding_common_mistakes?imm_mid=0e1cf3&cmp=em-prog-na-na-newsltr_20160319
