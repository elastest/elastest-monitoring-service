GRAMMAR_FILE=striver
PEG_FILE=${GRAMMAR_FILE}.peg
OUTPUT=parser.go


all: parser build

parser:
	pigeon -o ${OUTPUT} ${PEG_FILE}

build:
#	go build stamp.go
#	go build && go install
	go build 

clean:
	rm -f ${OUTPUT}
