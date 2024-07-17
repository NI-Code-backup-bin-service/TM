#!/bin/bash
GOOS=linux go test -c -o LdapTest.exe
# This upload needs to be run from the WL office for it to work;
# otherwise, you can upload manually via the usual process of
# uploading to Artifactory when working from home
curl -urobert_wl_test:AP4dJ1fydBqae3PHeZNgFZn9n3C -f -T LdapTest.exe https://networkinternational.jfrog.io/networkinternational/WL_Test_Repo/misc/LdapTest.exe
