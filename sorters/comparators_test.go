package sorters

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestComparators(t *testing.T) {
	const (
		a = "a"
		b = "b"

		fra = "aé"
		frb = "aè"

		ia = 1
		ib = 2
	)
	pa := swag.String(a)
	pb := swag.String(b)
	pia := swag.Int(ia)
	pib := swag.Int(ib)
	ba := []byte(a)
	bb := []byte(b)
	bfra := []byte(fra)
	bfrb := []byte(frb)

	t.Run("Reverse[T any](in Comparison[T]) Comparison[T]", func(t *testing.T) {
		comparator := Reverse(StringsComparator())

		require.Equal(t, 0, comparator(a, a))
		require.Equal(t, -1, comparator(b, a))
		require.Equal(t, 1, comparator(a, b))
	})

	t.Run("ReversePtr[T any](in Comparison[*T]) Comparison[*T]", func(t *testing.T) {
		comparator := ReversePtr(StringsPtrComparator())

		require.Equal(t, 0, comparator(nil, nil))
		require.Equal(t, -1, comparator(pa, nil))
		require.Equal(t, 1, comparator(nil, pa))

		require.Equal(t, 0, comparator(pa, pa))
		require.Equal(t, 1, comparator(pb, pa))
		require.Equal(t, -1, comparator(pa, pb))

	})

	t.Run("OrderedComparator[T Ordered]() Comparison[T]", func(t *testing.T) {
		comparator := OrderedComparator[int]()

		require.Equal(t, 0, comparator(ia, ia))
		require.Equal(t, 1, comparator(ib, ia))
		require.Equal(t, -1, comparator(ia, ib))
	})

	t.Run("OrderedPtrComparator[T Ordered]() Comparison[*T]", func(t *testing.T) {
		comparator := OrderedPtrComparator[int]()

		require.Equal(t, 0, comparator(nil, nil))
		require.Equal(t, 1, comparator(pia, nil))
		require.Equal(t, -1, comparator(nil, pia))

		require.Equal(t, 0, comparator(pia, pia))
		require.Equal(t, 1, comparator(pib, pia))
		require.Equal(t, -1, comparator(pia, pib))
	})

	t.Run("BytesComparator(opts ...CompareOption) Comparison[[]byte]", func(t *testing.T) {
		comparator := BytesComparator()

		require.Equal(t, 0, comparator(nil, nil))
		require.Equal(t, 1, comparator(ba, nil))
		require.Equal(t, -1, comparator(nil, ba))

		require.Equal(t, 0, comparator(ba, ba))
		require.Equal(t, 1, comparator(bb, ba))
		require.Equal(t, -1, comparator(ba, bb))
	})

	t.Run("BytesComparator(opts ...CompareOption) with LocaleTag", func(t *testing.T) {
		comparator := BytesComparator(WithLocaleTag(language.French))

		require.Equal(t, 0, comparator(nil, nil))
		require.Equal(t, 1, comparator(bfra, nil))
		require.Equal(t, -1, comparator(nil, bfra))

		require.Equal(t, 0, comparator(bfra, bfra))
		require.Equal(t, 1, comparator(bfrb, bfra))
		require.Equal(t, -1, comparator(bfra, bfrb))
	})

	t.Run("BytesComparator(opts ...CompareOption) with Locale", func(t *testing.T) {
		comparator := BytesComparator(WithLocale("fr"))

		require.Equal(t, 0, comparator(nil, nil))
		require.Equal(t, 1, comparator(bfra, nil))
		require.Equal(t, -1, comparator(nil, bfra))

		require.Equal(t, 0, comparator(bfra, bfra))
		require.Equal(t, 1, comparator(bfrb, bfra))
		require.Equal(t, -1, comparator(bfra, bfrb))
	})

	t.Run("BytesComparator(opts ...CompareOption) with default Locale ", func(t *testing.T) {
		comparator := BytesComparator(WithLocale("zorg"))

		require.Equal(t, 0, comparator(bfra, bfra))
		require.Equal(t, 1, comparator(bfrb, bfra))
		require.Equal(t, -1, comparator(bfra, bfrb))
	})

	t.Run("StringsComparator(opts ...CompareOption) Comparison[string]", func(t *testing.T) {
		comparator := StringsComparator()

		require.Equal(t, 0, comparator(a, a))
		require.Equal(t, 1, comparator(b, a))
		require.Equal(t, -1, comparator(a, b))
	})

	t.Run("StringsPtrComparator(opts ...CompareOption) Comparison[*string]", func(t *testing.T) {
		comparator := StringsPtrComparator()

		require.Equal(t, 0, comparator(nil, nil))
		require.Equal(t, 1, comparator(pa, nil))
		require.Equal(t, -1, comparator(nil, pa))

		require.Equal(t, 0, comparator(pa, pa))
		require.Equal(t, 1, comparator(pb, pa))
		require.Equal(t, -1, comparator(pa, pb))
	})

	t.Run("BoolsComparator(opts ...CompareOption) Comparison[bool]", func(t *testing.T) {
		comparator := BoolsComparator()

		require.Equal(t, 0, comparator(true, true))
		require.Equal(t, 0, comparator(false, false))
		require.Equal(t, 1, comparator(true, false))
		require.Equal(t, -1, comparator(false, true))
	})

	t.Run("BoolsPtrComparator(opts ...CompareOption) Comparison[*bool]", func(t *testing.T) {
		fa := swag.Bool(false)
		tr := swag.Bool(true)

		comparator := BoolsPtrComparator()

		require.Equal(t, 0, comparator(nil, nil))
		require.Equal(t, 1, comparator(tr, nil))
		require.Equal(t, -1, comparator(nil, tr))

		require.Equal(t, 0, comparator(tr, tr))
		require.Equal(t, 0, comparator(fa, fa))
		require.Equal(t, 1, comparator(tr, fa))
		require.Equal(t, -1, comparator(fa, tr))
	})
}
