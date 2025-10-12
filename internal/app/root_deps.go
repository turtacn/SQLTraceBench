package app

import (
	"github.com/turtacn/SQLTraceBench/internal/app/conversion"
	"github.com/turtacn/SQLTraceBench/internal/app/execution"
	"github.com/turtacn/SQLTraceBench/internal/app/generation"
	"github.com/turtacn/SQLTraceBench/internal/app/validation"
	"github.com/turtacn/SQLTraceBench/internal/domain/services"
	"github.com/turtacn/SQLTraceBench/internal/infrastructure/parsers"
)

// Root holds the application's dependencies.
type Root struct {
	// domain services
	Parser services.Parser

	// app services
	Conversion conversion.Service
	Execution  execution.Service
	Generation generation.Service
	Validation validation.Service
}

// NewRoot creates a new Root with all the application's dependencies instantiated.
func NewRoot() *Root {
	// simple instantiation
	parser := parsers.NewAntlrParser()
	return &Root{
		Parser:     parser,
		Conversion: conversion.NewService(parser),
		Execution:  execution.NewService(),
		Generation: generation.NewService(),
		Validation: validation.NewService(),
	}
}