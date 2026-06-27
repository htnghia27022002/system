# Implementation Plan: Authentication

**Feature**: `001-auth` | **Date**: 2026-06-27 | **Spec**: [spec.md](spec.md)

## Summary

Complete auth UX: email login/register (mostly done), JWT refresh (done), Google OAuth FE flow, logout with refresh revocation, modern auth UI.

## Technical context

| Layer | Stack |
|-------|-------|
| BE | Go 1.22, Gin, JWT, golang.org/x/oauth2 |
| FE | Next.js 15, React 19, TanStack Query, Zustand, shadcn/ui |
| API | `/api/auth/*`, camelCase JSON |

## Approach

1. **BE**: Small addition — OAuth providers discovery endpoint; core auth unchanged.
2. **FE**: OAuth redirect → `/auth/callback` → POST callback → store tokens; polish auth shell.
3. **No schema changes** expected for v1.

## Constitution check

- Package independence: BE/FE communicate via HTTP only — pass
- English docs — pass
- Spec → tasks → plan → implement — pass

See [be-implement.md](be-implement.md) and [fe-implement.md](fe-implement.md) for file-level detail.
