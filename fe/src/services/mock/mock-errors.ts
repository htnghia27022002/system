// Error types thrown by the mock API layer. Kept in a tiny module so app code can
// `instanceof`-check them without statically importing (and bundling) the mocks.

export class MockAuthError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'MockAuthError'
    this.status = status
  }
}

export class MockAccessControlError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'MockAccessControlError'
    this.status = status
  }
}
