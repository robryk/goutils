package lastline

import "bufio"
import "io"

type LastLine struct {
	lastLineChan chan string
	lastLine     *string
	io.WriteCloser
}

func NewLastLine() LastLine {
	emptyString := ""
	ll := LastLine{
		lastLineChan: make(chan string),
		lastLine:     &emptyString,
	}

	var pr io.Reader
	pr, ll.WriteCloser = io.Pipe()

	bufr := bufio.NewReader(pr)

	logchan := make(chan string)

	go func() {
		defer close(logchan)
		for {
			line, err := bufr.ReadString('\n')
			if err == io.EOF {
				logchan <- line
				return
			}
			if err != nil {
				return
			}
			logchan <- line[:len(line)-1]
		}
	}()

	go func() {
		defer close(ll.lastLineChan)
		for {
			select {
			case line, ok := <-logchan:
				if ok {
					*ll.lastLine = line
				} else {
					return
				}
			case ll.lastLineChan <- *ll.lastLine:
			}
		}
	}()

	return ll
}

func (ll LastLine) GetLastLine() string {
	line, ok := <-ll.lastLineChan
	if !ok {
		line = *ll.lastLine
	}
	return line
}
