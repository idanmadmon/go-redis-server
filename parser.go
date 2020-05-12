package go_redis_server

import (
	"errors"
	"strconv"

	uuid "github.com/satori/go.uuid"
)

var (
	simpleStringSign = '+'
	errorSign = '-'
	integerSign = ':'
	bulkStringSign = '$'
	arraySign = '*'
)

var parsec chan Request

type (
	Request struct {
		id 		uuid.UUID
		message	string
	}

	Parse struct {
		cfg			Redis
		Interrupt	*bool
	}
)

func (p *Parse) Start() error {
	if parsec == nil {
		parsec = make(chan Request, 0)
	}

	for i := 0; i < cfg.ParseWorkers; i++ {
		go p.run()
	}
	return nil
}

func (p *Parse) run() {
	for !*p.Interrupt {
		r := <-parsec
		cmds, err := ParseRequest(r.message)
		if err != nil {
			ReplyError(err, r.id)
		}
		cmdsc <- Command{r.id, cmds}
	}
}

func buildRespNullBulkString() string {
	return string(bulkStringSign) + "-1\r\n"
}

func ParseRequest(r string) ([]*string, error) {
	if err := requestValidation(r); err != nil {
		return nil, err
	}

	cmds := make([]*string, 0)
	index := 0
	toContinue := 0
	for i, c := range r {
		if toContinue > 0 {
			toContinue--
			continue
		}

		err := parseTypes(r, c, i, &cmds, &toContinue, &index)
		if err != nil {
			return nil, err
		}
	}

	return cmds, nil
}

func requestValidation(r string) error {
	if r == "" {
		return errors.New("empty request")
	}
	if r[0] != '*' {
		return errors.New("bad index, request must start with array")
	}
	return nil
}

func parseTypes(r string, c int32, i int, cmds *[]*string, toContinue, index *int) error {
	switch c {
	case simpleStringSign:
		err := parseSimpleTypes(r, i, cmds, toContinue, index)
		if err != nil {
			return err
		}
		break
	case errorSign:
		err := parseSimpleTypes(r, i, cmds, toContinue, index)
		if err != nil {
			return err
		}
		break
	case integerSign:
		err := parseSimpleTypes(r, i, cmds, toContinue, index)
		if err != nil {
			return err
		}
		break
	case bulkStringSign:
		err := parseBulkType(r, i, cmds, toContinue, index)
		if err != nil {
			return err
		}
		break
	case arraySign:
		if cmds == nil || len(*cmds) == 0 {
			err := parseArrayType(r, i, cmds, toContinue)
			if err != nil {
				return err
			}
			break
		}

		// Not necessary in the exercise
		return errors.New("bad index, can't use array as value type")
	}

	return nil
}

func parseSimpleTypes(r string, i int, cmds *[]*string, toContinue, index *int) error {
	n, err := getLengthUntilCRLF(r[i+1:])
	if err != nil {
		return err
	}

	val := r[i+1:i+1+n]
	(*cmds)[*index] = &val //escape pointer analysis
	*toContinue = 2+n //value + CRLF
	*index = *index+1
	return nil
}

func getLength(r string, i int) (int, int, error) {
	digits, err := getLengthUntilCRLF(r[i+1:])
	if err != nil {
		return 0, 0, err
	}

	n, err := strconv.Atoi(r[i+1:i+1+digits])
	if err != nil {

		return 0, 0, errors.New("bad index, length not integer, length: " + r[i+1:i+1+digits])
	}

	return digits, n, nil
}

func parseBulkType(r string, i int, cmds *[]*string, toContinue, index *int) error {
	digits, n, err := getLength(r, i)
	if err != nil {
		return err
	}

	if n == -1 {
		(*cmds)[*index] = nil
		*toContinue = 4 //-1 + CRLF
	} else {
		val := r[i+4:i+4+n]
		(*cmds)[*index] = &val //escape pointer analysis
		*toContinue = digits+4+n //count + CRLF + value + CRLF
	}

	*index = *index+1
	return nil
}

func parseArrayType(r string, i int, cmds *[]*string, toContinue *int) error {
	digits, n, err := getLength(r, i)
	if err != nil {
		return err
	}

	cmdsVal := make([]*string, n)
	*cmds = cmdsVal

	*toContinue = digits+2 //count + CRLF
	return nil
}

func getLengthUntilCRLF(s string) (int, error) {
	for i, c := range s {
		if c == '\r' {
			return i, nil
		}
	}
	return 0, errors.New("request didn't end with CRLF")
}
