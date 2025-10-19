function showError(message) {
  const errorDiv = document.createElement('div')
  errorDiv.className = 'error-toast'
  errorDiv.textContent = message
  document.body.appendChild(errorDiv)

  setTimeout(() => errorDiv.remove(), 5000)
}

function showSuccess(message) {
  const successDiv = document.createElememt('div')
  successDiv.className = 'success-toast'
  successDiv.textContent = message
  document.body.appendChild(successDiv)

  setTimeout(() => successDiv.remove(), 5000)
}

function showWarning(message) {
  const warningDiv = document.createElement('div')
  warningDiv.className = 'success-toast'
  warningDiv.textContent = message
  document.body.appendChild(warningDiv)

  setTimeout(() => warningDiv.remove(), 5000)
}

async function fetchJSON(url, options = {}) {
  const response = await fetch(url, options)

  if (!response.ok) {
    throw new Error(`HTTP ${response.status}: ${response.statusText}`)
  }

  return await response.json()
}
