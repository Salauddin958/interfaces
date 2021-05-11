package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type bot interface {
	getGreeting() string
}

type englishBot struct {
}

type spanishBot struct {
}

type shape interface {
	getArea() float64
}

type triangle struct {
	base   float64
	height float64
}

type square struct {
	sideLength float64
}

func (t triangle) getArea() float64 {
	return 0.5 * t.base * t.height
}

func (s square) getArea() float64 {
	return s.sideLength * s.sideLength
}

func printArea(s shape) {
	fmt.Println(s.getArea())
}

func (englishBot) getGreeting() string {
	return "hello"
}

func (spanishBot) getGreeting() string {
	return "hola"
}

func printGreeting(b bot) {
	fmt.Println(b.getGreeting())
}

type logWriter struct {
}

type Aboutable interface {
	About() string
}

type Book struct {
	name string
}

func (book *Book) About() string {
	return "Book name is : " + book.name
}

type Filter interface {
	About() string
	Process([]int) []int
}

type UniqueFilter struct {
}

func (UniqueFilter) About() string {
	return "remove duplicate numbers"
}

func (UniqueFilter) Process(inputs []int) []int {
	outs := make([]int, 0, len(inputs))
	pushed := make(map[int]bool)
	for _, n := range inputs {
		if !pushed[n] {
			outs = append(outs, n)
			pushed[n] = true
		}
	}
	return outs
}

type MultipleFilter int

func (mf MultipleFilter) About() string {
	return fmt.Sprintf("keep multiples of %v", mf)
}
func (mf MultipleFilter) Process(inputs []int) []int {
	var outs = make([]int, 0, len(inputs))
	for _, n := range inputs {
		if n%int(mf) == 0 {
			outs = append(outs, n)
		}
	}
	return outs
}

func filterAndPrint(fltr Filter, unfiltered []int) []int {
	filtered := fltr.Process(unfiltered)
	fmt.Println(fltr.About()+":\n\t", filtered)
	return filtered
}

func main() {
	typeSwitchExample()
	compareInterfaces()
	numbers := []int{12, 7, 21, 12, 12, 26, 25, 21, 30}
	fmt.Println("before filtering:\n\t", numbers)

	filters := []Filter{
		UniqueFilter{},
		MultipleFilter(2),
		MultipleFilter(3),
	}
	for _, filter := range filters {
		numbers = filterAndPrint(filter, numbers)
	}
	var a Aboutable = &Book{"Go 101"}
	fmt.Println(a)
	var i interface{} = &Book{"Rust 101"}
	fmt.Println(i)
	i = a
	fmt.Println(i)
	eb := englishBot{}
	sb := spanishBot{}
	printGreeting(eb)
	printGreeting(sb)
	readingData()
	tr := triangle{base: 3, height: 4}
	sq := square{sideLength: 5}
	printArea(tr)
	printArea(sq)
	//readingFromFile()

}

func readingData() {
	res, err := http.Get("http://google.com")
	if err != nil {
		fmt.Println("Error : ", err)
		os.Exit(1)
	}

	// bs := make([]byte, 999999)
	// res.Body.Read(bs)
	// fmt.Println(string(bs))

	lw := logWriter{}

	io.Copy(lw, res.Body)
}

func (logWriter) Write(bs []byte) (int, error) {
	fmt.Println(string(bs))
	fmt.Println("len : ", len(bs))
	return len(bs), nil
}

func readingFromFile() {
	fmt.Println("Opening File")
	filename := os.Args[1]
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error in printing file")
		os.Exit(1)
	}
	io.Copy(os.Stdout, file)
}

func typeSwitchExample() {
	values := []interface{}{
		456, "abc", true, 0.33, int32(789),
		[]int{1, 2, 3}, map[int]bool{}, nil,
	}
	for _, x := range values {
		// Here, v is declared once, but it denotes
		// different variables in different branches.
		switch v := x.(type) {
		case []int: // a type literal
			// The type of v is "[]int" in this branch.
			fmt.Println("int slice:", v)
		case string: // one type name
			// The type of v is "string" in this branch.
			fmt.Println("string:", v)
		case int, float64, int32: // multiple type names
			// The type of v is "interface{}",
			// the same as x in this branch.
			fmt.Println("number:", v)
		case nil:
			// The type of v is "interface{}",
			// the same as x in this branch.
			fmt.Println(v)
		default:
			// The type of v is "interface{}",
			// the same as x in this branch.
			fmt.Println("others:", v)
		}
		// Note, each variable denoted by v in the
		// last three branches is a copy of x.
	}
}

func compareInterfaces() {
	var a, b, c interface{} = "abc", 123, "a" + "b" + "c"
	// A case of step 2.
	fmt.Println(a == b) // false
	// A case of step 3.
	fmt.Println(a == c) // true

	var x *int = nil
	var y *bool = nil
	var ix, iy interface{} = x, y
	var i interface{} = nil
	// A case of step 2.
	fmt.Println(ix == iy) // false
	// A case of step 1.
	fmt.Println(ix == i) // false
	// A case of step 1.
	fmt.Println(iy == i) // false

	// []int is an incomparable type
	var s []int = nil
	i = s
	// A case of step 1.
	fmt.Println(i == nil) // false
	// A case of step 3.
	//fmt.Println(i == i) // will panic
}
