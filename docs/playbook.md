# Updating the Docker test images

cd .circleci/images

update builduploadall.sh with a new version (e.g., v4.4.0 -> v4.5.0)

--- a/.circleci/images/builduploadall.sh
+++ b/.circleci/images/builduploadall.sh
-VERSION=v4.4.0
+VERSION=v4.5.0

run that script. It may fail due to packages becoming unavailable
(e.g. python is always fun) but, if it works, the new docker
images will be uploaded.

Once this is done, update .circleci/config.yml to use the new container version

+++ .circleci/config.yml
@@ -104,7 +104,7 @@ jobs:
-      - image: uroottest/test-image-tamago:v4.4.0
+      - image: uroottest/test-image-tamago:v4.5.0
@@ -125,7 +125,7 @@ jobs:
-      - image: uroottest/test-image-tamago:v4.4.0
+      - image: uroottest/test-image-tamago:v4.5.0


Then make the commit and push the changes to these files.
See https://github.com/u-root/u-root/pull/2732 for an example.

## Updating tamago release
Sometimes you will need a new version of Tamago, in which case
you need to change the builduploadall.sh script as above, but also
specify a new Tamago version. You must also update the sha256sum
in the test-image-tamago/Dockerfile file.

Find the latest tamago release at https://github.com/usbarmory/tamago-go/
get that image, e.g.
TAMAGO_VERSION=1.20.6 wget -O tamago-go.tgz https://github.com/usbarmory/tamago-go/releases/download/tamago-go${TAMAGO_VERSION}/tamago-go${TAMAGO_VERSION}.linux-amd64.tar.gz
get the sha256sum:
sha256sum tamago-go.tgz

and update these two variables in .circleci/images/test-image-tamago/Dockerfile.

In this case, it looks like this:
--- a/.circleci/images/test-image-tamago/Dockerfile
+++ b/.circleci/images/test-image-tamago/Dockerfile
-ENV TAMAGO_VERSION="1.20.3"
-ENV TAMAGO_CHECKSUM="7657d39b8e062f85433a48367012d138ef5b928408b0b41e15384fe150495082"
+ENV TAMAGO_VERSION="1.20.6"
+ENV TAMAGO_CHECKSUM="6319b1778e93695b62bb63946c5dd28c4d8f3c1ac3c4bf28e49cb967d570dfd5"

Then follow the instructions, above, for updating the test images.
