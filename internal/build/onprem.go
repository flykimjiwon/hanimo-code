// Copyright 2025-2026 Kim Jiwon (김지원). All rights reserved.
// Licensed under the Apache License, Version 2.0.
// SPDX-License-Identifier: Apache-2.0

//go:build onprem

// Package build exposes compile-time profile flags so the rest of the
// codebase can branch on whether this binary was produced for the on-prem
// LTS distribution channel.
//
// On-prem builds are activated by `go build -tags=onprem` (see Makefile
// target `build-onprem`). They are intended for air-gapped / closed-network
// deployments and assume:
//
//   - zero outbound network calls on startup (no telemetry, no update probes,
//     no font/CDN fetches);
//   - LLM endpoints are user-configured and reachable on the internal network;
//   - desktop variants do not bootstrap WebView2 / install runtimes —
//     the host environment is assumed to provide them.
//
// See docs/policy/lts-onprem.md for the LTS cadence, EOL policy, and
// supported quarterly tags.
package build

const (
	// Onprem is true when this binary was built with the `onprem` tag.
	Onprem = true

	// ProfileName is a short identifier for diagnostics and `/version` output.
	ProfileName = "onprem"
)
