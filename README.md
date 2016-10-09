# alytfeed
========

This repository contains the RSS feed for the [At Least You're Trying](http://gtradio.net/alyt) podcast.
It may seem odd to have my podcast feed on github, but I assure you that it's a key piece of an overly-complex (but valuable) workflow.
I configured the github webook to notify the server hosting the official podcast feed every time I push a new episode into the feed.
A Sinatra service running on my VPS over at [Linode](https://www.linode.com/?r=30991a143a3c99716fbc7fdcf81355338c4d2b64) takes care of the subsequent cloning and file system housekeeping.
This scheme also happens to provide a full version history and an offsite backup of this precious XML file.

**But unit tests on a podcast feed?**

Yes, this feed has some tests written in Go.
This was my first experience writing Go code, and I chose it due to its simple testing story...`go test`.
Since I update the feed by hand, I don't want to be able to make a typo in the publication date and have my podcasting empire come crashing down (this has happened a few times).
For now, the tests run on a local git commit hook. In the future, I hope to move their execution up to Github/TravisCI.

Here are the contents of .git/hooks/pre-commit (don't forget to grant execution permissions!):

    #!/bin/sh

    status=0
    badtest="$(go test | grep "FAIL:" | awk '{print $3}')"
    if test -n "$badtest" ; then
        for bad in $badtest; do
            echo "git pre-commit check failed: go test failed: $bad"
        done
        status=1
    fi

    exit $status

## Become a listener...Official Feed Locations:
- RSS, www.kaplon.us/alytfeed
- iTunes, https://itunes.apple.com/us/podcast/at-least-youre-trying/id702153446?mt=2&ign-mpt=uo%3D4

## Interact
- http://gtradio.net/alyt
- https://www.facebook.com/AtLeastYoureTrying
- alyt.show@gmail.com

Thanks to [archive.org](https://archive.org/donate) for media hosting and bandwidth.

