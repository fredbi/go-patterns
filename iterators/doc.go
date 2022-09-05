// Package iterators provides generic iterators, that is
// constructs that may be iterated over using Next() in a
// loop, then calling some Item() (T, error) method to
// retrieve the current iterated value.
//
// When iteration is complete, the iterator resources should be
// relinquished using Close().
package iterators
