/*
 *  Porter Stemmer Algorithm in Go (Test)
 *  Developed by Suhail Patel <me@suhailpatel.com>
 */
package stemmer

import (
    "testing"
    "bufio"
    "os"
    "strings"
)

// I actually wanted to have a test containing all the words used in the
// original paper but Martin Porter (author of the original algorithm)
// has provided a massive collection of ~23000 words with the expected
// output of the algorithm so we can test our implementation against
// the original paper implementation
func TestCorpus(t *testing.T) {
	input, errIn := os.Open("corpus/test_input.txt")
    output, errOut := os.Open("corpus/test_output.txt")
    
    defer input.Close()
    defer output.Close()
    
    if errIn != nil || errOut != nil {
        t.Fatalf("Could not read input or output test files [%s, %s]", errIn, errOut)
    }
    
    inScan := bufio.NewScanner(input)
    outScan := bufio.NewScanner(output)
    
    for inScan.Scan() && outScan.Scan() {
        in := inScan.Text()
        out := outScan.Text()
        stemmed := Stem(in)
        
        t.Logf("[PASS] Input: %s → Expected: %s, Stemmed: %s\n", in, out, stemmed)
        
        if (!strings.EqualFold(out, stemmed)) {
            t.Errorf("[FAIL] Expected %s but got %s for input %s\n", out, stemmed, in)
        }
    }
    
    if inScan.Err() != nil || outScan.Err() != nil {
        t.Fatalf("Could not open scanner for input or output test files [%s, %s]", inScan.Err(), outScan.Err())
    }
}