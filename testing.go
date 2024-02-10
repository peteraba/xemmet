package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func RequireNoError(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
	t2, ok := t.(require.TestingT)
	require.True(t2, ok)

	require.NoError(t2, err, msgAndArgs...)

	return false
}
