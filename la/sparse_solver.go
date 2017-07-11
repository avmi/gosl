// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package la

import (
	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/mpi"
)

// real ////////////////////////////////////////////////////////////////////////////////////////////

// SparseSolver solves sparse linear systems using UMFPACK or MUMPS
//
//   Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b
//
type SparseSolver interface {
	Init(t *Triplet, symmetric, verbose bool, ordering, scaling string, comm *mpi.Communicator) error
	Free()
	Fact() error
	Solve(x, b Vector, bIsDistr bool) error
}

// spSolverMaker defines a function that makes spSolvers
type spSolverMaker func() SparseSolver

// spSolverDB implements a database of SparseSolver makers
var spSolverDB = make(map[string]spSolverMaker)

// NewSparseSolver finds a SparseSolver in database or panic
//   kind -- "umfpack" or "mumps"
func NewSparseSolver(kind string) SparseSolver {
	if maker, ok := spSolverDB[kind]; ok {
		return maker()
	}
	chk.Panic("cannot find SparseSolver named %q in database", kind)
	return nil
}

// complex /////////////////////////////////////////////////////////////////////////////////////////

// SparseSolverC solves sparse linear systems using UMFPACK or MUMPS (complex version)
//
//   Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b
//
type SparseSolverC interface {
	Init(t *TripletC, symmetric, verbose bool, ordering, scaling string, comm *mpi.Communicator) error
	Free()
	Fact() error
	Solve(x, b VectorC, bIsDistr bool) error
}

// spSolverMakerC defines a function that makes spSolvers (complex version)
type spSolverMakerC func() SparseSolverC

// spSolverDBc implements a database of SparseSolver makers (complex version)
var spSolverDBc = make(map[string]spSolverMakerC)

// NewSparseSolverC finds a SparseSolver in database or panic
func NewSparseSolverC(kind string) SparseSolverC {
	if maker, ok := spSolverDBc[kind]; ok {
		return maker()
	}
	chk.Panic("cannot find SparseSolverC named %q in database", kind)
	return nil
}

// high-level functions ////////////////////////////////////////////////////////////////////////////

// SpSolve solves a sparse linear system (using UMFPACK)
//
//   Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b
//
func SpSolve(A *Triplet, b Vector) (x Vector, err error) {

	// allocate solver
	o := NewSparseSolver("umfpack")
	defer o.Free()

	// initialise solver
	err = o.Init(A, false, false, "", "", nil)
	if err != nil {
		return
	}

	// factorise
	err = o.Fact()
	if err != nil {
		return
	}

	// solve
	x = NewVector(len(b))
	err = o.Solve(x, b, false) // x := inv(A) * b
	return
}

// SpSolveC solves a sparse linear system (using UMFPACK) (complex version)
//
//   Given:  A ⋅ x = b    find x   such that   x = A⁻¹ ⋅ b
//
func SpSolveC(A *TripletC, b VectorC) (x VectorC, err error) {

	// allocate solver
	o := NewSparseSolverC("umfpack")
	defer o.Free()

	// initialise solver
	err = o.Init(A, false, false, "", "", nil)
	if err != nil {
		return
	}

	// factorise
	err = o.Fact()
	if err != nil {
		return
	}

	// solve
	x = NewVectorC(len(b))
	err = o.Solve(x, b, false) // x := inv(A) * b
	return
}
