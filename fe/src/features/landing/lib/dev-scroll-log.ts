export function devScrollLog(section: string, progress: number) {
  if (process.env.NODE_ENV === 'development') {
    // eslint-disable-next-line no-console
    console.debug(`[${section}] scroll progress:`, progress.toFixed(3))
  }
}
