FROM golang:1.15
COPY ./main /go/src/main
COPY ./mapreduce /go/src/mapreduce
CMD cd "$GOPATH/src/main" && go run wc.go master sequential pg-*.txt && sort -n -k2 mrtmp.wcseq | tail -10 > top_output.txt && diff top_output.txt mr-testout.txt && if [ $? -eq 0 ]; then echo "Successful"; fi
