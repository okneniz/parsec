package tml

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShift(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		d      int
		cutoff int
		in     Expr
		want   Expr
	}{
		{
			name:   "below the cutoff stays",
			d:      1,
			cutoff: 1,
			in:     Tuple{Items: []Expr{Var{Index: 0}, Var{Index: 1}}},
			want:   Tuple{Items: []Expr{Var{Index: 0}, Var{Index: 2}}},
		},
		{
			name:   "negative shift",
			d:      -1,
			cutoff: 1,
			in:     Tuple{Items: []Expr{Var{Index: 0}, Var{Index: 2}}},
			want:   Tuple{Items: []Expr{Var{Index: 0}, Var{Index: 1}}},
		},
		{
			name: "a fn binder shields its variable",
			d:    1,
			in: Fn{
				Arg:  VarPat{Name: "x"},
				Body: Tuple{Items: []Expr{Var{Index: 0}, Var{Index: 1}}},
			},
			want: Fn{
				Arg:  VarPat{Name: "x"},
				Body: Tuple{Items: []Expr{Var{Index: 0}, Var{Index: 2}}},
			},
		},
		{
			name: "a tuple pattern binds two",
			d:    1,
			in: Fn{
				Arg: TuplePat{Items: []Pat{
					VarPat{Name: "a"},
					VarPat{Name: "b"},
				}},
				Body: Tuple{Items: []Expr{
					Var{Index: 0},
					Var{Index: 1},
					Var{Index: 2},
				}},
			},
			want: Fn{
				Arg: TuplePat{Items: []Pat{
					VarPat{Name: "a"},
					VarPat{Name: "b"},
				}},
				Body: Tuple{Items: []Expr{
					Var{Index: 0},
					Var{Index: 1},
					Var{Index: 3},
				}},
			},
		},
		{
			name: "let declarations raise the cutoff for the body",
			d:    1,
			in: Let{
				Decls: []Decl{ValDecl{Binds: []ValBind{
					{Pat: VarPat{Name: "x"}, E: Var{Index: 1}},
				}}},
				Body: Tuple{Items: []Expr{Var{Index: 0}, Var{Index: 1}}},
			},
			want: Let{
				Decls: []Decl{ValDecl{Binds: []ValBind{
					{Pat: VarPat{Name: "x"}, E: Var{Index: 2}},
				}}},
				Body: Tuple{Items: []Expr{Var{Index: 0}, Var{Index: 2}}},
			},
		},
		{
			name: "a fun binds its name and its arguments",
			d:    1,
			in: Let{
				Decls: []Decl{FunDecl{
					Name: "f",
					Args: []Pat{VarPat{Name: "x"}},
					Body: Tuple{Items: []Expr{
						Var{Index: 0},
						Var{Index: 1},
						Var{Index: 2},
					}},
				}},
				Body: Tuple{Items: []Expr{Var{Index: 0}, Var{Index: 1}}},
			},
			want: Let{
				Decls: []Decl{FunDecl{
					Name: "f",
					Args: []Pat{VarPat{Name: "x"}},
					Body: Tuple{Items: []Expr{
						Var{Index: 0},
						Var{Index: 1},
						Var{Index: 3},
					}},
				}},
				Body: Tuple{Items: []Expr{Var{Index: 0}, Var{Index: 2}}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := shift(test.d, test.cutoff, test.in)

			require.Equal(t, test.want, got)
		})
	}
}

func TestSubst(t *testing.T) {
	t.Parallel()

	one := IntLit{Text: "1", Value: 1}

	tests := []struct {
		name string
		j    int
		v    Expr
		in   Expr
		want Expr
	}{
		{
			name: "the hit replaces, the miss keeps its index",
			j:    1,
			v:    one,
			in:   Tuple{Items: []Expr{Var{Index: 1}, Var{Index: 0}}},
			want: Tuple{Items: []Expr{one, Var{Index: 0}}},
		},
		{
			name: "the value climbs one binder",
			j:    0,
			v:    Var{Index: 5},
			in: Fn{
				Arg:  VarPat{Name: "x"},
				Body: Tuple{Items: []Expr{Var{Index: 0}, Var{Index: 1}}},
			},
			want: Fn{
				Arg:  VarPat{Name: "x"},
				Body: Tuple{Items: []Expr{Var{Index: 0}, Var{Index: 6}}},
			},
		},
		{
			name: "the value climbs the binders of a let",
			j:    0,
			v:    one,
			in: Let{
				Decls: []Decl{ValDecl{Binds: []ValBind{
					{Pat: VarPat{Name: "x"}, E: Var{Index: 0}},
				}}},
				Body: Var{Index: 0},
			},
			want: Let{
				Decls: []Decl{ValDecl{Binds: []ValBind{
					{Pat: VarPat{Name: "x"}, E: one},
				}}},
				Body: Var{Index: 0},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := subst(test.j, test.v, test.in)

			require.Equal(t, test.want, got)
		})
	}
}

func TestSubstitutionTop(t *testing.T) {
	t.Parallel()

	three := IntLit{Text: "3", Value: 3}

	tests := []struct {
		name string
		v    Expr
		in   Expr
		want Expr
	}{
		{
			name: "the beta step",
			v:    three,
			in:   Infix{Op: "+", L: Var{Index: 0}, R: Var{Index: 1}},
			want: Infix{Op: "+", L: three, R: Var{Index: 0}},
		},
		{
			name: "the binder disappears and the rest drops",
			v:    three,
			in:   Tuple{Items: []Expr{Var{Index: 0}, Var{Index: 1}}},
			want: Tuple{Items: []Expr{three, Var{Index: 0}}},
		},
		{
			name: "a closed body keeps nothing",
			v:    three,
			in:   Tuple{Items: []Expr{Var{Index: 0}, IntLit{Text: "9", Value: 9}}},
			want: Tuple{Items: []Expr{three, IntLit{Text: "9", Value: 9}}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := substitutionTop(test.in, test.v)

			require.Equal(t, test.want, got)
		})
	}
}

func TestSubstitutionTops(t *testing.T) {
	t.Parallel()

	one := IntLit{Text: "1", Value: 1}
	two := IntLit{Text: "2", Value: 2}
	three := IntLit{Text: "3", Value: 3}

	tests := []struct {
		name string
		vs   []Expr
		in   Expr
		want Expr
	}{
		{
			name: "the values a tuple pattern matched replace its binders",
			vs:   []Expr{one, two},
			in: Tuple{Items: []Expr{
				Var{Index: 1},
				Var{Index: 0},
				Var{Index: 2},
			}},
			want: Tuple{Items: []Expr{one, two, Var{Index: 0}}},
		},
		{
			name: "three binders at once",
			vs:   []Expr{one, two, three},
			in: Tuple{Items: []Expr{
				Var{Index: 0},
				Var{Index: 1},
				Var{Index: 2},
				Var{Index: 3},
			}},
			want: Tuple{Items: []Expr{three, two, one, Var{Index: 0}}},
		},
		{
			name: "the value is a term with free indices",
			vs:   []Expr{Var{Index: 7}},
			in: Fn{
				Arg:  VarPat{Name: "x"},
				Body: Tuple{Items: []Expr{Var{Index: 0}, Var{Index: 1}}},
			},
			want: Fn{
				Arg:  VarPat{Name: "x"},
				Body: Tuple{Items: []Expr{Var{Index: 0}, Var{Index: 8}}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := substitutionTops(test.in, test.vs)

			require.Equal(t, test.want, got)
		})
	}
}
