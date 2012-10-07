package lastline

import "fmt"
import "io"
import "runtime"
import "testing"

func generateIntLines(t *testing.T, count int, output io.Writer) {
	for i := 0; i < count; i++ {
		_, err := fmt.Fprintf(output, "%010d\n", i)
		if err != nil {
			t.Fatalf("Error outputting to LastLine: %v", err)
		}
		runtime.Gosched()
	}
}

func TestWrite(t *testing.T) {
	ll := NewLastLine()
	go generateIntLines(t, 1000, ll)
}

func TestWriteRead(t *testing.T) {
	ll := NewLastLine()
	go generateIntLines(t, 1000, ll)
	lastValue := ""
	for i := 0; i < 1000; i++ {
		value := ll.GetLastLine()
		t.Logf("Got value: %s", value)
		if value < lastValue {
			t.Errorf("Last lines aren't increasing: Got %s after %s.", value, lastValue)
		}
		lastValue = value
	}
}

func TestCopy(t *testing.T) {
	ll := NewLastLine()
	llCopy := ll
	fmt.Fprintf(ll, "ala")
	ll.Close()
	for ll.GetLastLine() == "" {
		runtime.Gosched()
	}
	if ll.GetLastLine() != "ala" {
		t.Errorf("Last line is %s instead of ala", llCopy.GetLastLine())
	}
	if llCopy.GetLastLine() != "ala" {
		t.Errorf("Last line on copy is %s instead of ala", llCopy.GetLastLine())
	}
}
