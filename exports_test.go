package main

import "testing"

func checkResult(t *testing.T, code string, expected cell) {
	var actual cell
	var err error

	actual, err = Execute(code)

	if err != nil {
		t.Errorf("Error occured while executing code: %s\n", err.Error())
	}

	if !actual.Equal(expected) {
		t.Errorf("Expected %s but got %s\n", expected, actual)
	}
}

func TestOperators(t *testing.T) {
	// Lists
	checkResult(t, "(cons 1 2 3)", newConsCell(
		newFixNumCell(1),
		newConsCell(
			newFixNumCell(2),
			newConsCell(
				newFixNumCell(3),
				newNilCell(),
			),
		),
	))

	checkResult(t, "(car (cons 1 2))", newFixNumCell(1))
	checkResult(t, "(cdr (cons 1 2))", newConsCell(newFixNumCell(2), newNilCell()))

	checkResult(t, "(atom 5)", newTrueCell())
	checkResult(t, "(atom ())", newNilCell())
	checkResult(t, "(atom (cons 1 2))", newNilCell())

	checkResult(t, "(append () ())", newNilCell())

	checkResult(t, "(append (cons 1 2) (cons 3 4) (cons 5 6))", newConsCell(
		newFixNumCell(1),
		newConsCell(
			newFixNumCell(2),
			newConsCell(
				newFixNumCell(3),
				newConsCell(
					newFixNumCell(4),
					newConsCell(
						newFixNumCell(5),
						newConsCell(
							newFixNumCell(6),
							newNilCell(),
						),
					),
				),
			),
		),
	))

	checkResult(t, "(append (cons 1 2) () (cons 5 6))", newConsCell(
		newFixNumCell(1),
		newConsCell(
			newFixNumCell(2),
			newConsCell(
				newFixNumCell(5),
				newConsCell(
					newFixNumCell(6),
					newNilCell(),
				),
			),
		),
	))

	// Arithmetic
	checkResult(t, "(+ 6 5)", newFixNumCell(11))
	checkResult(t, "(- 6 5)", newFixNumCell(1))
	checkResult(t, "(* 6 5)", newFixNumCell(30))
	checkResult(t, "(/ 30 5)", newFixNumCell(6))
	checkResult(t, "(- 6)", newFixNumCell(-6))

	// Equality
	checkResult(t, "(= 5 5)", newTrueCell())
	checkResult(t, "(= 5 6)", newNilCell())
	checkResult(t, "(= (cons 1 2) (cons 1 2))", newTrueCell())
	checkResult(t, "(= () ())", newTrueCell())
	checkResult(t, "(= (cons 1 2) ())", newNilCell())

	checkResult(t, "(!= 5 5)", newNilCell())
	checkResult(t, "(!= 5 6)", newTrueCell())
	checkResult(t, "(!= (cons 1 2) (cons 1 2))", newNilCell())
	checkResult(t, "(!= () ())", newNilCell())
	checkResult(t, "(!= (cons 1 2) ())", newTrueCell())

	checkResult(t, "(> 2 1)", newTrueCell())
	checkResult(t, "(> 1 1)", newNilCell())
	checkResult(t, "(> 1 2)", newNilCell())

	checkResult(t, "(>= 2 1)", newTrueCell())
	checkResult(t, "(>= 1 1)", newTrueCell())
	checkResult(t, "(>= 1 2)", newNilCell())

	checkResult(t, "(< 2 1)", newNilCell())
	checkResult(t, "(< 1 1)", newNilCell())
	checkResult(t, "(< 1 2)", newTrueCell())

	checkResult(t, "(<= 2 1)", newNilCell())
	checkResult(t, "(<= 1 1)", newTrueCell())
	checkResult(t, "(<= 1 2)", newTrueCell())
}

func TestIfCondition(t *testing.T) {
	checkResult(t, "(if (= 10 10) 1 2)", newFixNumCell(1))
	checkResult(t, "(if (= 10 11) 1 2)", newFixNumCell(2))
	checkResult(t, "(if (= 10 10) (* 2 2) 2)", newFixNumCell(4))
	checkResult(t, "(if (= 10 11) 1 (* 4 4))", newFixNumCell(16))
}

func TestFunctionApplication(t *testing.T) {
	checkResult(t, "((fn () (+ 6 5)))", newFixNumCell(11))
	checkResult(t, "((fn (a b) (+ a b)) 6 5)", newFixNumCell(11))
}

func TestLet(t *testing.T) {
	checkResult(t, "(let (a 1) (b (+ 2 2)) (c (* 3 3)) (* b (+ a c)))", newFixNumCell(40))
}

func TestSequences(t *testing.T) {
	checkResult(t, "(- 3 2)\n(+ 6 5)", newFixNumCell(11))
	checkResult(t, "((fn (a b) (- a b) (+ a b)) 6 5)", newFixNumCell(11))
}

func TestStorage(t *testing.T) {
	checkResult(t, "(def a 1) (def b 2) (+ a b)", newFixNumCell(3))
	checkResult(t, "(def add (a b) (+ a b))\n(def square (x) (* x x))\n(square (add 6 5))", newFixNumCell(121))
}

func TestRecursion(t *testing.T) {
	code := `
(def fac (n)
     (if (= n 0)
         1
         (* n (fac (- n 1)))))

(fac 6)
`

	checkResult(t, code, newFixNumCell(720))

	code = `
(def filter (f ls)
     (if ls
         (let
             (head (car ls))
             (tail (cdr ls))
             (if (f head)
                 (append (cons head) (filter f tail))
                 (filter f tail)))
     ls))

(def qsort (ls)
     (if ls
         (let
             (head (car ls))
             (tail (cdr ls))
             (let
                 (lo (filter (fn (n) (< n head)) tail))
                 (hi (filter (fn (n) (>= n head)) tail))
                 (append (qsort lo) (cons head) (qsort hi))))
         ls))

(qsort (cons 2 5 3 6 1 4))
`

	checkResult(t, code, newConsCell(
		newFixNumCell(1),
		newConsCell(
			newFixNumCell(2),
			newConsCell(
				newFixNumCell(3),
				newConsCell(
					newFixNumCell(4),
					newConsCell(
						newFixNumCell(5),
						newConsCell(
							newFixNumCell(6),
							newNilCell(),
						),
					),
				),
			),
		),
	))
}
