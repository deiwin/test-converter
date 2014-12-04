# test-converter

This is a little command I quickly wrote to convert an array-driven test to a [table-driven test](http://talks.golang.org/2014/testing.slide#5).

The current implementation is very specific to an [issue (asaskevich/govalidator#13)](https://github.com/asaskevich/govalidator/issues/13) for which it was created.

Running the command:

    go get github.com/deiwin/test-converter
    go install github.com/deiwin/test-converter
    curl https://raw.githubusercontent.com/asaskevich/govalidator/a14ae891c5a02e8de7c94aabef47279b001277a4/validator_test.go | test-converter
