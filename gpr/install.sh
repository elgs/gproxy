export GOPATH=/Volumes/User/Home/Dev/Go/gpr
rm -Rf $GOPATH/bin

GOOS=windows GOARCH=386 go install
GOOS=windows GOARCH=amd64 go install

GOOS=linux GOARCH=386 go install
GOOS=linux GOARCH=amd64 go install
GOOS=linux GOARCH=arm go install

GOOS=darwin GOARCH=386 go install
GOOS=darwin GOARCH=amd64 go install

GOOS=freebsd GOARCH=386 go install
GOOS=freebsd GOARCH=amd64 go install

GOOS=netbsd GOARCH=386 go install
GOOS=netbsd GOARCH=amd64 go install

GOOS=openbsd GOARCH=386 go install
GOOS=openbsd GOARCH=amd64 go install

cp gpr.json $GOPATH/bin

mkdir $GOPATH/bin/darwin_amd64
mv $GOPATH/bin/gpr $GOPATH/bin/darwin_amd64
