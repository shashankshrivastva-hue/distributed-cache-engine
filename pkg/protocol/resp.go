package protocol

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

// RESP Data Types
const (
	SimpleString = '+'
	Error        = '-'
	Integer      = ':'
	BulkString   = '$'
	Array        = '*'
)

type Parser struct {
	reader *bufio.Reader
}

func NewParser(rd io.Reader) *Parser {
	return &Parser{reader: bufio.NewReader(rd)}
}

func (p *Parser) ParseCommand() ([]string, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if len(line) < 3 || line[0] != Array {
		return nil, errors.New("invalid RESP array command")
	}

	count, err := strconv.Atoi(line[1 : len(line)-2])
	if err != nil {
		return nil, err
	}

	args := make([]string, 0, count)
	for i := 0; i < count; i++ {
		bulkHeader, err := p.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if bulkHeader[0] != BulkString {
			return nil, errors.New("expected bulk string token")
		}
		strLen, _ := strconv.Atoi(bulkHeader[1 : len(bulkHeader)-2])
		buf := make([]byte, strLen+2)
		if _, err := io.ReadFull(p.reader, buf); err != nil {
			return nil, err
		}
		args = append(args, string(buf[:strLen]))
	}
	return args, nil
}

func FormatSimpleString(msg string) string {
	return fmt.Sprintf("+%s\r\n", msg)
}

func FormatBulkString(val string) string {
	return fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)
}

func FormatNullBulk() string {
	return "$-1\r\n"
}
