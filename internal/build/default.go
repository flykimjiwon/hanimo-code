// Copyright 2025-2026 Kim Jiwon (김지원). All rights reserved.
// Licensed under the Apache License, Version 2.0.
// SPDX-License-Identifier: Apache-2.0

//go:build !onprem

// Package build exposes compile-time profile flags. This file is the
// default (non-onprem) variant — see onprem.go for the closed-network
// counterpart and docs/policy/lts-onprem.md for the LTS strategy.
package build

const (
	// Onprem is false for the default open-network distribution.
	Onprem = false

	// ProfileName is a short identifier for diagnostics and `/version` output.
	ProfileName = "default"
)
