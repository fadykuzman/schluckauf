function showError(message) {
  const errorDiv = document.createElement('div')
  errorDiv.className = 'error-toast'
  errorDiv.textContent = message
  document.body.appendChild(errorDiv)

  setTimeout(() => errorDiv.remove(), 5000)
}

async function fetchJSON(url, options = {}) {
  const response = await fetch(url, options)

  if (!response.ok) {
    throw new Error(`HTTP ${response.status}: ${response.statusText}`)
  }

  return await response.json()
}
