#  Getting started with ~~ u-root ~~


Prerequisites
=============

To execute u-root and play with it, you need to have git, golang installed. 
On a Debian, Ubuntu or other .deb system, you should be able to get going with:

	sudo aptitude install git golang build-essential


U-root requires go1.5 and the go source tree needs to be present.

GERRIT
======

We use gerrithub.io for code-review. If you want to submit changes, go to

	https://review.gerrithub.io/#/admin/projects/u-root/u-root

and check out the repository from gerrithub rather than github. The clone
command will probably look something like this:

	git clone ssh://USERNAME@review.gerrithub.io:29418/u-root/u-root

you'll need to run a few commands inside the top-level directory to get set
up for code-review:

	cd u-root
	curl -Lo .git/hooks/commit-msg http://review.gerrithub.io/tools/hooks/commit-msg
	chmod u+x .git/hooks/commit-msg
	git config remote.origin.push HEAD:refs/for/master
	git config remote.origin.receivepack "git receive-pack --reviewer rminnich --reviewer rhiguita"

You're now all set, you can build the whole thing just by running the README file, or you could just type:

	go run scripts/ramfs.go -test


	
